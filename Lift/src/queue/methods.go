package queue

import (
	def "../config"
	"fmt"
	"log"
)

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
	// Ignore if order to be set is equal to order already in queue:
	if q.isActiveOrder(floor, button) == status.Active {
		return
	}

	q.Q[floor][button] = status
	def.SyncLightsChan <- true
	suggestBackup()
}

func (q *queue) isActiveOrder(floor, button int) bool { // todo: consider rename
	return q.Q[floor][button].Active
}

func (q *queue) isOrdersAbove(floor int) bool {
	for f := floor + 1; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
				return true
			}
		}
	}
	return false
}

func (q *queue) isOrdersBelow(floor int) bool {
	for f := 0; f < floor; f++ {
		for b := 0; b < def.NumButtons; b++ {
			if q.isActiveOrder(f, b) {
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
		return q.isActiveOrder(floor, def.ButtonDown) ||
			q.isActiveOrder(floor, def.ButtonCommand) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case def.DirUp:
		return q.isActiveOrder(floor, def.ButtonUp) ||
			q.isActiveOrder(floor, def.ButtonCommand) ||
			floor == def.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case def.DirStop:
		return q.isActiveOrder(floor, def.ButtonDown) ||
			q.isActiveOrder(floor, def.ButtonUp) ||
			q.isActiveOrder(floor, def.ButtonCommand)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}

func (q *queue) deepCopy() *queue {
	var copy queue
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
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
	} else if dir != def.DirStop {
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
