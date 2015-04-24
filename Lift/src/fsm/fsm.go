package fsm

import (
	"../defs"
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
	dir = defs.DirStop
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
			dir = queue.ChooseDirection(floor, dir)
			state = moving
		}
	case moving:
		// ignore
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			doorReset <- true
		}
	}
}

func eventFloorReached(newFloor int) {
	fmt.Printf("\n\n   ☺      Event floor %d reached in state %s\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			hw.SetMotorDirection(defs.DirStop)
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
	defs.SyncLightsChan <- true
}

func eventDoorTimeout() { //this happens for each external order
	fmt.Printf("\n\n   ☺      Event door timeout in state %s\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		dir = queue.ChooseDirection(floor, dir)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(dir)
		if dir == defs.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
	defs.SyncLightsChan <- true
}
