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
	floor             int
	dir               int
	motorDir          chan<- int
	doorOpenLamp      chan<- bool
	floorCompleted    chan<- int
	doorReset         chan bool
	eventNewOrder     <-chan bool
	eventFloorReached <-chan int
	eventDoorTimeout  <-chan bool
}

// stateFunc represents the state of the lift
// as a function that returns the next state.
type stateFunc func(*lift) stateFunc

var l *lift

var eventDoorTimeout = make(chan bool)

func Init(
	startFloor int,
	motorDir chan<- int,
	doorOpenLamp chan<- bool,
	floorCompleted chan int) (
	eventNewOrder chan bool,
	eventFloorReached chan int) {

	eventNewOrder = make(chan bool)
	eventFloorReached = make(chan int)

	doorReset := make(chan bool)

	l = &lift{
		floor:             startFloor,
		dir:               def.DirStop,
		motorDir:          motorDir,
		doorOpenLamp:      doorOpenLamp,
		floorCompleted:    floorCompleted,
		doorReset:         doorReset,
		eventNewOrder:     eventNewOrder,
		eventFloorReached: eventFloorReached,
		eventDoorTimeout:  eventDoorTimeout,
	}
	
	go startTimer(doorReset, eventDoorTimeout)
	go l.run()

	log.Println("fsm.Init() returning...")
	return eventNewOrder, eventFloorReached
}

func (l *lift) run() {
	for state := idle; state != nil; {
		state = state(l)
	}
}

func idle(l *lift) stateFunc {
	log.Println("State idle")
	l.dir = def.DirStop
	l.motorDir <- def.DirStop
	l.doorOpenLamp <- false
	log.Println("State idle will select now")
	select {

	case <-l.eventNewOrder:
		if queue.ShouldStop(l.floor, l.dir) {
			return open
		} else {
			l.dir = queue.ChooseDirection(l.floor, l.dir)
			return moving
		}

	case l.floor = <-l.eventFloorReached:
		log.Printf("Makes no sense to arrive at a floor (%v) in state idle.\n", l.floor)
		return idle

	case <-l.eventDoorTimeout:
		log.Printf("Makes no sense to time out door timer when in state idle.\n")
		return idle
	}
}

func moving(l *lift) stateFunc {
	log.Println("State moving")
	l.motorDir <- l.dir
	l.doorOpenLamp <- false

	select {

	case <-l.eventNewOrder:
		return moving

	case l.floor = <-l.eventFloorReached:
		if queue.ShouldStop(l.floor, l.dir) {
			return open
		} else {
			l.dir = queue.ChooseDirection(l.floor, l.dir)
			return moving
		}

	case <-l.eventDoorTimeout:
		log.Println("Makes no sense to time out door timer when in state moving.")
		return moving
	}
}

func open(l *lift) stateFunc {
	log.Println("State open")
	l.dir = def.DirStop
	l.motorDir <- def.DirStop
	l.doorOpenLamp <- true
	l.doorReset <- true
	queue.LiftArrivedAt(l.floor, l.floorCompleted)
	log.Println("2349507234623457913845")
	log.Println("will select")
	select {

	case <-l.eventNewOrder:
		if queue.ShouldStop(l.floor, l.dir) {
			l.dir = def.DirStop
			l.doorReset <- true
			queue.LiftArrivedAt(l.floor, l.floorCompleted)
		}
		return open

	case l.floor = <-l.eventFloorReached:
		log.Printf("Makes no sense to arrive at a floor (%v) in state door open.\n", l.floor)
		return open

	case <-l.eventDoorTimeout:
		if l.dir = queue.ChooseDirection(l.floor, l.dir); l.dir == def.DirStop {
			return idle
		} else {
			return moving
		}
	}
}

func Floor() int {
	return l.floor
}

func Dir() int {
	return l.dir
}
