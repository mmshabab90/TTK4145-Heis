package queue

import (
	"../elev"
	"log"
)

var queue = [NumFloors][NumButtons]bool{
	{false, false, false},
	{false, false, false},
	{false, false, false},
	{false, false, false},
}

func AddOrder(floor int, button Elev_button_type_t) {
	queue[floor][button] = true
}

func ChooseDirection(currFloor int, currDir Elev_motor_direction_t) Elev_motor_direction_t {
	if !isAnyOrders() {
		return DirnStop
	}
	switch currDir {
	case DirnDown:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return DirnDown
		} else {
			return DirnUp
		}
	case DirnUp:
		if isOrdersAbove(currFloor) && currFloor < NumFloors-1 {
			return DirnUp
		} else {
			return DirnDown
		}
	case DirnStop:
		if isOrdersAbove(currFloor) {
			return DirnUp
		} else if isOrdersBelow(currFloor) {
			return DirnDown
		} else {
			return DirnStop
		}
	default:
		log.Printf("Orders_chooseDirection called with unexpected direction %d!", currDir)
		return DirnStop
	}
}

func ShouldStop(floor int, direction Elev_motor_direction_t) bool {
	switch direction {
	case DirnDown:
		return queue[floor][ButtonCallDown] ||
			queue[floor][ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case DirnUp:
		return queue[floor][ButtonCallUp] ||
			queue[floor][ButtonCommand] ||
			floor == NumFloors-1 ||
			!isOrdersAbove(floor)
	case DirnStop:
		return queue[floor][ButtonCallDown] ||
			queue[floor][ButtonCallUp] ||
			queue[floor][ButtonCommand]
	default:
		log.Printf("Orders_shouldStop called with unexpected direction %d!", direction)
		return false
	}
}

func RemoveOrdersAt(floor int) {
	for b := 0; b < NumButtons; b++ {
		queue[floor][b] = false
	}
}

func RemoveAll() {
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			queue[f][b] = false
		}
	}
}

func IsOrder(floor int, button Elev_button_type_t) bool {
	return queue[floor][button]
}

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < NumFloors; f++ {
		for b := 0; b < NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}

func main() {
}
