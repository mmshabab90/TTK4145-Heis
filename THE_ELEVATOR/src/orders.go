package main

import (
	"log"
	"./driver"
)

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

var orders = [N_FLOORS][N_BUTTONS] bool {
	{false, false, false},
	{false, false, false},
	{false, false, false},
	{false, false, false},
}

func Orders_addOrder(floor int, button Elev_button_type_t) {
	orders[floor][button] = true;
}

func Orders_chooseDirection(currFloor int, currDir Elev_motor_direction_t) Elev_motor_direction_t {
	if !isAnyOrders() {
		return DIRN_STOP
	}
	switch currDir {
	case DIRN_DOWN:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return DIRN_DOWN
		} else {
			return DIRN_UP
		}
	case DIRN_UP:
		if isOrdersAbove(currFloor) && currFloor < N_FLOORS - 1 {
			return DIRN_UP
		} else {
			return DIRN_DOWN
		}
	case DIRN_STOP:
		if isOrdersAbove(currFloor) {
			return DIRN_UP
		} else if isOrdersBelow(currFloor) {
			return DIRN_DOWN
		} else {
			return DIRN_STOP
		}
	default:
		log.Printf("Orders_chooseDirection called with unexpected direction %d!", currDir)
		return DIRN_STOP
	}
}

func Orders_shouldStop(floor int, direction Elev_motor_direction_t) bool {
	switch direction {
	case DIRN_DOWN:
		return	orders[floor][BUTTON_CALL_DOWN] ||
				orders[floor][BUTTON_COMMAND] 	||
				floor == 0						||
				!isOrdersBelow(floor)
	case DIRN_UP:
		return	orders[floor][BUTTON_CALL_UP]	||
				orders[floor][BUTTON_COMMAND]	||
				floor == N_FLOORS - 1			||
				!isOrdersAbove(floor)
	case DIRN_STOP:
		return	orders[floor][BUTTON_CALL_DOWN]	||
				orders[floor][BUTTON_CALL_UP]	||
				orders[floor][BUTTON_COMMAND]
	default:
		log.Printf("Orders_shouldStop called with unexpected direction %d!", direction)
		return false
	}
}

func Orders_removeOrder(floor int) {
	for b := 0; b < N_BUTTONS; b++ {
		orders[floor][b] = false;
	}
}

func Orders_removeAll() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			orders[f][b] = false;
		}
	}
}

func Orders_isOrder(floor int, button Elev_button_type_t) bool {
	return orders[floor][button];
}

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			if orders[f][b] {
				return true
			}
		}
	}
	return false
}


func main() {
}