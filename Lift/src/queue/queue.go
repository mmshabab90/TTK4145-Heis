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
var OrderTimeoutChan = make(chan def.Keypress)
var newOrder = make(chan bool)

func Init(newOrderChan chan bool) {
	newOrder = newOrderChan
	go updateLocalQueue()
	go runBackup()
	log.Println("Queue initialized.")
}

// IsLocalOrder returns whether there in an order with the given floor and
// button in the local queue.
func IsLocalOrder(floor, button int) bool {
	return local.isOrder(floor, button)
}

// IsRemoteOrder returns true if there is a order with the given floor and
// button in the remote queue.
func IsRemoteOrder(floor, button int) bool {
	return remote.isOrder(floor, button)
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, "", nil})

	newOrder <- true
}

// AddRemoteOrder adds an order to the remote queue.
func AddRemoteOrder(floor, button int, addr string) {
	alreadyExist := remote.isOrder(floor, button)
	remote.setOrder(floor, button, orderStatus{true, addr, nil})
	if !alreadyExist {
		go remote.startTimer(floor, button)
		fmt.Printf("\n--------------------\n")
		fmt.Println("New order timer made")
		fmt.Printf("--------------------\n\n")
	}
	updateLocal <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		remote.stopTimer(floor, b)
		remote.setOrder(floor, b, blankOrder)
	}

	updateLocal <- true
}

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < def.NumButtons; b++ {
		remote.stopTimer(floor, b)
		local.setOrder(floor, b, blankOrder)
		remote.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // todo: fix bad abstraction here?

	suggestBackup()
}

// ShouldStop returns whether the lift should stop when it reaches the given
// floor, going in the given direction.
func ShouldStop(floor, dir int) bool {
	return local.shouldStop(floor, dir)
}

// ChooseDirection returns the direction the lift should continue after the
// current floor, going in the given direction.
func ChooseDirection(floor, dir int) int {
	return local.chooseDirection(floor, dir)
}

// ReassignOrders finds all orders assigned to a dead lift, removes them from
// the remote queue, and sends them on the network as new, unassigned orders.
func ReassignOrders(deadAddr string) {
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if remote.Q[f][b].Addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				def.OutgoingMsg <- def.Message{
					Description: def.NewOrder,
					Floor:       f,
					Button:      b}
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	orderComplete := def.Message{Description: def.CompleteOrder, Floor: floor, Button: -1, Cost: -1}
	def.OutgoingMsg <- orderComplete
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
				if remote.isOrder(f, b) {
					if b != def.ButtonIn && remote.Q[f][b].Addr == def.Laddr {
						local.setOrder(f, b, orderStatus{true, "", nil})
						newOrder <- true
					}
				}
			}
		}
	}
}
