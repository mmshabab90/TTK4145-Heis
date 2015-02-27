package main

import (
	"log"
	"reflect"
	"runtime"
)

var _ = log.Fatal // For debugging only, remove when done
var _ = reflect.ValueOf
var _ = runtime.FuncForPC

// THESE ARE IN THE DRIVER; REMOVE FROM THIS FILE
// AFTER DRIVER IS IMPORTED
const N_BUTTONS = 3
const N_FLOORS  = 4
type Elev_button_type_t int
type Elev_motor_direction_t int
const (
	DIRN_DOWN Elev_motor_direction_t = -1
	DIRN_STOP Elev_motor_direction_t  = 0
	DIRN_UP Elev_motor_direction_t = 1
)
const (
	BUTTON_CALL_UP Elev_button_type_t = iota
	BUTTON_CALL_DOWN
	BUTTON_COMMAND
)
/////////////////////////////////////////////////

type State_t int
const (
	INIT State_t = iota
	IDLE
	DOOROPEN
	MOVING
	EMERGENCY
	STOPPEDBETWEENFLOORS
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

func Fsm_init() {
	state = INIT
	direction = DIRN_STOP
	floor = -1
	departDirection = DIRN_DOWN
	Orders_removeAll()
}

func Fsm_eventButtonPressed(floor int, button Elev_button_type_t) {

}

func main() {}