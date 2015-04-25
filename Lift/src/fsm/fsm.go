// Package fsm implements a finite state machine for the behaviour of a lift.
// The lift runs based on a queue stored and managed by the queue package.
package fsm

import (
	def "config"
	"log"
	"queue"
)

const (
	idle int = iota
	moving
	doorOpen
)

var state int
var floor int
var dir int

type Channels struct {
	// Events
	NewOrder     chan bool
	FloorReached chan int
	doorTimeout  chan bool
	// Hardware interaction
	MotorDir  chan int
	FloorLamp chan int
	DoorLamp  chan bool
	// Door timer
	doorTimerReset chan bool
}

func Init(c Channels, startFloor int) {
	state = idle
	dir = def.DirStop
	floor = startFloor

	go syncLights()

	c.doorTimeout = make(chan bool)
	c.doorTimerReset = make(chan bool)

	go doorTimer(c.doorTimeout, c.doorTimerReset)
	go run(c)

	log.Println("FSM initialised.")
}

func run(c Channels) {
	for {
		select {
		case <-c.NewOrder:
			eventNewOrder(c)
		case f := <-c.FloorReached:
			eventFloorReached(c, f)
		case <-c.doorTimeout:
			eventDoorTimeout(c)
		}
	}
}

func eventNewOrder(e Channels) {
	log.Printf("EVENT: New order in state %v.\n\n", stateString(state))
	switch state {
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(floor, dir) {
			e.DoorLamp <- true
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			e.doorTimerReset <- true
			state = doorOpen
		} else {
			e.MotorDir <- dir
			state = moving
		}
	case moving:
		// Ignore
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			queue.RemoveOrdersAt(floor)
			e.doorTimerReset <- true
		}
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("This state doesn't exist")
	}
}

func eventFloorReached(e Channels, newFloor int) {
	log.Printf("EVENT: Floor %d reached in state %s.\n\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	e.FloorLamp <- floor
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			dir = def.DirStop
			e.MotorDir <- dir
			e.DoorLamp <- true
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			e.doorTimerReset <- true
			state = doorOpen
		}
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("Makes no sense to arrive at a floor in state %s.\n", stateString(state))
	}
}

func eventDoorTimeout(e Channels) {
	log.Printf("EVENT: Door timeout in state %s.\n\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		dir = queue.ChooseDirection(floor, dir)
		e.DoorLamp <- false
		e.MotorDir <- dir
		if dir == def.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
}
