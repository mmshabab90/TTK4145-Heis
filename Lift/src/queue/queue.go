// This is a complete rewrite of the queue package
package queue

import (
	def "../config"
	"fmt"
	"log"
	"time"
)

var _ = fmt.Printf
var _ = log.Printf

type orderStatus struct {
	Active bool
	Addr   string
}

var blankOrder = orderStatus{false, ""}

type queue struct {
	Q [def.NumFloors][def.NumButtons]orderStatus
}

var local queue
var remote queue

var updateLocal = make(chan bool)
var backup = make(chan bool)

func init() {
	go runBackup()
	go updateLocalQueue()
}

func NewKeypress(floor, button) (notifyFsm bool) { // todo: finish this
	notifyFsm = false
	// add order to local if internal and no identical order exists
	switch button {
	case def.ButtonIn:
		if !local.isOrder(floor, button) {
			local.setOrder(floor, button, orderStatus{true, ""})
			notifyFsm = true
		}
	case def.ButtonDown, button == def.ButtonUp:
		if !shared.isOrder(floor, button) {
			// send on network
		}
	}

	return notifyFsm
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, ""})

	backup <- true
	syncLights()
}

// AddRemoteOrder adds the given order to the remote queue.
func AddRemoteOrder(floor, button int, addr string) {
	remote.setOrder(floor, button, orderStatus{true, addr})
	updateLocal <- true

	backup <- true
	syncLights()
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		remote.setOrder(floor, b, blankOrder)
	}
	updateLocal <- true

	backup <- true
	syncLights()
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		local.setOrder(floor, b, blankOrder)
		remote.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // bad abstraction
	backup <- true
}

// ChooseDirection returns the direction the lift should continue after the
// current floor.
func ChooseDirection(floor, dir int) int {
	return local.chooseDirection(floor, dir)
}

// ShouldStop returns whether the lift should stop at the given floor, if
// going in the given direction.
func ShouldStop(floor, dir int) bool {
	return local.shouldStop(floor, dir)
}

// ReassignOrders finds all orders assigned to the given dead lift, removes
// them from the remote queue, and sends them on the network as new, un-
// assigned orders.
func ReassignOrders(deadAddr string) {
	// loop thru remote queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if remote.Q[f][b].Addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				reassignMessage := def.Message{
					Kind:   def.NewOrder,
					Floor:  f,
					Button: b}
				def.Outgoing <- reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	orderComplete := def.Message{Kind: def.CompleteOrder, Floor: floor, Button: -1, Cost: -1}
	def.Outgoing <- orderComplete
}

// Print prints local and remote queue to screen in a somewhat legible
// manner.
func Print() {
	fmt.Println("Local   Remote")
	for f := def.NumFloors - 1; f >= 0; f-- {
		lifts := "   "

		if local.isOrder(f, def.ButtonUp) {
			fmt.Printf("↑")
		} else {
			fmt.Printf(" ")
		}
		if local.isOrder(f, def.ButtonIn) {
			fmt.Printf("×")
		} else {
			fmt.Printf(" ")
		}
		if local.isOrder(f, def.ButtonDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}
		if remote.isOrder(f, def.ButtonUp) {
			fmt.Printf("↑")
			lifts += "(↑ " + remote.Q[f][def.ButtonUp].Addr[12:15] + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isOrder(f, def.ButtonDown) {
			fmt.Printf("↓")
			lifts += "(↓ " + remote.Q[f][def.ButtonDown].Addr[12:15] + ")"
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
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.Q[f][b].Active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(floor, button int, status orderStatus) {
	q.Q[floor][button] = status
}

func (q *queue) isOrder(floor, button int) bool {
	return q.Q[floor][button].Active
}

func (q *queue) isOrdersAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.isOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(floor, dir int) int {
	if q.isEmpty() {
		return def.DirStop
	}
	switch dir {
	case def.DirDown:
		if q.isOrdersBelow(floor) && floor > 0 {
			return def.DirDown
		} else {
			return def.DirUp
		}
	case def.DirUp:
		if q.isOrdersAbove(floor) && floor < def.NumFloors-1 {
			return def.DirUp
		} else {
			return def.DirDown
		}
	case def.DirStop:
		if q.isOrdersAbove(floor) {
			return def.DirUp
		} else if q.isOrdersBelow(floor) {
			return def.DirDown
		} else {
			return def.DirStop
		}
	default:
		log.Printf("ChooseDirection(): called with invalid direction %d, returning stop\n", dir)
		return def.DirStop
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return q.isOrder(floor, def.ButtonDown) ||
			q.isOrder(floor, def.ButtonIn) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case def.DirUp:
		return q.isOrder(floor, def.ButtonUp) ||
			q.isOrder(floor, def.ButtonIn) ||
			floor == def.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case def.DirStop:
		return q.isOrder(floor, def.ButtonDown) ||
			q.isOrder(floor, def.ButtonUp) ||
			q.isOrder(floor, def.ButtonIn)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}
