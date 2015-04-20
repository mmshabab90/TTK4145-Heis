package fsm

import (
	"../defs"
	"../hw"
	"../queue"
	"fmt"
	"log"
	"time"
)

type stateType int // kill this

const (
	idle stateType = iota
	moving
	doorOpen
)

var state stateType
var floor int
var direction int
var departDirection int

var doorReset = make(chan bool)

const doorOpenTime = 1 * time.Second

// --------------- PUBLIC: ---------------

var DoorTimeoutChan = make(chan bool)

func Init() {
	log.Println("FSM Init")
	go startTimer()
	state = idle
	direction = defs.DirnStop
	floor = hw.Floor()
	if floor == -1 {
		floor = hw.MoveToDefinedState()
	}
	departDirection = defs.DirnDown
	syncLights()
}

func EventInternalButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("\n\nEvent internal button (floor %d %s) pressed in state %s\n",
		buttonFloor, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle:
		queue.AddInternalOrder(buttonFloor, buttonType)
		switch direction := queue.ChooseDirection(floor, direction); direction {
		case defs.DirnStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		case defs.DirnUp, defs.DirnDown:
			hw.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}		
	case doorOpen:
		if floor == buttonFloor {
			doorReset <- true
		} else {
			queue.AddInternalOrder(buttonFloor, buttonType)
		}
	case moving:
		queue.AddInternalOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}
	syncLights()
}

func EventExternalButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("\n\nEvent external button (floor %d %s) pressed in state %s\n",
		buttonFloor, buttonString(buttonType), stateString(state))
	queue.Print()
	switch state {
	case idle, doorOpen, moving:
		// send order on network
		message := &defs.Message{Kind: defs.NewOrder, Floor: buttonFloor, Button: buttonType}
		defs.MessageChan <- *message
	default:
		//
	}
	syncLights()
}

func EventExternalOrderGivenToMe() {
	queue.Print()
	if queue.IsLocalEmpty() {
		// strange
	}
	switch state {
	case idle:
		switch direction := queue.ChooseDirection(floor, direction); direction {
		case defs.DirnStop:
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		case defs.DirnUp, defs.DirnDown:
			hw.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	default:
		fmt.Println("EventExternalOrderGivenToMe(): Not in idle, will ignore.")
	}
	syncLights()
}

func EventFloorReached(newFloor int) {
	fmt.Printf("\n\nEvent floor %d reached in state %s\n", newFloor, stateString(state))
	queue.Print()
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			hw.SetMotorDirection(defs.DirnStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
			departDirection = direction
		}
	default:
		log.Printf("Makes no sense to arrive at a floor in state %s.\n", stateString(state))
	}
	syncLights()
}

func EventDoorTimeout() {
	fmt.Printf("\n\nEvent door timeout in state %s\n", stateString(state))
	// queue.Print()
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(direction)
		if direction == defs.DirnStop {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to time out when not in state door open\n")
	}
	syncLights()
}

func Direction() int {
	return direction
}

func DepartDirection() int {
	return departDirection
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

func syncLights() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if (b == defs.ButtonCallUp && f == defs.NumFloors-1) ||
				(b == defs.ButtonCallDown && f == 0) {
				continue
			} else {
				hw.SetButtonLamp(f, b, queue.IsOrder(f, b))
			}
		}
	}
}

func stateString(state stateType) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case doorOpen:
		return "door open"
	default:
		return "error: bad state"
	}
}

func buttonString(button int) string {
	switch button {
	case defs.ButtonCallUp:
		return "up"
	case defs.ButtonCallDown:
		return "down"
	case defs.ButtonCommand:
		return "command"
	default:
		return "error: bad button"
	}
}
