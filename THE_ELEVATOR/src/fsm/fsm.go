package fsm

import (
	"elev"
	"queue"
	"log"
	"reflect"
	"runtime"
	"fmt"
)

var _ = log.Fatal         // For debugging only, remove when done
var _ = reflect.ValueOf   // For debugging only, remove when done
var _ = runtime.FuncForPC // For debugging only, remove when done

type State_t int

const (
	idle State_t = iota
	moving
	doorOpen
)

var state State_t
var direction elev.MotorDirnType
var floor int
var departDirection elev.MotorDirnType

const doorOpenTime = 3.0

func Init() {
	fmt.Println("Init")
	state = idle
	direction = elev.DirnStop
	floor = elev.GetFloor()
	departDirection = elev.DirnDown
	queue.RemoveAll()
	syncLights()
}

func EventButtonPressed(buttonFloor int, buttonType elev.ButtonType) {
	fmt.Print("Event button pressed in state ")
	switch state {
	case idle:
		fmt.Println("idle")
		queue.AddOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == elev.DirnStop {
			elev.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// timer.Start(doorOpenTime)
			state = doorOpen
		} else {
			elev.SetMotorDirection(direction)
			departDirection = direction
			state = moving
		}
	case doorOpen:
		fmt.Println("door open")
		if floor == buttonFloor {
			// timer.Start(doorOpenTime)
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
	elev.SetFloorIndicator(floor)
	switch state {
	case moving:
		fmt.Println("moving")
		if queue.ShouldStop(floor, direction) {
			elev.SetMotorDirection(elev.DirnStop)
			elev.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// timer.Start(doorOpenTime)
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
		fmt.Println("Washington")
		direction = queue.ChooseDirection(floor, direction)
		elev.SetDoorOpenLamp(false)
		elev.SetMotorDirection(direction)
		if direction == elev.DirnStop {
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

func syncLights() {
	var b elev.ButtonType
	for f := 0; f < elev.NumFloors; f++ {
		for b = 0; b < elev.NumButtons; b++ {
			if (b == elev.ButtonCallUp && f == elev.NumFloors-1) ||
				(b == elev.ButtonCallDown && f == 0) {
				continue
			} else {
				elev.SetButtonLamp(f, b, queue.IsOrder(f, b))
			}
		}
	}
}
