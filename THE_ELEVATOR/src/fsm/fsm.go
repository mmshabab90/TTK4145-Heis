package fsm

import (
	"../hw"
	"../queue"
	"fmt"
	"log"
	"time"
	"../defs"
)

// --------------- PUBLIC: ---------------

func Init() {
	log.Println("FSM Init")
	queue.Init()
	runTimer()
	state = idle
	direction = hw.DirnStop
	floor = hw.GetFloor()
	if floor == -1 {
		floor = hw.MoveToDefinedState()
	}
	departDirection = hw.DirnDown
	syncLights()
}

func EventButtonPressed(buttonFloor int, buttonType int) {
	fmt.Printf("Event button pressed in state %s\n", getStateName(state))
	switch state {
	case idle:
		queue.AddOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == hw.DirnStop {
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
			queue.AddOrder(buttonFloor, buttonType)
		}
	case moving:
		queue.AddOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}
	syncLights()
}

func EventFloorReached(newFloor int) {
	fmt.Printf("Event button pressed in state %s\n", getStateName(state))
	floor = newFloor
	hw.SetFloorLamp(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			hw.SetMotorDirection(hw.DirnStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// send completed order-message:
			queue.SendOrderCompleteMessage(floor)
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
	fmt.Printf("Event button pressed in state %s\n", getStateName(state))
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		hw.SetDoorOpenLamp(false)
		hw.SetMotorDirection(direction)
		if direction == hw.DirnStop {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to time out when not in doorOpen\n")
	}
	syncLights()
}

func GetDirection() int {
	return direction
}

func GetFloor() int {
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
var doorTimeout = make(chan bool)

const doorOpenTime = 3 * time.Second

func runTimer() {
	timer := time.NewTimer(0)
	timer.Stop()

	go func() {
		for {
			select {
			case <-doorReset:
				timer.Reset(doorOpenTime)
			case <-timer.C:
				doorTimeout <- true
				timer.Stop()
			}
		}
	}()
}

func syncLights() {
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if (b == hw.ButtonCallUp && f == hw.NumFloors-1) ||
				(b == hw.ButtonCallDown && f == 0) {
				continue
			} else {
				hw.SetButtonLamp(f, b, queue.IsOrder(f, b))
			}
		}
	}
}

func getStateName(state stateType) string {
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
