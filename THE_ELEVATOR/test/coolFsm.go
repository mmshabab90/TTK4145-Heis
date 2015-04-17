package fsm

import (
	"../hardware"
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

const doorOpenTime = 3.0

type stateData struct {
	direction Elev_motor_direction_t
	departDirection Elev_motor_direction_t
	floor int

}

var state State_t

func syncLights() {
	for f := 0; f < nFloors; f++ {
		for b := 0; b < nButtons; b++ {
			if (b == ButtonCallUp && f == nFloors-1) ||
				(b == ButtonCallDown && f == 0) {
				continue
			} else {
				hardware.SetButtonLamp(b, f, queue.IsOrder(f, b))
			}
		}
	}
}

func Init() {
	state = idle
	direction = DirnStop
	floor = -1
	departDirection = DirnDown
	queue.RemoveAll()
}

func EventButtonPressed(buttonFloor int, buttonType Elev_button_type_t) {
	switch state {
	case idle:
		queue.AddOrder(buttonFloor, buttonType)
		direction = queue.ChooseDirection(floor, direction)
		if direction == DirnStop {
			driver.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// timer.Start(doorOpenTime)
			state = doorOpen
		} else {
			driver.SetMotorDirection(direction)
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

func EventArrivedAtFloor(newFloor int) {
	floor = newFloor
	driver.SetFloorIndicator(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			driver.SetMotorDirection(DirnStop)
			driver.SetDoorOpenLamp(true)
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

func EventTimerTimeOut() {
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		driver.SetDoorOpenLamp(false)
		driver.SetMotorDirection(direction)
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
