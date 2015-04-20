// This is a complete rewrite of the queue package
package queue

import (
	"../defs"
	"fmt"
	"log"
)

var _ = fmt.Printf
var _ = log.Printf

type orderStatus struct {
	active bool
	addr   string
}

var blankOrder = orderStatus{false, ""}

type queue struct {
	q [defs.NumFloors][defs.NumButtons]orderStatus
}

var local queue
var remote queue

var updateLocalChan = make(chan bool)

func init() {
	go updateLocalQueue()
}

// AddInternalOrder adds an order to the local queue.
func AddInternalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, ""})
}

// RemoveInternalOrder removes an order from the local queue.
func RemoveInternalOrder(floor int, button int) {
	local.setOrder(floor, button, blankOrder)
}

// AddRemoteOrder adds the given order to the remote queue.
func AddRemoteOrder(floor, button int, addr string) {
	remote.q[floor][button] = orderStatus{true, addr}

	defs.SyncLightsChan <- true
	updateLocalChan <- true
}

func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		remote.setOrder(floor, b, blankOrder)
	}

	defs.SyncLightsChan <- true
	updateLocalChan <- true
}

// ChooseDirection returns the direction the lift should continue after the
// current floor.
func ChooseDirection(currFloor, currDir int) int {
	return local.chooseDirection(currFloor, currDir)
}

// ShouldStop returns whether the lift should stop at the given floor, if
// going in the given direction.
func ShouldStop(floor, dir int) bool {
	return local.shouldStop(floor, dir)
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		local.setOrder(floor, b, blankOrder)
		remote.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // bad abstraction
}

// IsOrder returns whether there in an order with the given floor and button
// in the local queue.
func IsOrder(floor, button int) bool { // Rename to IsLocalOrder
	return local.isActiveOrder(floor, button)
}

// Blah blah blah
func IsLocalEmpty() bool {
	return local.isEmpty()
}

// IsRemoteOrder returns true if there is a order with the given floor and
// button in the remote queue.
func IsRemoteOrder(floor, button int) bool {
	return remote.isActiveOrder(floor, button)
}

// ReassignOrders finds all orders assigned to the given dead lift, removes
// them from the remote queue, and sends them on the network as new, un-
// assigned orders.
func ReassignOrders(deadAddr string) {
	// loop thru remote queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if remote.q[f][b].addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				reassignMessage := &defs.Message{
					Kind:   defs.NewOrder,
					Floor:  f,
					Button: b}
				defs.MessageChan <- *reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	message := &defs.Message{Kind: defs.CompleteOrder, Floor: floor}
	defs.MessageChan <- *message
}

func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	return local.deepCopy().calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir)
}

