package fsm

import (
	"log"
)

type keypress struct {
	floor  int
	button int
}

const (
	dirDown int = iota - 1
	dirStop
	dirUp
)

const (
	idle int = iota
	moving
	open
)

var floor int
var dir int
var departDir int

func Init(startFloor int,
	eventNewOrder <-chan bool,
	eventFloorReached <-chan int) {
	log.Println("fsm.Init()")

	eventDoorTimeout := make(chan bool)

	state = idle
	floor = startFloor
	dir = dirStop
	departDir = dirDown

	go run(eventNewOrder, eventFloorReached, eventDoorTimeout)
}

func run(eventNewOrder <-chan bool,
	eventFloorReached <-chan int,
	eventDoorTimeout <-chan bool) {
	for {
		select {
		case <-eventNewOrder:
			eventNewOrder()
		case floor := <-eventFloorReached:
			eventFloorReached(floor)
		case <-eventDoorTimeout:
			eventDoorTimeout()
		}
	}
}

func eventNewOrder() {
	switch state {
	case idle:
		switch dir = queue.ChooseDirection(floor, dir); dir {
		case dirStop:
			eventFloorReached <- floor
		case dirDown, dirUp:
			// motorDir <- dir
			departDir = dir
			state = moving
		}
	case open:
		if queue.ShouldStop() {
			dir = dirStop // redundant
			doorReset <- true
			queue.RemoveOrdersAt(floor)
		}
	case moving:
		// ignore
	}
}

func eventFloorReached(reachedFloor int) {
	floor = reachedFloor
	// floorLamp <- floor
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			// motorDir <- dir
			// doorOpenLamp <- true
			queue.RemoveOrdersAt(floor)
			// orderComplete <- floor
			doorReset <- true
			state = open
		} else {
			departDir = dir
		}
	default:
		log.Printf("Makes no sense to arrive at a floor in state %s.\n", stateString(state))
	}
}

func eventDoorTimeout() {
	// doorOpenLamp <- false
	switch state {
	case open:
		dir = queue.ChooseDirection(floor, dir)
		// motorDir <- dir
		if dir == dirStop {
			state = idle
		} else {
			state = moving
			departDir = dir
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
}
