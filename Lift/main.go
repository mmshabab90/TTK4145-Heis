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
)

const debugPrint = false

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

// var orderTimeoutChan = make(chan order)

func main() {
	motorDir := make(chan int)
	doorOpenLamp := make(chan bool)
	floorLamp := make(chan int)
	deathChan := make(chan string)
	floorCompleted := queue.Init(floorLamp, deathChan)

	floor := hw.Init(motorDir, doorOpenLamp, floorLamp)
	eventNewOrder, eventFloorReached := fsm.Init(floor)

	messageChan := make(chan network.Message)
	addRemoteOrder := make(chan network.RemoteOrder)
	network.Init(orderComplete, deathChan, messageChan, addRemoteOrder)
	run(eventNewOrder, eventFloorReached, messageChan)
}

func run(eventNewOrder <-chan bool, eventFloorReached <-chan int, messageChan <-chan network.Message) {
	
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case key := <-buttonChan:
			queue.NewKeypress(key)
			eventNewOrder <- true
		case floor := <-floorChan:
			eventFloorReached <- floor
		case msg := <-messageChan:
			handleMessage(msg)
		case order := <-addRemoteOrder:
			queue.AddRemoteOrder(order.Floor, order.Button, order.Addr)
		}
	}
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
							c <- defKeypress{Floor: f, Button: b}
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
