package main

import (
	def "config"
	"errors"
	"fmt"
	"fsm"
	"hw"
	"liftAssigner"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
	"time"
)

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

var onlineLifts = make(map[string]network.UdpConnection)
var numberOfOnlineLifts int

var deadChan = make(chan network.UdpConnection)
var costChan = make(chan def.Message)
var outgoingMsg = make(chan def.Message, 10)
var incomingMsg = make(chan def.Message, 10)

func main() {
	var floor int
	var err error
	floor, err = hw.Init()
	if err != nil {
		def.Restart.Run()
		log.Fatal(err)
	}

	e := fsm.Channels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan int, 10),
		FloorLamp:    make(chan int, 10),
		DoorLamp:     make(chan bool, 10),
		OutgoingMsg:  outgoingMsg,
	}
	fsm.Init(e, floor)

	network.Init(outgoingMsg, incomingMsg)

	go liftAssigner.Run(costChan, &numberOfOnlineLifts)
	go eventHandler(e)
	go syncLights()

	queue.Init(e.NewOrder, outgoingMsg)

	// Handle CTRL+C
	go safeKill() //bad name?

	for { //nicer solution?
		time.Sleep(100 * time.Second)
	}
}

func eventHandler(c fsm.Channels) {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			switch keypress.Button {
			case def.BtnInside:
				queue.AddLocalOrder(keypress.Floor, keypress.Button)
			case def.BtnUp, def.BtnDown:
				outgoingMsg <- def.Message{
					Category: def.NewOrder,
					Floor:    keypress.Floor,
					Button:   keypress.Button}
			}
		case floor := <-floorChan:
			c.FloorReached <- floor
		case message := <-incomingMsg:
			handleMessage(message)
		case connection := <-deadChan:
			handleDeadLift(connection.Addr)
		case order := <-queue.OrderTimeoutChan:
			log.Println(def.ColR, "Order timeout, I can do it myself!", def.ColN)
			queue.RemoveRemoteOrdersAt(order.Floor)
			queue.AddRemoteOrder(order.Floor, order.Button, def.Laddr)
		case dir := <-c.MotorDir:
			hw.SetMotorDirection(dir)
		case floor := <-c.FloorLamp:
			hw.SetFloorLamp(floor)
		case value := <-c.DoorLamp:
			hw.SetDoorOpenLamp(value)
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
					if (f == 0 && b == def.BtnDown) ||
						(f == def.NumFloors-1 && b == def.BtnUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- def.Keypress{Button: b, Floor: f}
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

func handleMessage(message def.Message) { // consider moving each case into a function
	const aliveTimeout = 2 * time.Second

	switch message.Category {
	case def.Alive:
		if connection, exist := onlineLifts[message.Addr]; exist {
			connection.Timer.Reset(aliveTimeout)
		} else {
			newConnection := network.UdpConnection{message.Addr, time.NewTimer(aliveTimeout)}
			onlineLifts[message.Addr] = newConnection
			numberOfOnlineLifts = len(onlineLifts)
			go connectionTimer(&newConnection)
			log.Printf("%sConnection to IP %s established!%s", def.ColG, message.Addr[0:15], def.ColN)
		}
	case def.NewOrder:
		// log.Printf("handleMessage(): NewOrder message: f=%d b=%d from lift %s\n",
		//	message.Floor+1, message.Button, message.Addr[12:15])

		cost := queue.CalculateCost(message.Floor, message.Button, fsm.Floor(), hw.Floor(), fsm.Direction())

		costMessage := def.Message{
			Category: def.Cost,
			Floor:    message.Floor,
			Button:   message.Button,
			Cost:     cost}
		// log.Printf("handleMessage(): NewOrder sends cost message: f=%d b=%d (with cost %d) from me\n", costMessage.Floor+1, costMessage.Button, costMessage.Cost)
		outgoingMsg <- costMessage
	case def.CompleteOrder:
		queue.RemoveRemoteOrdersAt(message.Floor)
	case def.Cost:
		// log.Printf("handleMessage(): Cost message: f=%d b=%d with cost %d from lift %s\n", message.Floor+1, message.Button, message.Cost, message.Addr[12:15])
		costChan <- message
	}
}

// handleDeadLift removes the lift that have timed out from the onlineLifts
// and reassigns the dead lifts orders
func handleDeadLift(deadAddr string) {
	log.Printf("%sConnection to IP %s is dead!%s", def.ColR, deadAddr[0:15], def.ColN)
	delete(onlineLifts, deadAddr)
	numberOfOnlineLifts = len(onlineLifts)
	queue.ReassignOrders(deadAddr, outgoingMsg)
}

// connectionTimer is a go-routine for detecting that lifts aren't on the network anymore 1
func connectionTimer(connection *network.UdpConnection) {
	<-connection.Timer.C
	deadChan <- *connection
}

// safeKill stops the motor if the program is killed with CTRL+C.
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
	hw.SetMotorDirection(def.DirStop)
	log.Fatal(def.ColR, "User terminated program", def.ColN)
}

func syncLights() {
	for {
		<-def.SyncLightsChan

		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if (b == def.BtnUp && f == def.NumFloors-1) ||
					(b == def.BtnDown && f == 0) {
					continue
				} else {
					switch b {
					case def.BtnInside:
						hw.SetButtonLamp(f, b, queue.IsLocalOrder(f, b))
					case def.BtnUp, def.BtnDown:
						hw.SetButtonLamp(f, b, queue.IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
