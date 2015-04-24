package queue

import (
	def "../config"
	"fmt"
)

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	q := local.deepCopy()

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
