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
}

var blankOrder = orderStatus{false, ""}

type queue struct {
	Q [defs.NumStoreys][defs.NumButtons]orderStatus
}

var local queue
var remote queue

var updateLocal = make(chan bool)
var backup = make(chan bool)

func init() {
	go runBackup()
	go updateLocalQueue()
}

// AddLocalOrder adds an order to the local queue.
func AddLocalOrder(storey int, button int) {
	local.setOrder(storey, button, orderStatus{true, ""})

	backup <- true
}

// AddRemoteOrder adds the given order to the remote queue.
func AddRemoteOrder(storey, button int, addr string) {
	remote.setOrder(storey, button, orderStatus{true, addr})

	defs.SyncLightsChan <- true
	updateLocal <- true
	backup <- true
}

// RemoveRemoteOrdersAt removes all orders at the given storey from the remote
// queue.
func RemoveRemoteOrdersAt(storey int) {
	for b := 0; b < defs.NumButtons; b++ {
		remote.setOrder(storey, b, blankOrder)
	}

	defs.SyncLightsChan <- true
	updateLocal <- true
	backup <- true
}

// ChooseDirection returns the direction the lift should continue after the
// current storey.
func ChooseDirection(storey, dir int) int {
	return local.chooseDirection(storey, dir)
}

// ShouldStop returns whether the lift should stop at the given storey, if
// going in the given direction.
func ShouldStop(storey, dir int) bool {
	return local.shouldStop(storey, dir)
}

// RemoveOrdersAt removes all orders at the given storey in local and remote queue.
func RemoveOrdersAt(storey int) {
	for b := 0; b < defs.NumButtons; b++ {
		local.setOrder(storey, b, blankOrder)
		remote.setOrder(storey, b, blankOrder)
	}
	SendOrderCompleteMessage(storey) // bad abstraction
	backup <- true
}

// IsOrder returns whether there in an order with the given storey and button
// in the local queue.
func IsOrder(storey, button int) bool { // Rename to IsLocalOrder
	return local.isActiveOrder(storey, button)
}

// Blah blah blah
func IsLocalEmpty() bool {
	return local.isEmpty()
}

// IsRemoteOrder returns true if there is a order with the given storey and
// button in the remote queue.
func IsRemoteOrder(storey, button int) bool {
	return remote.isActiveOrder(storey, button)
}

// ReassignOrders finds all orders assigned to the given dead lift, removes
// them from the remote queue, and sends them on the network as new, un-
// assigned orders.
func ReassignOrders(deadAddr string) {
	// loop thru remote queue
	// remove all orders assigned to the dead lift
	// send neworder-message for each removed order
	for f := 0; f < defs.NumStoreys; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if remote.Q[f][b].Addr == deadAddr {
				remote.setOrder(f, b, blankOrder)
				reassignMessage := defs.Message{
					Kind:   defs.NewOrder,
					Storey: f,
					Button: b}
				defs.MessageChan <- reassignMessage
			}
		}
	}
}

// SendOrderCompleteMessage communicates to the network that this lift has
// taken care of orders at the given storey.
func SendOrderCompleteMessage(storey int) {
	orderComplete := defs.Message{Kind: defs.CompleteOrder, Storey: storey, Button: -1, Cost: -1}
	defs.MessageChan <- orderComplete
}

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// storeys on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between storeys.
func CalculateCost(targetStorey, targetButton, prevStorey, currStorey, currDir int) int {
	return local.deepCopy().calculateCost(targetStorey, targetButton, prevStorey, currStorey, currDir)
}

// Print prints local and remote queue to screen in a somewhat legible
// manner.
func Print() {
	fmt.Println("Local   Remote")
	for f := defs.NumStoreys - 1; f >= 0; f-- {
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

func (q *queue) isEmpty() bool {
	for f := 0; f < defs.NumStoreys; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.Q[f][b].Active {
				return false
			}
		}
	}
	return true
}

func (q *queue) setOrder(storey, button int, status orderStatus) {
	q.Q[storey][button] = status
}

func (q *queue) isActiveOrder(storey, button int) bool {
	return q.Q[storey][button].Active
}

