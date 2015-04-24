// This is a complete rewrite of the queue package
package queue

import (
	"../defs"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

var _ = fmt.Printf
var _ = log.Printf

const diskDebug = false

type orderStatus struct {
	Active bool
	Addr   string
	Timer  *time.Timer
}

var blankOrder = orderStatus{false, "", nil}

type queue struct {
	Q [defs.NumFloors][defs.NumButtons]orderStatus
}

var local queue
var remote queue

var updateLocal = make(chan bool)
var backup = make(chan bool)
var OrderStatusTimeoutChan = make(chan orderStatus) //overkill name?
var newOrder = make(chan bool)

func Init(newOrderChan chan bool) {
	newOrder = newOrderChan
	go runBackup()
	go updateLocalQueue()
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(floor int, button int) {
	local.setOrder(floor, button, orderStatus{true, "", nil})

	newOrder <- true
	backup <- true
}

// AddRemoteOrder adds the given order to the remote queue.
func AddRemoteOrder(floor, button int, addr string) {
	remote.setOrder(floor, button, orderStatus{true, addr, time.NewTimer(10 * time.Second)})
	//go remote.startTimer(floor, button)

	updateLocal <- true
	// newOrder <- true
	backup <- true
}

// RemoveRemoteOrdersAt removes all orders at the given floor from the remote
// queue.
func RemoveRemoteOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		//remote.stopTimer(floor, b)
		remote.setOrder(floor, b, blankOrder)
	}

	updateLocal <- true
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

// RemoveOrdersAt removes all orders at the given floor in local and remote queue.
func RemoveOrdersAt(floor int) {
	for b := 0; b < defs.NumButtons; b++ {
		//remote.stopTimer(floor, b)
		local.setOrder(floor, b, blankOrder)
		remote.setOrder(floor, b, blankOrder)
	}
	SendOrderCompleteMessage(floor) // bad abstraction
	backup <- true
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
			if remote.Q[f][b].Addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				reassignMessage := defs.Message{
					Kind:   defs.NewOrder,
					Floor:  f,
					Button: b}
				defs.MessageChan <- reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given floor.
func SendOrderCompleteMessage(floor int) {
	orderComplete := defs.Message{Kind: defs.CompleteOrder, Floor: floor, Button: -1, Cost: -1}
	defs.MessageChan <- orderComplete
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
	for f := defs.NumFloors - 1; f >= 0; f-- {
		lifts := "   "

		if local.isActiveOrder(f, defs.ButtonUp) {
			fmt.Printf("↑")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, defs.ButtonCommand) {
			fmt.Printf("×")
		} else {
			fmt.Printf(" ")
		}
		if local.isActiveOrder(f, defs.ButtonDown) {
			fmt.Printf("↓   %d  ", f+1)
		} else {
			fmt.Printf("    %d  ", f+1)
		}
		if remote.isActiveOrder(f, defs.ButtonUp) {
			fmt.Printf("↑")
			lifts += "(↑ " + defs.LastPartOfIp(remote.Q[f][defs.ButtonUp].Addr) + ")"
		} else {
			fmt.Printf(" ")
		}
		if remote.isActiveOrder(f, defs.ButtonDown) {
			fmt.Printf("↓")
			lifts += "(↓ " + defs.LastPartOfIp(remote.Q[f][defs.ButtonDown].Addr) + ")"
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

func (q *queue) startTimer(floor, button int) {
	fmt.Println("run startTimer()")
	//q.Q[floor][button].Timer = time.NewTimer(10*time.Second)
	<-q.Q[floor][button].Timer.C
	OrderStatusTimeoutChan <- q.Q[floor][button]
}

func (q *queue) stopTimer(floor, button int) {
	fmt.Println("run stopTimer()")
	if q.Q[floor][button].Timer != nil {
		removed := q.Q[floor][button].Timer.Stop()
		if removed {
			fmt.Println("timer removed")
		} else {
			fmt.Println("Timer not removed")
		}
	} else {
		fmt.Println("Timer was nil")
	}
}

func (q *queue) isEmpty() bool {
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.Q[f][b].Active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(floor, button int, status orderStatus) {
	q.Q[floor][button] = status
	defs.SyncLightsChan <- true
	/*if button == defs.ButtonCommand {
		newOrder <- true
	} else if status.Active && status.Addr == defs.Laddr.String() {
		newOrder <- true
	}*/
}

func (q *queue) isActiveOrder(floor, button int) bool {
	return q.Q[floor][button].Active
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
		return defs.DirStop
	}
	switch dir {
	case defs.DirDown:
		if q.isOrdersBelow(floor) && floor > 0 {
			return defs.DirDown
		} else {
			return defs.DirUp
		}
	case defs.DirUp:
		if q.isOrdersAbove(floor) && floor < defs.NumFloors-1 {
			return defs.DirUp
		} else {
			return defs.DirDown
		}
	case defs.DirStop:
		if q.isOrdersAbove(floor) {
			return defs.DirUp
		} else if q.isOrdersBelow(floor) {
			return defs.DirDown
		} else {
			return defs.DirStop
		}
	default:
		log.Printf("ChooseDirection(): called with invalid direction %d, returning stop\n", dir)
		return defs.DirStop
	}
}

func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case defs.DirDown:
		return q.isActiveOrder(floor, defs.ButtonDown) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case defs.DirUp:
		return q.isActiveOrder(floor, defs.ButtonUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand) ||
			floor == defs.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case defs.DirStop:
		return q.isActiveOrder(floor, defs.ButtonDown) ||
			q.isActiveOrder(floor, defs.ButtonUp) ||
			q.isActiveOrder(floor, defs.ButtonCommand)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}

func (q *queue) deepCopy() *queue {
	var copy queue
	for f := 0; f < defs.NumFloors; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			copy.Q[f][b] = q.Q[f][b]
		}
	}
	return &copy
}

// this should run on a copy of local queue
func (q *queue) calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q.setOrder(targetFloor, targetButton, orderStatus{true, "", nil})

	cost := 0
	floor := prevFloor
	dir := currDir

	if currFloor == -1 {
		// Between floors, add 1 cost
		cost++
	} else if dir != defs.DirStop {
		// At floor, but moving, add 2 cost
		cost += 2
	}

	floor, dir = incrementFloor(floor, dir)
	fmt.Printf("Cost floor sequence: %v →  %v", currFloor, floor)

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
	// fmt.Printf("(incr:f%v d%v)", floor, dir)
	switch dir {
	case defs.DirDown:
		floor--
	case defs.DirUp:
		floor++
	case defs.DirStop:
		// fmt.Println("incrementFloor(): direction stop, not incremented (this is okay)")
	default:
		fmt.Println("incrementFloor(): invalid direction, not incremented")
	}

	if floor <= 0 && dir == defs.DirDown {
		dir = defs.DirUp
		floor = 0
	}
	if floor >= defs.NumFloors-1 && dir == defs.DirUp {
		dir = defs.DirDown
		floor = defs.NumFloors - 1
	}
	return floor, dir
}

func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < defs.NumFloors; f++ {
			for b := 0; b < defs.NumButtons; b++ {
				if remote.isActiveOrder(f, b) {
					if b != defs.ButtonCommand && remote.Q[f][b].Addr == defs.Laddr.String() {
						local.setOrder(f, b, orderStatus{true, "", nil})
						newOrder <- true
					}
				}
			}
		}
		time.Sleep(time.Millisecond)
	}
}

// runBackup loads queue data from file if file exists once and saves backups
// whenever its asked to.
func runBackup() {
	filenameLocal := "localQueueBackup"
	filenameRemote := "remoteQueueBackup"

	local.loadFromDisk(filenameLocal)
	// remote.loadFromDisk(filenameRemote)

	for {
		<-backup
		if err := local.saveToDisk(filenameLocal); err != nil {
			fmt.Println(err)
		}
		if err := remote.saveToDisk(filenameRemote); err != nil {
			fmt.Println(err)
		}
	}
}

func (q *queue) saveToDisk(filename string) error {
	fi, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fi.Close()

	if err := gob.NewEncoder(fi).Encode(q); err != nil {
		return err
	}

	if diskDebug {
		fmt.Printf("Successful save of file %s\n", filename)
	}
	return nil
}

// loadFromDisk checks if a file of the given name is available on disk, and
// saves its contents to the queue it's invoked on if the file is present.
func (q *queue) loadFromDisk(filename string) error {
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("Backup file %s exists, processing...\n", filename)
		fi, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer fi.Close()

		if err := gob.NewDecoder(fi).Decode(&q); err != nil {
			return err
		}
	}

	// Ny ide: If not empty, event button pressed.

	return nil
}
