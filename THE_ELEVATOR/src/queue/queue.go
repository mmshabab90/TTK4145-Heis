package queue

import (
	"../hw"
	"../network"
	"log"
)

type sharedOrder struct {
	isOrderActive bool
	assignedLiftAddr  string
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

func AddOrder(floor int, button int) {
	// New AddOrder should:
		// Send message about new order to all lifts
		// Another func should receive replies and
		// assign the order to the lift with the lowest
		// cost (or lowest ip if several lifts have same cost)

	if button == hw.ButtonCommand {
		localQueue[floor][button] = true
	} else {
		message := network.Message{Kind: network.NewOrder, Floor: floor, Button: button}
		network.Send(message)
	}
}

func ChooseDirection(currFloor int, currDir int) int {
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

func ShouldStop(floor int, direction int) bool {
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

func IsOrder(floor int, button int) bool {
	return localQueue[floor][button]
}

func ReassignOrders(deadAddr string) { // better name plz
	// loop thru shared queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if sharedQueue[f][b].assignedLiftAddr == deadAddr {
				sharedQueue[f][b].isOrderActive = false
				sharedQueue[f][b].assignedLiftAddr = ""

                reassignMessage := network.Message{
                	Kind: network.NewOrder,
                	Floor: f,
                	Button: b}
				network.Send(reassignMessage)
			}
		}
	}
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
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			if b != hw.ButtonCommand &&
				sharedQueue[f][b].isOrderActive &&
				sharedQueue[f][b].assignedLiftAddr == network.Laddr.String() {
				localQueue[f][b] = true
			}
		}
	}
}

func RemoveSharedOrder(floor int, button int) {
	if button == hw.ButtonCommand {
		// error
		return
	}

	sharedQueue[floor][button].isOrderActive = false
	sharedQueue[floor][button].assignedLiftAddr = ""
}

func resetLocalQueue() {
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func resetSharedQueue() {
	blankOrder := sharedOrder{isOrderActive: false, assignedLiftAddr: ""}
	for f := 0; f < hw.NumFloors; f++ {
		for b := 0; b < hw.NumButtons; b++ {
			sharedQueue[f][b] = blankOrder
		}
	}
}

/*func updateSharedQueue(floor int, button int) {
	// If order completed was assigned to this elevator: Remove from shared queue
	if button == hw.ButtonCommand {
		// error
		return
	}

	if sharedQueue[floor][button].isOrderActive
	&& sharedQueue[floor][button].assignedLiftAddr == laddr {
		sharedQueue[floor][button].isOrderActive = false
		sharedQueue[floor][button].assignedLiftAddr = invalidAddr
	}
}*/