func Print() {
	fmt.Println("Local   Remote")
	for f := defs.NumFloors - 1; f >= 0; f-- {
		lifts := "   "

		if local.isActiveOrder(f, defs.ButtonCallUp) {
			fmt.Printf("↑")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, defs.ButtonCommand) {
			fmt.Printf("×")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, defs.ButtonCallDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}
		if remote.isActiveOrder(f, defs.ButtonCallUp) {
			fmt.Printf("↑")
			lifts += "(↑ " + defs.LastPartOfIp(remote.q[f][defs.ButtonCallUp].addr) + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isActiveOrder(f, defs.ButtonCallDown) {
			fmt.Printf("↓")
			lifts += "(↓ " + defs.LastPartOfIp(remote.q[f][defs.ButtonCallDown].addr) + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", lifts)
		fmt.Println()
	}
}

/*
 * Methods on queue struct:
 */

func (q *queue) isEmpty() bool {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.q[f][b].active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(floor, button int, status orderStatus) {
	q.q[floor][button] = status
}

func (q *queue) isActiveOrder(floor, button int) bool {
	return q.q[floor][button].active
}

func (q *queue) isOrdersAbove(floor int) bool {
	for f := floor + 1; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor, dir int) int {
	if q.isEmpty() {
		log.Println("ChooseDirection(): empty queue, returning stop")
		return defs.DirnStop
	}
	switch dir {
	case defs.DirnDown:
		if q.isOrdersBelow(floor) && floor > 0 {
			return defs.DirnDown
		} else {
			return defs.DirnUp
		}
	case defs.DirnUp:
		if q.isOrdersAbove(floor) && floor < defs.NumFloors-1 {
			return defs.DirnUp
		} else {
			return defs.DirnDown
		}
	case defs.DirnStop:
		if q.isOrdersAbove(floor) {
			return defs.DirnUp
		} else if q.isOrdersBelow(floor) {
			return defs.DirnDown
		} else {
			return defs.DirnStop
		}
	default:
		log.Printf("ChooseDirection(): called with invalid direction %d, returning stop\n", dir)
		return defs.DirnStop
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case defs.DirnDown:
		return q.isActiveOrder(floor, defs.ButtonCallDown) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case defs.DirnUp:
		return q.isActiveOrder(floor, defs.ButtonCallUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == defs.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case defs.DirnStop:
		return q.isActiveOrder(floor, defs.ButtonCallDown) ||
			q.isActiveOrder(floor, defs.ButtonCallUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}

func (q *queue) print() {
	var status string
	var lifts string

	for f := defs.NumFloors - 1; f >= 0; f-- {
		if q.q[f][defs.ButtonCallUp].active {
			status += "↑"
			lifts += "(↑ " + q.q[f][defs.ButtonCallUp].addr + ")"
		} else {
			status += " "
		}
		if q.q[f][defs.ButtonCommand].active {
			status += "×"
			lifts += "(× " + q.q[f][defs.ButtonCommand].addr + ")"
		} else {
			status += " "
		}
		if q.q[f][defs.ButtonCallDown].active {
			status += "↓   "
			lifts += "(↓ " + q.q[f][defs.ButtonCallDown].addr + ")"
		} else {
			status += " "
		}
		status += lifts + "\n"
		lifts = ""
	}
	fmt.Printf(status)
}

func (q *queue) deepCopy() *queue {
	var copy queue
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			copy.q[f][b] = q.q[f][b]
		}
	}
	return &copy
}

// this should run on a copy of local queue
func (q *queue) calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q.setOrder(targetFloor, targetButton, orderStatus{true, ""})

	cost := 0
	floor := prevFloor
	dir := currDir

	fmt.Printf("Cost floor sequence: %v", currFloor)
	// Go to valid state (a floor/dir that mirrors a button)
	if currFloor == -1 {
		// Between floors, add 1 cost
		cost++
	} else if dir != defs.DirnStop {
		// At floor, but moving, add 2 cost
		cost += 2

		if currFloor != prevFloor {
			fmt.Println("not goode: currFloor != prevFloor")
		}
	}

	floor, dir = incrementFloor(floor, dir)
	fmt.Printf(" →  %v", floor)

	for !(floor == targetFloor && q.shouldStop(floor, dir)) {
		if q.shouldStop(floor, dir) {
			cost += 2
			fmt.Printf("(S)")
		}
		dir = q.chooseDirection(floor, dir)
		floor, dir = incrementFloor(floor, dir)
		cost += 2
		fmt.Printf(" →  %v", floor)
	}
	fmt.Printf(" = cost %v\n", cost)
	return cost
}

func incrementFloor(floor, dir int) (int, int) {
	fmt.Printf("(incr:f%v d%v)", floor, dir)
	switch dir {
	case defs.DirnDown:
		floor--
	case defs.DirnUp:
		floor++
	case defs.DirnStop:
		// fmt.Println("incrementFloor(): direction stop, not incremented (this is okay)")
	default:
		fmt.Println("incrementFloor(): invalid direction, not incremented")
	}

	if floor == 0 && dir == defs.DirnDown {
		dir = defs.DirnUp
	}
	if floor == defs.NumFloors-1 && dir == defs.DirnUp {
		dir = defs.DirnDown
	}
	return floor, dir
}

func updateLocalQueue() {
	for {
		<-updateLocalChan

		for f := 0; f < defs.NumFloors; f++ {
			for b := 0; b < defs.NumButtons; b++ {
				if remote.isActiveOrder(f, b) {
					if b != defs.ButtonCommand && remote.q[f][b].addr == defs.Laddr.String() {
						// set local order f b
						local.setOrder(f, b, orderStatus{true, ""})
					}
				}
			}
		}
	}
}
