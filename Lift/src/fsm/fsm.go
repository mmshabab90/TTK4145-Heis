// This finite state machine is based on code and ideas presented in
// Rob Pike's talk 'Lexical Scanning in Go':
// https://www.youtube.com/watch?v=HxaD_trXwRE
package fsm

import (
	def "../config"
	"../queue"
	"log"
)

type lift struct {
	floor int
	dir   int
}

// TODO These must be connected to the rest of the system (just placeholders now)
var motorDir = make(chan int)
var doorOpenLamp = make(chan bool)

// stateFunc represents the state of the lift
// as a function that returns the next state.
type stateFunc func(*lift) stateFunc

var eventNewOrder = make(chan bool)
var eventFloorReached = make(chan int)
var eventDoorTimeout = make(chan bool)

func Init(startFloor int) (eventNewOrder <-chan bool, eventFloorReached <-chan int) {
	l := &lift{
		floor: startFloor,
		dir:   def.DirStop,
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
	l.dir = def.DirStop
	motorDir <- def.DirStop
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
	l.dir = def.DirStop
	motorDir <- def.DirStop
	doorOpenLamp <- true
	// orderComplete <- l.floor // order complete message should be invoked from queue package
	doorReset <- true
	queue.LiftArrivedAt(l.floor)

	select {

	case <-eventNewOrder:
		if queue.ShouldStop(l.floor, l.dir) {
			l.dir = def.DirStop
			doorReset <- true
			queue.LiftArrivedAt(l.floor)
		}
		return open

	case l.floor = <-eventFloorReached:
		log.Printf("Makes no sense to arrive at a floor (%v) in state door open.\n", l.floor)
		return open

	case <-eventDoorTimeout:
		if l.dir = queue.ChooseDirection(l.floor, l.dir); l.dir == def.DirStop {
			return idle
		} else {
			return moving
		}
	}
}
