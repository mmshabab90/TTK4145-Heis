package fsm

import (
	"../elev"
	"../queue"
	"log"
	"reflect"
	"runtime"
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
var direction Elev_motor_direction_t
var floor int
var departDirection Elev_motor_direction_t

const doorOpenTime = 3.0

func Init() {
	state = idle
	direction = DirnStop
	floor = elev.GetFloor()
	departDirection = DirnDown
	queue.RemoveAll()
}

func EventButtonPressed(buttonFloor int, buttonType Elev_button_type_t) {
	switch state {
	case idle:
		queue.AddOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == DirnStop {
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
		if floor == buttonFloor {
			// timer.Start(doorOpenTime)
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
	floor = newFloor
	elev.SetFloorIndicator(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			elev.SetMotorDirection(DirnStop)
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
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		elev.SetDoorOpenLamp(false)
		elev.SetMotorDirection(direction)
		if direction == DirnStop {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		log.Fatalf("Makes no sense to time out when not in doorOpen\n")
	}
}

func syncLights() {
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			if (b == ButtonCallUp && f == NumFloors-1) ||
				(b == ButtonCallDown && f == 0) {
				continue
			} else {
				elev.SetButtonLamp(f, b, queue.IsOrder(f, b))
			}
		}
	}
}
