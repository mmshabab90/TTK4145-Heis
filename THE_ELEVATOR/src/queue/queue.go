package queue

import (
	"../elev"
	"../network"
	"log"
	"encoding/json"
)

var laddr string

const invalidAddr = "0.0.0.0"

type messageType int
const (
	alive messageType = iota
	newOrder
	completeOrder
	cost
)

type sharedOrder struct {
	isOrderActive bool
	elevatorAddr  string
}

type message struct { // maybe should make this public
	kind messageType
	floor int
	button elev.ButtonType
	cost int
}

var localQueue [elev.NumFloors][elev.NumButtons]bool
// Internal orders in shared queue are unused, but present for better indexing:
var sharedQueue [elev.NumFloors][elev.NumButtons]sharedOrder

// --------------- PUBLIC: ---------------
func Init() {
	resetLocalQueue()
	resetSharedQueue()
	// set laddr variable
}

func AddOrder(floor int, button elev.ButtonType) {
	// New AddOrder should:
		// Send message about new order to all lifts
		// Another func should receive replies and
		// assign the order to the lift with the lowest
		// cost (or lowest ip if several lifts have same cost)

	if button == elev.ButtonCommand {
		localQueue[floor][button] = true
	} else {
		msg := message{kind: newOrder, floor: floor, button: button, cost: -1}
		jsonMsg, err := json.Marshal(msg)
		if err {
			// worry
		}
		network.SendMsg(jsonMsg)
	}
}

func ChooseDirection(currFloor int, currDir elev.DirnType) elev.DirnType {
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
		log.Printf("localQueue: ChooseDirection called with invalid direction %d!\n", currDir)
		return elev.DirnStop
	}
}

func ShouldStop(floor int, direction elev.DirnType) bool {
	switch direction {
	case elev.DirnDown:
		return localQueue[floor][elev.ButtonCallDown] ||
			localQueue[floor][elev.ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case elev.DirnUp:
		return localQueue[floor][elev.ButtonCallUp] ||
			localQueue[floor][elev.ButtonCommand] ||
			floor == elev.NumFloors-1 ||
			!isOrdersAbove(floor)
	case elev.DirnStop:
		return localQueue[floor][elev.ButtonCallDown] ||
			localQueue[floor][elev.ButtonCallUp] ||
			localQueue[floor][elev.ButtonCommand]
	default:
		log.Printf("localQueue: ShouldStop called with invalid direction %d!\n", direction)
		return false
	}
}

func RemoveOrdersAt(floor int) {
	for b := 0; b < elev.NumButtons; b++ {
		localQueue[floor][b] = false
	}
}

func IsOrder(floor int, button elev.ButtonType) bool {
	return localQueue[floor][button]
}

// --------------- PRIVATE: ---------------

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func updateLocalQueue() {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			if b != elev.ButtonCommand &&
				sharedQueue[f][b].isOrderActive &&
				sharedQueue[f][b].elevatorAddr == laddr {
				localQueue[f][b] = true
			}
		}
	}
}

func removeSharedOrder(floor int, button elev.ButtonType) {
	if button == elev.ButtonCommand {
		// error
		return
	}

	sharedQueue[floor][button].isOrderActive = false
	sharedQueue[floor][button].elevatorAddr = invalidAddr
}

func resetLocalQueue() {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func resetSharedQueue() {
	blankOrder := sharedOrder{isOrderActive: false, elevatorAddr: invalidAddr}
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			sharedQueue[f][b] = blankOrder
		}
	}
}



/*func updateSharedQueue(floor int, button elev.ButtonType) {
	// If order completed was assigned to this elevator: Remove from shared queue
	if button == elev.ButtonCommand {
		// error
		return
	}

	if sharedQueue[floor][button].isOrderActive
	&& sharedQueue[floor][button].elevatorAddr == laddr {
		sharedQueue[floor][button].isOrderActive = false
		sharedQueue[floor][button].elevatorAddr = invalidAddr
	}
}*/
