package queue

import (
	"../hw"
	"../network"
	"log"
	"encoding/json"
	"fmt"
)

var laddr string

const invalidAddr = "0.0.0.0"

type sharedOrder struct {
	isOrderActive bool
	elevatorAddr  string
}

var localQueue [hw.NumFloors][hw.NumButtons]bool
// Internal orders in shared queue are unused, but present for better indexing:
var sharedQueue [hw.NumFloors][hw.NumButtons]sharedOrder

// --------------- PUBLIC: ---------------

func Init() {
	resetLocalQueue()
	resetSharedQueue()
	// set laddr variable
}

func AddOrder(floor int, button hw.ButtonType) {
	// New AddOrder should:
		// Send message about new order to all lifts
		// Another func should receive replies and
		// assign the order to the lift with the lowest
		// cost (or lowest ip if several lifts have same cost)

	if button == hw.ButtonCommand {
		localQueue[floor][button] = true
	} else {
		message := network.Message{kind: newOrder, floor: floor, button: button}
		network.Send(message)
	}
}

func ChooseDirection(currFloor int, currDir hw.DirnType) hw.DirnType {
	if !isAnyOrders() {
		return hw.DirnStop
	}
	switch currDir {
	case hw.DirnDown:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return hw.DirnDown
		} else {
			return hw.DirnUp
		}
	case hw.DirnUp:
		if isOrdersAbove(currFloor) && currFloor < hw.NumFloors-1 {
			return hw.DirnUp
		} else {
			return hw.DirnDown
		}
	case hw.DirnStop:
		if isOrdersAbove(currFloor) {
			return hw.DirnUp
		} else if isOrdersBelow(currFloor) {
			return hw.DirnDown
		} else {
			return hw.DirnStop
		}
	default:
		log.Printf("localQueue: ChooseDirection called with invalid direction %d!\n", currDir)
		return hw.DirnStop
	}
}

func ShouldStop(floor int, direction hw.DirnType) bool {
	switch direction {
	case hw.DirnDown:
		return localQueue[floor][hw.ButtonCallDown] ||
			localQueue[floor][hw.ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case hw.DirnUp:
		return localQueue[floor][hw.ButtonCallUp] ||
			localQueue[floor][hw.ButtonCommand] ||
			floor == hw.NumFloors-1 ||
			!isOrdersAbove(floor)
	case hw.DirnStop:
		return localQueue[floor][hw.ButtonCallDown] ||
			localQueue[floor][hw.ButtonCallUp] ||
			localQueue[floor][hw.ButtonCommand]
	default:
		log.Printf("localQueue: ShouldStop called with invalid direction %d!\n", direction)
		return false
	}
}

func RemoveOrdersAt(floor int) {
	for b := 0; b < hw.NumButtons; b++ {
		localQueue[floor][b] = false
	}
}

func IsOrder(floor int, button hw.ButtonType) bool {
	return localQueue[floor][button]
}

// --------------- PRIVATE: ---------------

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func updateLocalQueue() {
	var b hw.ButtonType
	for f := 0; f < hw.NumFloors; f++ {
		for b = 0; b < hw.NumButtons; b++ {
			if b != hw.ButtonCommand &&
				sharedQueue[f][b].isOrderActive &&
				sharedQueue[f][b].elevatorAddr == laddr {
				localQueue[f][b] = true
			}
		}
	}
}

func RemoveSharedOrder(floor int, button hw.ButtonType) {
	if button == hw.ButtonCommand {
		// error
		return
	}

	sharedQueue[floor][button].isOrderActive = false
	sharedQueue[floor][button].elevatorAddr = invalidAddr
}

func resetLocalQueue() {
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func resetSharedQueue() {
	blankOrder := sharedOrder{isOrderActive: false, elevatorAddr: invalidAddr}
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			sharedQueue[f][b] = blankOrder
		}
	}
}



/*func updateSharedQueue(floor int, button hw.ButtonType) {
	// If order completed was assigned to this elevator: Remove from shared queue
	if button == hw.ButtonCommand {
		// error
		return
	}

	if sharedQueue[floor][button].isOrderActive
	&& sharedQueue[floor][button].elevatorAddr == laddr {
		sharedQueue[floor][button].isOrderActive = false
		sharedQueue[floor][button].elevatorAddr = invalidAddr
	}
}*/