func (q *queue) isOrdersAbove(storey int) bool {
	for f := storey + 1; f < defs.NumStoreys; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(storey int) bool {
	for f := 0; f < storey; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) chooseDirection(storey, dir int) int {
	if q.isEmpty() {
		return defs.DirStop
	}
	switch dir {
	case defs.DirDown:
		if q.isOrdersBelow(storey) && storey > 0 {
			return defs.DirDown
		} else {
			return defs.DirUp
		}
	case defs.DirUp:
		if q.isOrdersAbove(storey) && storey < defs.NumStoreys-1 {
			return defs.DirUp
		} else {
			return defs.DirDown
		}
	case defs.DirStop:
		if q.isOrdersAbove(storey) {
			return defs.DirUp
		} else if q.isOrdersBelow(storey) {
			return defs.DirDown
		} else {
			return defs.DirStop
		}
	default:
		log.Printf("ChooseDirection(): called with invalid direction %d, returning stop\n", dir)
		return defs.DirStop
	}
}

func (q *queue) shouldStop(storey, dir int) bool {
	switch dir {
	case defs.DirDown:
		return q.isActiveOrder(storey, defs.ButtonDown) ||
			q.isActiveOrder(storey, defs.ButtonCommand) ||
			storey == 0 ||
			!q.isOrdersBelow(storey)
	case defs.DirUp:
		return q.isActiveOrder(storey, defs.ButtonUp) ||
			q.isActiveOrder(storey, defs.ButtonCommand) ||
			storey == defs.NumStoreys-1 ||
			!q.isOrdersAbove(storey)
	case defs.DirStop:
		return q.isActiveOrder(storey, defs.ButtonDown) ||
			q.isActiveOrder(storey, defs.ButtonUp) ||
			q.isActiveOrder(storey, defs.ButtonCommand)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}

func (q *queue) deepCopy() *queue {
	var copy queue
	for f := 0; f < defs.NumStoreys; f++ {
		for b := 0; b < defs.NumButtons; b++ {
			copy.Q[f][b] = q.Q[f][b]
		}
	}
	return &copy
}

// this should run on a copy of local queue
func (q *queue) calculateCost(targetStorey, targetButton, prevStorey, currStorey, currDir int) int {
	q.setOrder(targetStorey, targetButton, orderStatus{true, ""})

	cost := 0
	storey := prevStorey
	dir := currDir

	if currStorey == -1 {
		// Between storeys, add 1 cost
		cost++
	} else if dir != defs.DirStop {
		// At storey, but moving, add 2 cost
		cost += 2
	}

	storey, dir = incrementStorey(storey, dir)
	fmt.Printf("Cost storey sequence: %v →  %v", currStorey, storey)

	for !(storey == targetStorey && q.shouldStop(storey, dir)) {
		if q.shouldStop(storey, dir) {
			cost += 2
			fmt.Printf("(S)")
		}
		dir = q.chooseDirection(storey, dir)
		storey, dir = incrementStorey(storey, dir)
		cost += 2
		fmt.Printf(" →  %v", storey)
	}
	fmt.Printf(" = cost %v\n", cost)
	return cost
}

func incrementStorey(storey, dir int) (int, int) {
	// fmt.Printf("(incr:f%v d%v)", storey, dir)
	switch dir {
	case defs.DirDown:
		storey--
	case defs.DirUp:
		storey++
	case defs.DirStop:
		// fmt.Println("incrementStorey(): direction stop, not incremented (this is okay)")
	default:
		fmt.Println("incrementStorey(): invalid direction, not incremented")
	}

	if storey <= 0 && dir == defs.DirDown {
		dir = defs.DirUp
		storey = 0
	}
	if storey >= defs.NumStoreys-1 && dir == defs.DirUp {
		dir = defs.DirDown
		storey = defs.NumStoreys - 1
	}
	return storey, dir
}

func updateLocalQueue() {
	fmt.Println("updateLocalQueue() routine running...")
	for {
		<-updateLocal

		for f := 0; f < defs.NumStoreys; f++ {
			for b := 0; b < defs.NumButtons; b++ {
				if remote.isActiveOrder(f, b) {
					if b != defs.ButtonCommand && remote.Q[f][b].Addr == defs.Laddr.String() {
						local.setOrder(f, b, orderStatus{true, ""})
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
