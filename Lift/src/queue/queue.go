// This is a complete rewrite of the queue package
package queue

import (
	def "../config"
	"fmt"
	"log"
	"time"
)

var _ = log.Println

type orderStatus struct {
	Active bool
	Addr   string
	Timer  *time.Timer
}

var blankOrder = orderStatus{false, "", nil}

type queue struct {
	Q [def.NumFloors][def.NumButtons]orderStatus
}

var local queue
var remote queue

var updateLocal = make(chan bool)
var backupChan = make(chan bool)
var OrderStatusTimeoutChan = make(chan orderStatus) //overkill name?
var newOrder = make(chan bool)

func Init(newOrderChan chan bool) {
	newOrder = newOrderChan
	go updateLocalQueue()
	runBackup()
	fmt.Println("queue initialized")
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, "", nil})

	newOrder <- true
}

// AddRemoteOrder adds the given order to the remote queue.
func AddRemoteOrder(floor, button int, addr string) {
	//if IsRemoteOrder(floor, button) {
	remote.setOrder(floor, button, orderStatus{true, addr /*time.NewTimer(10 * time.Second)*/, nil})
	//go remote.startTimer(floor, button)
	//}
	updateLocal <- true
	// newOrder <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		//remote.stopTimer(floor, b)
		remote.setOrder(floor, b, blankOrder)
	}

	updateLocal <- true
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

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		//remote.stopTimer(floor, b)
		local.setOrder(floor, b, blankOrder)
		remote.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // bad abstraction

	suggestBackup()
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
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if remote.Q[f][b].Addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				reassignMessage := def.Message{
					Kind:   def.NewOrder,
					Floor:  f,
					Button: b}
				def.MessageChan <- reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	orderComplete := def.Message{Kind: def.CompleteOrder, Floor: floor, Button: -1, Cost: -1}
	def.MessageChan <- orderComplete
}

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	return local.deepCopy().calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir)
}

// Print prints local and remote queue to screen in a somewhat legible
// manner.
func Print() {
	fmt.Println("Local   Remote")
	for f := def.NumFloors - 1; f >= 0; f-- {
		lifts := "   "

		if local.isActiveOrder(f, def.ButtonUp) {
			fmt.Printf("↑")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, def.ButtonCommand) {
			fmt.Printf("×")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, def.ButtonDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}
		if remote.isActiveOrder(f, def.ButtonUp) {
			fmt.Printf("↑")
			lifts += "(↑ " + def.LastPartOfIp(remote.Q[f][def.ButtonUp].Addr) + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isActiveOrder(f, def.ButtonDown) {
			fmt.Printf("↓")
			lifts += "(↓ " + def.LastPartOfIp(remote.Q[f][def.ButtonDown].Addr) + ")"
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", lifts)
		fmt.Println()
	}
}

func incrementFloor(floor, dir int) (int, int) {
	// fmt.Printf("(incr:f%v d%v)", floor, dir)
	switch dir {
	case def.DirDown:
		floor--
	case def.DirUp:
		floor++
	case def.DirStop:
		// fmt.Println("incrementFloor(): direction stop, not incremented (this is okay)")
	default:
		fmt.Println("incrementFloor(): invalid direction, not incremented")
	}

	if floor <= 0 && dir == def.DirDown {
		dir = def.DirUp
		floor = 0
	}
	if floor >= def.NumFloors-1 && dir == def.DirUp {
		dir = def.DirDown
		floor = def.NumFloors - 1
	}
	return floor, dir
}

func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if remote.isActiveOrder(f, b) {
					if b != def.ButtonCommand && remote.Q[f][b].Addr == def.Laddr.String() {
						local.setOrder(f, b, orderStatus{true, "", nil})
						newOrder <- true
					}
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}
