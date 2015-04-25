// Package fsm implements a finite state machine for the behaviour of a lift.
// The lift runs based on a queue stored and managed by the queue package.
package fsm

import (
	def "config"
	"fmt"
	"log"
	"queue"
	"time"
)

const (
	idle int = iota
	moving
	doorOpen
)

const doorOpenTime = 1 * time.Second // todo move to doorTimer.go

var state int
var floor int
var dir int

type Channels struct {
	// Events
	NewOrder     chan bool
	FloorReached chan int
	DoorTimeout  chan bool
	// Hardware interation
	MotorDir  chan int
	FloorLamp chan int
	DoorLamp  chan bool
}

var doorReset = make(chan bool)

func Init(e Channels, startFloor int) {
	state = idle
	dir = def.DirStop
	floor = startFloor

	go syncLights()

	e.DoorTimeout = make(chan bool)
	go startDoorTimer(e.DoorTimeout)
	go run(e)

	log.Println("fsm initialized")
}

func run(e Channels) {
	for {
		select {
		case <-e.NewOrder:
			eventNewOrder(e)
		case f := <-e.FloorReached:
			eventFloorReached(e, f)
		case <-e.DoorTimeout:
			eventDoorTimeout(e)
		}
	}
}

func eventNewOrder(e Channels) {
	fmt.Printf("\nEVENT: New order in state %v.\n\n", stateString(state))
	switch state {
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if queue.ShouldStop(floor, dir) {
			// hw.SetDoorOpenLamp(true)
			e.DoorLamp <- true
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
			// hw.SetMotorDirection(dir)
			e.MotorDir <- dir
			state = moving
		}
	case moving:
		// ignore
	case doorOpen:
		if queue.ShouldStop(floor, dir) {
			queue.RemoveOrdersAt(floor)
			doorReset <- true
		}
	default:
		//komenter inn dette når vi vil kjøre restarts
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("This state doesn't exist")
	}
}

func eventFloorReached(e Channels, newFloor int) {
	fmt.Printf("\nEVENT: Floor %d reached in state %s.\n\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	// hw.SetFloorLamp(floor)
	e.FloorLamp <- floor
	switch state {
	case moving:
		if queue.ShouldStop(floor, dir) {
			dir = def.DirStop
			// hw.SetMotorDirection(dir)
			e.MotorDir <- dir
			// hw.SetDoorOpenLamp(true)
			e.DoorLamp <- true
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
		}
	default:
		//komenter inn dette når vi vil kjøre restarts
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("Makes no sense to arrive at a floor in state %s.\n", stateString(state))
	}
}

func eventDoorTimeout(e Channels) {
	fmt.Printf("\nEVENT: Door timeout in state %s.\n\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		dir = queue.ChooseDirection(floor, dir)
		// hw.SetDoorOpenLamp(false)
		e.DoorLamp <- false
		// hw.SetMotorDirection(dir)
		e.MotorDir <- dir
		if dir == def.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		//komenter inn dette når vi vil kjøre restarts
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
}
