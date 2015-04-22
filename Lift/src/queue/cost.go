package queue

import (
	def "../config"
	"fmt"
	"log"
)

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	return local.deepCopy().calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir)
}

func (q *queue) calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q.setOrder(targetFloor, targetButton, orderStatus{true, ""})

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
		log.Println("incrementFloor(): invalid direction, not incremented")
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

func (q *queue) deepCopy() *queue {
	var queueCopy queue
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			queueCopy.Q[f][b] = q.Q[f][b]
		}
	}
	return &queueCopy
}
