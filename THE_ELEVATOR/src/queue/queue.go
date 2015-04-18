package queue

import (
	"../defs"
	"log"
)

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

	if button == defs.ButtonCommand {
		localQueue[floor][button] = true
	} else {
		message := &defs.Message{Kind: defs.NewOrder, Floor: floor, Button: button}
		defs.MessageChan <- *message
		//network.Send(message)
	}
}

func ChooseDirection(currFloor int, currDir int) int {
	if !isAnyOrders() {
		return defs.DirnStop
	}
	switch currDir {
	case defs.DirnDown:
		if isOrdersBelow(currFloor) && currFloor > 0 {
			return defs.DirnDown
		} else {
			return defs.DirnUp
		}
	case defs.DirnUp:
		if isOrdersAbove(currFloor) && currFloor < defs.NumFloors-1 {
			return defs.DirnUp
		} else {
			return defs.DirnDown
		}
	case defs.DirnStop:
		if isOrdersAbove(currFloor) {
			return defs.DirnUp
		} else if isOrdersBelow(currFloor) {
			return defs.DirnDown
		} else {
			return defs.DirnStop
		}
	default:
		log.Printf("localQueue: ChooseDirection called with invalid direction %d!\n", currDir)
		return defs.DirnStop
	}
}

func ShouldStop(floor int, direction int) bool {
	switch direction {
	case defs.DirnDown:
		return localQueue[floor][defs.ButtonCallDown] ||
			localQueue[floor][defs.ButtonCommand] ||
			floor == 0 ||
			!isOrdersBelow(floor)
	case defs.DirnUp:
		return localQueue[floor][defs.ButtonCallUp] ||
			localQueue[floor][defs.ButtonCommand] ||
			floor == defs.NumFloors-1 ||
			!isOrdersAbove(floor)
	case defs.DirnStop:
		return localQueue[floor][defs.ButtonCallDown] ||
			localQueue[floor][defs.ButtonCallUp] ||
			localQueue[floor][defs.ButtonCommand]
	default:
		log.Printf("localQueue: ShouldStop called with invalid direction %d!\n", direction)
		return false
	}
}

func RemoveOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
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
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if sharedQueue[f][b].assignedLiftAddr == deadAddr {
				sharedQueue[f][b] = blankOrder
				reassignMessage := &defs.Message{
					Kind:   defs.NewOrder,
					Floor:  f,
					Button: b}
				//network.Send(reassignMessage)
				defs.MessageChan <- *reassignMessage
			}
		}
	}
}

func SendOrderCompleteMessage(floor int) {
	message := &defs.Message{Kind: defs.CompleteOrder, Floor: floor}
	//network.Send(message)
	defs.MessageChan <- *message
}

// --------------- PRIVATE: ---------------

var blankOrder = sharedOrder{isOrderActive: false, assignedLiftAddr: ""}

type sharedOrder struct {
	isOrderActive    bool
	assignedLiftAddr string
}

var localQueue [defs.NumFloors][defs.NumButtons]bool

// Internal orders in shared queue are unused, but present for better indexing:
var sharedQueue [defs.NumFloors][defs.NumButtons]sharedOrder

func isOrdersAbove(floor int) bool {
	for f := floor + 1; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func isAnyOrders() bool {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if localQueue[f][b] {
				return true
			}
		}
	}
	return false
}

func updateLocalQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if b != defs.ButtonCommand &&
				sharedQueue[f][b].isOrderActive &&
				sharedQueue[f][b].assignedLiftAddr == defs.Laddr.String() {
				localQueue[f][b] = true
			}
		}
	}
}

func RemoveSharedOrder(floor int, button int) {
	if button == defs.ButtonCommand {
		// error
		return
	}

	sharedQueue[floor][button] = blankOrder
}

func resetLocalQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			localQueue[f][b] = false
		}
	}
}

func resetSharedQueue() {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			sharedQueue[f][b] = blankOrder
		}
	}
}

/*func updateSharedQueue(floor int, button int) {
	// If order completed was assigned to this elevator: Remove from shared queue
	if button == defs.ButtonCommand {
		// error
		return
	}

	if sharedQueue[floor][button].isOrderActive
	&& sharedQueue[floor][button].assignedLiftAddr == laddr {
		sharedQueue[floor][button].isOrderActive = false
		sharedQueue[floor][button].assignedLiftAddr = invalidAddr
	}
}*/
