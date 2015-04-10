package queue

import (
	"../elev"
	"log"
)

var queue [elev.NumFloors][elev.NumButtons]bool

func AddOrder(floor int, button elev.ButtonType) {
	queue[floor][button] = true
}

func ChooseDirection(currFloor int, currDir elev.MotorDirnType) elev.MotorDirnType {
	if !isAnyOrders() {
		return elev.DirnStop
	}
	switch currDir {
	case elev.DirnDown:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return elev.DirnDown
		} else {
			return elev.DirnUp
		}
	case elev.DirnUp:
		if isOrdersAbove(currFloor) && currFloor < elev.NumFloors-1 {
			return elev.DirnUp
		} else {
			return elev.DirnDown
		}
	case elev.DirnStop:
		if isOrdersAbove(currFloor) {
			return elev.DirnUp
		} else if isOrdersBelow(currFloor) {
			return elev.DirnDown
		} else {
			return elev.DirnStop
		}
	default:
		log.Printf("queue: ChooseDirection called with invalid direction %d!\n", currDir)
		return elev.DirnStop
	}
}

func ShouldStop(floor int, direction elev.MotorDirnType) bool {
	switch direction {
	case elev.DirnDown:
		return queue[floor][elev.ButtonCallDown] ||
			queue[floor][elev.ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case elev.DirnUp:
		return queue[floor][elev.ButtonCallUp] ||
			queue[floor][elev.ButtonCommand] ||
			floor == elev.NumFloors-1 ||
			!isOrdersAbove(floor)
	case elev.DirnStop:
		return queue[floor][elev.ButtonCallDown] ||
			queue[floor][elev.ButtonCallUp] ||
			queue[floor][elev.ButtonCommand]
	default:
		log.Printf("queue: ShouldStop called with invalid direction %d!\n", direction)
		return false
	}
}

func RemoveOrdersAt(floor int) {
	for b := 0; b < elev.NumButtons; b++ {
		queue[floor][b] = false
	}
}

func RemoveAll() {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			queue[f][b] = false
		}
	}
}

func IsOrder(floor int, button elev.ButtonType) bool {
	return queue[floor][button]
}

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if queue[f][b] {
				return true
			}
		}
	}
	return false
}
