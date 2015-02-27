package main

import (
	"./driver"
	"log"
)

var orders = [N_FLOORS][N_BUTTONS] bool {
	{false, false, false},
	{false, false, false},
	{false, false, false},
	{false, false, false},
}

func AddOrder(floor int, button Elev_button_type_t) {
	orders[floor][button] = true;
}

func ChooseDirection(currFloor int, currDir Elev_motor_direction_t) Elev_motor_direction_t {
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

func ShouldStop(floor int, direction Elev_motor_direction_t) bool {
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

func RemoveOrder(floor int) {
	for b := 0; b < N_BUTTONS; b++ {
		orders[floor][b] = false;
	}
}

func RemoveAll() {
	for f := 0; f < N_FLOORS; f++ {
		for b := 0; b < N_BUTTONS; b++ {
			orders[f][b] = false;
		}
	}
}

func IsOrder(floor int, button Elev_button_type_t) bool {
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