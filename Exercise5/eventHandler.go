package driver
// check package naming!

import (
	"log"
)

type State_t int

type elevator_state struct (
	floor int
	currentDirection Elev_motor_direction_t
)

// const (
// 	IDLE State_t = iota
// 	RUNNING
// 	OPENDOOR
// 	STOPPED
// )

func idle() State_t {

}

func fsm_run() {
	for state := idle; state != nil; {
		state = state()
	}
}

func 

func EventHandler(){
	for {
		if Elev_get_stop_signal() {
			// Worry about stop signal.
		}

		for button := 0; button < N_BUTTONS; button++ {
			for floor := 0; floor < N_FLOORS; floor++ {
				if isValidOrderButton(button, floor) {
					// Worry about button.
				} else {
					continue
				}
			}
		}

		if Elev_get_floor_sensor_signal() != -1 {
			// Worry about reaching a floor.
		}


	}
}

func init(){
	Elev_init()
	Elev_set_speed(-300)
	for{
		floor := Elev_get_floor_sensor_signal()
		if floor != -1{
		
		}
	}
}

func isValidOrderButton(button Elev_button_type_t, floor int) {
	if floor >= N_FLOORS {
		return false
	} else if floor < 0 {
		return false
	} else if floor == 0 && button == BUTTON_CALL_DOWN {
		return false
	} else if floor == N_FLOORS - 1 && button == BUTTON_CALL_UP {
		return false
	} else if button != BUTTON_CALL_UP
	&& button != BUTTON_CALL_DOWN
	&& button != BUTTON_COMMAND {
		return false
	}
	return true
}