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
var direction int

type Events struct {
	NewOrder     <-chan bool
	FloorReached <-chan int
	DoorTimeout  <-chan bool
}

var doorReset = make(chan bool)
var DoorTimeoutChan = make(chan bool)

func Init() {
	log.Println("fsm.Init() starting")
	go startTimer()
	state = idle
	direction = defs.DirStop
	floor = hw.Floor()
	if floor == -1 {
		floor = hw.MoveToDefinedState()
	}
	go syncLights()
}

func EventInternalButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event internal button (floor %d %s) pressed in state %s\n",
		buttonFloor, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle:
		queue.AddLocalOrder(buttonFloor, buttonType)
		switch direction = queue.ChooseDirection(floor, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		case defs.DirUp, defs.DirDown:
			hw.SetMotorDirection(direction)
			state = moving
		}
	case doorOpen:
		if floor == buttonFloor {
			doorReset <- true
		} else {
			queue.AddLocalOrder(buttonFloor, buttonType)
		}
	case moving:
		queue.AddLocalOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}

	defs.SyncLightsChan <- true
}

func EventExternalButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event external button (floor %d %s) pressed in state %s\n",
		buttonFloor, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle, doorOpen, moving:
		// send order on network
		message := defs.Message{Kind: defs.NewOrder, Floor: buttonFloor, Button: buttonType, Cost: -1}
		defs.MessageChan <- message
	default:
		//
	}

	defs.SyncLightsChan <- true
}

func EventExternalOrderGivenToMe() {
	fmt.Printf("\n\n   ☺      Event external order given to me.\n")
	queue.Print()

	if queue.IsLocalEmpty() {
		// strange
	}
	switch state {
	case idle:
		switch direction = queue.ChooseDirection(floor, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor) //her tror jeg buggen ligger!
			doorReset <- true
			state = doorOpen
		case defs.DirUp, defs.DirDown:
			hw.SetMotorDirection(direction)
			state = moving
		}
	default:
		fmt.Println("   ☺      EventExternalOrderGivenToMe(): Not in idle, will ignore.")
	}
	defs.SyncLightsChan <- true
}

func EventNewOrder(floor, button int) {
	switch state {
	case idle:
		dir = queue.ChooseDirection(floor, dir)
		if dir == defs.DirDown {
			// todo more here
		} else {
			// todo more here
		}
	}
}

func EventFloorReached(newFloor int) {
	fmt.Printf("\n\n   ☺      Event floor %d reached in state %s\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
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

func EventDoorTimeout() { //this happens for each external order
	fmt.Printf("\n\n   ☺      Event door timeout in state %s\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(direction)
		if direction == defs.DirStop {
			state = idle
		} else {
			state = moving
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
	defs.SyncLightsChan <- true
}

func Direction() int {
	return direction
}

func Floor() int {
	return floor
}

func startTimer() {
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-doorReset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			EventDoorTimeout()
		}
	}
}
