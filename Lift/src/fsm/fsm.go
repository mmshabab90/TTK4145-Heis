// This finite state machine is based on code and ideas presented in Rob
// Pike's talk 'Lexical Scanning in Go':
// https://www.youtube.com/watch?v=HxaD_trXwRE
package fsm

import (
	"../queue"
	"log"
)

const (
	dirDown int = iota - 1
	dirStop
	dirUp
)

type lift struct {
	floor int
	dir   int
}

// stateFunc represents the state of the lift
// as a function that returns the next state.
type stateFunc func(*lift) stateFunc

var eventNewOrder = make(chan bool)
var eventFloorReached = make(chan int)
var eventDoorTimeout = make(chan bool)

func Init(startFloor int) (eventNewOrder <-chan bool, eventFloorReached <-chan int) {
	l := &lift{
		floor: startFloor,
		dir:   dirStop,
	}

	go l.run()

	return eventNewOrder, eventFloorReached
}

func (l *lift) run() {
	for state := idle; state != nil; {
		state = state(l)
	}
}

func idle(l *lift) stateFunc {
	l.dir = dirStop
	motorDir <- dirStop
	doorOpenLamp <- true

	select {

	case <-eventNewOrder:
		if queue.ShouldStop(l.floor, l.dir) {
			return open
		} else {
			l.dir = queue.ChooseDirection(l.floor, l.dir)
			return moving
		}

	case l.floor = <-eventFloorReached:
		log.Printf("Makes no sense to arrive at a floor (%v) in state idle.\n", l.floor)
		return idle

	case <-eventDoorTimeout:
		log.Printf("Makes no sense to time out door timer when in state idle.\n")
		return idle
	}
}

func moving(l *lift) stateFunc {
	motorDir <- l.dir
	doorOpenLamp <- false

	select {

	case <-eventNewOrder:
		return moving

	case l.floor = <-eventFloorReached:
		if queue.ShouldStop(l.floor, l.dir) {
			return open
		} else {
			l.dir = queue.ChooseDirection(l.floor, l.dir)
			return moving
		}

	case <-eventDoorTimeout:
		log.Println("Makes no sense to time out door timer when in state moving.")
		return moving
	}
}

func open(l *lift) stateFunc {
	l.dir = dirStop
	motorDir <- dirStop
	doorOpenLamp <- true
	orderComplete <- l.floor
	doorReset <- true
	queue.RemoveOrdersAt(l.floor) // maybe this should happen via 'orderComplete <- l.floor'

	select {

	case <-eventNewOrder:
		if queue.ShouldStop() {
			l.dir = l.dirStop // redundant
			doorReset <- true
			queue.RemoveOrdersAt(l.floor)
		}
		return open

	case l.floor = <-eventFloorReached:
		log.Printf("Makes no sense to arrive at a floor (%v) in state door open.\n", l.floor)
		return open

	case <-eventDoorTimeout:
		if l.dir = queue.ChooseDirection(l.floor, l.dir); l.dir == dirStop {
			return idle
		} else {
			return moving
		}
	}
}
