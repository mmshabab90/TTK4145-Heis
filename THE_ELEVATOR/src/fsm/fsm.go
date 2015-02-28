package fsm

import (
	"../elev"
	"../queue"
	"log"
	"reflect"
	"runtime"
)

var _ = log.Fatal // For debugging only, remove when done
var _ = reflect.ValueOf // For debugging only, remove when done
var _ = runtime.FuncForPC // For debugging only, remove when done

type State_t int
const (
	idle State_t = iota
	moving
	doorOpen
)

var		state			State_t
var		direction		Elev_motor_direction_t
var		floor			int
var		departDirection	Elev_motor_direction_t
const 	doorOpenTime	= 3.0

func syncLights() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if (b == BUTTON_CALL_UP && f == N_FLOORS-1) ||
			(b == BUTTON_CALL_DOWN && f == 0) {
			   	continue
			} else {
				Elev_set_button_lamp(b, f, Orders_isOrder(f, b))
			}
		}
	}
}

func Init() {
	state = idle
	direction = DIRN_STOP
	floor = -1
	departDirection = DIRN_DOWN
	Orders_removeAll()
}

func EventButtonPressed(buttonFloor int, buttonType Elev_button_type_t) {
	switch state {
		case idle:
			queue.AddOrder(buttonFloor, buttonType)
			direction = queue.ChooseDirection(floor, direction)
			if direction == DIRN_STOP {
				driver.SetDoorOpenLamp(true)
				queue.RemoveOrdersAt(floor)
				// timer.Start(doorOpenTime)
				state = doorOpen
			} else {
				driver.SetMotorDirection(direction)
				state = moving
				departDirection = direction
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
			// log error invalid state
	}
	syncLights()
}

func EventArrivedAtFloor(newFloor int) {
	floor = newFloor
	driver.SetFloorIndicator(floor)
	switch state {
	case moving:
		if queue.ShouldStop(floor, direction) {
			driver.SetMotorDirection(DIRN_STOP)
			driver.SetDoorOpenLamp(true)
			queue.RemoveOrdersAt(floor)
			// timer.Start(doorOpenTime)
			syncLights()
			state = doorOpen
		} else {
			departDirection = direction
		}
	default:
		// log error makes no sense to arrive at floor in state <state>
	}
}

func EventTimerTimeOut() {
	switch state {
	case doorOpen:
		direction = queue.ChooseDirection(floor, direction)
		driver.SetDoorOpenLamp(false)
		driver.SetMotorDirection(direction)
		if direction == DIRN_STOP {
			state = idle
		} else {
			state = moving
			departDirection = direction
		}
	default:
		// makes no sense
	}
}

func main() {}