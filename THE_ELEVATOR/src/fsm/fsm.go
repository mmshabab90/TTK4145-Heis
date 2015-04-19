package fsm

import (
	"../defs"
	"../hw"
	"../queue"
	"fmt"
	"log"
	"time"
)

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

func EventButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("Event button (floor %d %s) pressed in state %s\n", buttonFloor, buttonName(buttonType), stateName(state))
	switch state {
	case idle:
		queue.NewOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == defs.DirnStop {
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			doorReset <- true
			state = doorOpen
		} else {
			hw.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	case doorOpen:
		if floor == buttonFloor {
			doorReset <- true
		} else {
			queue.NewOrder(buttonFloor, buttonType)
		}
	case moving:
		queue.NewOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}
	syncLights()
}

func EventFloorReached(newFloor int) {
	fmt.Printf("Event floor %d reached in state %s\n", newFloor, stateName(state))
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			hw.SetMotorDirection(defs.DirnStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// send completed order-message:
			go queue.SendOrderCompleteMessage(floor)
			doorReset <- true
			state = doorOpen
		} else {
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to arrive at a floor in state %d", state)
	}
	syncLights()
}

func EventDoorTimeout() {
	fmt.Printf("Event door timeout in state %s\n", stateName(state))
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

func Floor() int {
	return floor
}

// --------------- PRIVATE: ---------------

type stateType int

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

const doorOpenTime = 3 * time.Second

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

func stateName(state stateType) string {
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

func buttonName(button int) string {
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
