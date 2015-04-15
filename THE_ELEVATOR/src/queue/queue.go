package localQueue

import (
	"../elev"
	"log"
)

var laddr string
const invalidAddr = "0.0.0.0"

// Burde ikke alive-telleren gå i en annen liste over heiser?

type sharedOrder struct {
	isOrderActive bool
	elevatorAddr string
}

var sharedQueue [elev.NumFloors][elev.NumButtons]sharedOrder // internal orders in this are not used
var localQueue [elev.NumFloors][elev.NumButtons]bool

func Init() {
	// set laddr variable
}

func AddOrder(floor int, button elev.ButtonType) {
	localQueue[floor][button] = true
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

func RemoveAll() {
	for f := 0; f < elev.NumFloors; f++ {
		for b := 0; b < elev.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func IsOrder(floor int, button elev.ButtonType) bool {
	return localQueue[floor][button]
}

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
			if b != elev.ButtonCommand
			&& sharedQueue[f][b].isOrderActive
			&& sharedQueue[f][b].elevatorAddr == laddr {
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
