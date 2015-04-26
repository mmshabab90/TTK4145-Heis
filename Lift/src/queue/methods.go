package queue

import (
	def "config"
	"log"
	"time"
)

func (q *queue) startTimer(floor, button int) {
	const orderTimeout = 30 * time.Second

	q.Q[floor][button].Timer = time.NewTimer(orderTimeout)
	<-q.Q[floor][button].Timer.C
	OrderTimeoutChan <- def.Keypress{Button: button, Floor: floor}
}

func (q *queue) stopTimer(floor, button int) {
	if q.Q[floor][button].Timer != nil {
		q.Q[floor][button].Timer.Stop()
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
	// Ignore if order to be set is equal to order already in queue.
	if q.isOrder(floor, button) == status.Active {
		return
	}
	q.Q[floor][button] = status
	def.SyncLightsChan <- true
	takeBackup <- true
	printQueues()
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
		return q.isOrder(floor, def.BtnDown) ||
			q.isOrder(floor, def.BtnInside) ||
			floor == 0 ||
			!q.isOrdersBelow(floor)
	case def.DirUp:
		return q.isOrder(floor, def.BtnUp) ||
			q.isOrder(floor, def.BtnInside) ||
			floor == def.NumFloors-1 ||
			!q.isOrdersAbove(floor)
	case def.DirStop:
		return q.isOrder(floor, def.BtnDown) ||
			q.isOrder(floor, def.BtnUp) ||
			q.isOrder(floor, def.BtnInside)
	default:
		def.CloseConnectionChan <- true
		def.Restart.Run()
		log.Fatalf(def.ColR, "This direction doesn't exist", def.ColN)
	}
}

func (q *queue) deepCopy() *queue {
	queueCopy := new(queue)
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			queueCopy.Q[f][b] = q.Q[f][b]
		}
	}
	return queueCopy
}
