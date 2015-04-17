package fsm

import (
	"../hw"
	"../queue"
	"../network"
	"fmt"
	"log"
	"time"
)

type stateType int
const (
	idle stateType = iota
	moving
	doorOpen
)

var state stateType
var floor int
var direction hw.DirnType
var departDirection hw.DirnType

var DoorTimeout = make(chan bool)
var doorReset = make(chan bool)
const doorOpenTime = 3 * time.Second

func Init() {
	log.Println("FSM Init")
	queue.Init()
	runTimer()
	state = idle
	direction = hw.DirnStop
	floor = hw.GetFloor()
	departDirection = hw.DirnDown
	syncLights()
}

func runTimer() {
	timer := time.NewTimer(0)
	timer.Stop()

	go func() {
		for {
			select {
			case <-doorReset:
				timer.Reset(doorOpenTime)
			case <-timer.C:
				DoorTimeout <- true
				timer.Stop()
			}
		}
	}()
}

func EventButtonPressed(buttonFloor int, buttonType int) {
	fmt.Print("Event button pressed in state ")
	switch state {
	case idle:
		fmt.Println("idle")
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
		fmt.Println("door open")
		if floor == buttonFloor {
			doorReset <- true
		} else {
			queue.AddOrder(buttonFloor, buttonType)
		}
	case moving:
		fmt.Println("moving")
		queue.AddOrder(buttonFloor, buttonType)
	default:
		log.Fatalf("State %d is invalid!\n", state)
	}
	syncLights()
}

func EventFloorReached(newFloor int) {
	fmt.Print("Event floor reached in state ")
	floor = newFloor
	hw.SetFloorIndicator(floor)
	switch state {
	case moving:
		fmt.Println("moving")
		if queue.ShouldStop(floor, direction) {
			hw.SetMotorDirection(hw.DirnStop)
			hw.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
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

func EventTimerOut() {
	fmt.Print("Event timer out in state ")
	switch state {
	case doorOpen:
		fmt.Println("door open")
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

func GetDirection() hw.DirnType {
	return direction
}

func GetFloor() int {
	return floor
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
