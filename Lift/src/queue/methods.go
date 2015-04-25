package queue

import (
	def "config"
	"fmt"
	"log"
	"time"
)

func (q *queue) startTimer(floor, button int) {
	//fmt.Println("run startTimer()")
	q.Q[floor][button].Timer = time.NewTimer(def.OrderTime)
	<-q.Q[floor][button].Timer.C
	OrderTimeoutChan <- def.Keypress{Button: button, Floor: floor} //bad abstraction?
}

func (q *queue) stopTimer(floor, button int) {
	//fmt.Println("run stopTimer()")
	if q.Q[floor][button].Timer != nil {
		removed := q.Q[floor][button].Timer.Stop()
		if removed {
			fmt.Printf("\n--------------------\n")
			fmt.Println("Order timer removed")
			fmt.Printf("--------------------\n\n")
		} else {
			fmt.Printf("\n--------------------\n")
			fmt.Println("Error: Order timer not removed")
			fmt.Printf("--------------------\n\n")
		}
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

// BUG(Whoever): Returns stop also if the lift should turn around, but this
// is okay(ish) because chooseDirection still returns the direction the lift
// should move (immediately) after stopping.
func (q *queue) shouldStop(floor, dir int) bool {
	switch dir {
	case def.DirDown:
		return q.isActiveOrder(floor, def.BtnDown) ||
			q.isActiveOrder(floor, def.BtnInside) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case def.DirUp:
		return q.isActiveOrder(floor, def.BtnUp) ||
			q.isActiveOrder(floor, def.BtnInside) ||
			floor == def.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case def.DirStop:
		return q.isActiveOrder(floor, def.BtnDown) ||
			q.isActiveOrder(floor, def.BtnUp) ||
			q.isActiveOrder(floor, def.BtnInside)
	default:
		log.Printf("shouldStop() called with invalid direction %d!\n", dir)
		return false
	}
}

func (q *queue) deepCopy() *queue {
	var queCopy queue
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			queCopy.Q[f][b] = q.Q[f][b]
		}
	}
	return &queCopy
}
