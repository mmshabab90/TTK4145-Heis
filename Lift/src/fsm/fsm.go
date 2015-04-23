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
var storey int
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
	storey = hw.Storey()
	if storey == -1 {
		storey = hw.MoveToDefinedState()
	}
	go syncLights()
}

func EventInternalButtonPressed(buttonStorey int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event internal button (storey %d %s) pressed in state %s\n",
		buttonStorey, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle:
		queue.AddLocalOrder(buttonStorey, buttonType)
		switch direction = queue.ChooseDirection(storey, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(storey)
			doorReset <- true
			state = doorOpen
		case defs.DirUp, defs.DirDown:
			hw.SetMotorDirection(direction)
			state = moving
		}
	case doorOpen:
		if storey == buttonStorey {
			doorReset <- true
		} else {
			queue.AddLocalOrder(buttonStorey, buttonType)
		}
	case moving:
		queue.AddLocalOrder(buttonStorey, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}

	defs.SyncLightsChan <- true
}

func EventExternalButtonPressed(buttonStorey int, buttonType int) {
	fmt.Printf("\n\n   ☺      Event external button (storey %d %s) pressed in state %s\n",
		buttonStorey, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle, doorOpen, moving:
		// send order on network
		message := defs.Message{Kind: defs.NewOrder, Storey: buttonStorey, Button: buttonType, Cost: -1}
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
		switch direction = queue.ChooseDirection(storey, direction); direction {
		case defs.DirStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(storey) //her tror jeg buggen ligger!
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

func EventStoreyReached(newStorey int) {
	fmt.Printf("\n\n   ☺      Event storey %d reached in state %s\n", newStorey, stateString(state))
	queue.Print()
	storey = newStorey
	hw.SetStoreyLamp(storey)
	switch state {
	case moving:
		if queue.ShouldStop(storey, direction) {
			hw.SetMotorDirection(defs.DirStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(storey)
			go queue.SendOrderCompleteMessage(storey)
			doorReset <- true
			state = doorOpen
		} else {
		}
	default:
		log.Printf("Makes no sense to arrive at a storey in state %s.\n", stateString(state))
	}
	defs.SyncLightsChan <- true
}

func EventDoorTimeout() { //this happens for each external order
	fmt.Printf("\n\n   ☺      Event door timeout in state %s\n", stateString(state))
	queue.Print()
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(storey, direction)
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

func Storey() int {
	return storey
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
