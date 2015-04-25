// Package fsm implements a finite state machine for the behaviour of a lift.
// The lift runs based on a queue stored and managed by the queue package.
//
// There are three states:
// Idle: Lift is stationary, at a floor, door closed, awaiting orders.
// Moving: Lift is moving, can be between floors or at a floor going past it.
// Door open: Lift is at a floor with the door open.
//
// And three events:
// New order: A new order is added to the queue.
// Floor reached: The lift reaches a floor.
// Door timeout: The door timer times out (the door should close).
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

	log.Println(def.ClrG, "FSM initialised.", def.ClrN)
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
	log.Printf("%sEVENT: New order in state %v.%s", def.ClrY, stateString(state), def.ClrN)
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
		// Ignore.
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
	log.Printf("%sEVENT: Floor %d reached in state %s.%s", def.ClrY, newFloor+1, stateString(state), def.ClrN)
	// queue.Print()
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
	log.Printf("%sEVENT: Door timeout in state %s.%s", def.ClrY, stateString(state), def.ClrN)
	// queue.Print()
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
