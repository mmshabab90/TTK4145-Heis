package main

import (
	def "./src/config"
	"./src/fsm"
	"./src/hw"
	"./src/network"
	"./src/queue"
	"errors"
	"fmt"
	"log"
	"time"
)

const debugPrint = false

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

var onlineLifts = make(map[string]*time.Timer)
var resetTime = 30 * time.Second

// var orderTimeoutChan = make(chan order)

func main() {
	log.Println("##################################################")
	motorDir := make(chan int)
	doorOpenLamp := make(chan bool)
	floorLamp := make(chan int)
	setButtonLamp := make(chan def.Keypress)
	deathChan := make(chan string)
	keypressChan := make(chan def.Keypress)
	messageChan := make(chan network.Message)
	addRemoteOrder := make(chan network.RemoteOrder)
	costMessage := make(chan network.Message)
	floorCompleted := make(chan int)

	floor := hw.Init(motorDir,
		doorOpenLamp,
		floorLamp)

	eventNewOrder, eventFloorReached := fsm.Init(
		floor,
		motorDir,
		doorOpenLamp
		floorCompleted)

	queue.Init(
		setButtonLamp,
		deathChan,
		keypressChan,
		eventNewOrder
		floorCompleted)

	network.Init(
		floorCompleted,
		deathChan,
		addRemoteOrder,
		costMessage,
		onlineLifts)

	run(
		eventNewOrder,
		eventFloorReached,
		messageChan,
		deathChan,
		costMessage,
		addRemoteOrder,
		keypressChan)
}

func run(
	eventNewOrder chan bool,
	eventFloorReached chan int,
	messageChan <-chan network.Message,
	deathChan chan<- string,
	costMessage chan<- network.Message,
	addRemoteOrder <-chan network.RemoteOrder,
	keypressChan chan <- def.Keypress) {

	buttonChan := pollButtons()
	floorChan := pollFloors()

	log.Println("Running...")

	for {
		select {
		case key := <-buttonChan:
			keypressChan <- key
		case floor := <-floorChan:
			eventFloorReached <- floor
		case msg := <-messageChan:
			handleMessage(deathChan, costMessage, msg)
		case order := <-addRemoteOrder:
			queue.AddRemoteOrder(order.Floor, order.Button, order.Addr)
		}
	}
}

func handleMessage(deathChan chan<- string, costMessage chan<- network.Message, msg network.Message) {
	switch msg.Kind {

	case network.Alive:
		if timer, exist := onlineLifts[msg.Addr]; exist {
			timer.Reset(resetTime)
		} else {
			timer := time.NewTimer(resetTime)
			onlineLifts[msg.Addr] = timer
			go waitForDeath(deathChan, onlineLifts, msg.Addr)
		}

	case network.NewOrder:
		cost := queue.CalculateCost(msg.Floor, msg.Button, fsm.Floor(), hw.Floor(), fsm.Dir())
		network.Outgoing <- network.Message{
			Kind:   network.Cost,
			Floor:  msg.Floor,
			Button: msg.Button,
			Cost:   cost}

	case network.CompleteOrder:
		queue.RemoveRemoteOrdersAt(msg.Floor)

	case network.Cost:
		costMessage <- msg
	}
}

func waitForDeath(deathChan chan<- string, onlineLifts map[string]*time.Timer, deadAddr string) {
	<-onlineLifts[deadAddr].C
	delete(onlineLifts, deadAddr)
	deathChan <- deadAddr
}

func pollButtons() <-chan def.Keypress {
	c := make(chan def.Keypress)

	go func() {
		var buttonState [def.NumFloors][def.NumButtons]bool

		for {
			for f := 0; f < def.NumFloors; f++ {
				for b := 0; b < def.NumButtons; b++ {
					if (f == 0 && b == def.ButtonDown) ||
						(f == def.NumFloors-1 && b == def.ButtonUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- def.Keypress{Floor: f, Button: b}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

func pollFloors() <-chan int {
	c := make(chan int)
	go func() {
		oldFloor := hw.Floor()

		for {
			newFloor := hw.Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()
	return c
}
