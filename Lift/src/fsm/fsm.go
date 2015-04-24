package fsm

import (
	def "../config"
	"../hw"
	"../queue"
	"fmt"
	"log"
	"time"
)

const (
	idle int = iota
	moving
	doorOpen
)

const doorOpenTime = 1 * time.Second

var state int
var floor int
var dir int

type EventChannels struct {
	NewOrder     chan bool
	FloorReached chan int
	DoorTimeout  chan bool
}

var doorReset = make(chan bool)

func Init(e EventChannels) {
	log.Println("fsm.Init() starting")

	state = idle
	dir = def.DirStop
	floor = hw.Floor()

	go syncLights()

	e.DoorTimeout = make(chan bool)
	go startDoorTimer(e.DoorTimeout)
	go run(e)
}

func run(e EventChannels) {
	for {
		select {
		case <-e.NewOrder:
			eventNewOrder()
		case f := <-e.FloorReached:
			eventFloorReached(f)
		case <-e.DoorTimeout:
			eventDoorTimeout()
		}
	}
}

func eventNewOrder() {
	fmt.Printf("\nEVENT: New order in state %v.\n\n", stateString(state))
	switch state {
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(floor, dir) {
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
			hw.SetMotorDirection(dir)
			state = moving
		}
	case moving:
		// ignore
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			fmt.Println("SPAM!")
			queue.RemoveOrdersAt(floor)
			//go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
		}
	}
}

func eventFloorReached(newFloor int) {
	fmt.Printf("\nEVENT: Floor %d reached in state %s.\n\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			dir = def.DirStop
			hw.SetMotorDirection(dir)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
		}
	default:
		log.Printf("Makes no sense to arrive at a floor in state %s.\n", stateString(state))
	}
}

func eventDoorTimeout() { //this happens for each external order
	fmt.Printf("\nEVENT: Door timeout in state %s.\n\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		dir = queue.ChooseDirection(floor, dir)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(dir)
		if dir == def.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
}
