package cost

import (
	def "../config"
	"fmt"
	"time"
)

var _ = time.Sleep

// CalculateCost returns how much effort it is for this lift to carry out
// the given order. Each sheduled stop and each travel between adjacent
// floors on the way towards target will add cost 2. Cost 1 is added if the
// lift starts between floors.
func CalculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	return local.deepCopy().calculateCost(targetFloor, targetButton, prevFloor, currFloor, currDir)
}

// CalculateCost calculates how much effort it takes this lift to carry out
// the given order. Each sheduled stop on the way there and each travel
// between adjacent floors will add cost 2. Cost 1 is added if the lift
// starts between floors.
// Parameters:
// targetFloor and targetButton are self-explanatory.
// prevFloor is the most recent floor the lift has reached (same as currFloor
// if lift is at a floor).
// currFloor is the true current floor, as reported by sensors (-1 if between
// floors)
// currDir is the true current direction
func Calculate(queue [def.NumFloors][def.NumButtons]bool, targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	floor := prevFloor
	dir := currDir

	queue[targetFloor][targetButton] = true

	cost := 0

	if currFloor == -1 {
		cost++
		floor = increment(floor, dir)
	} else if currDir != def.DirnStop {
		cost += 2
		floor = increment(floor, dir)
	}

	for !(shouldStop(queue, floor, dir) && floor == targetFloor) {
		if floor <= 0 {
			dir = def.DirnUp
		} else if floor >= def.NumFloors-1 {
			dir = def.DirnDown
		}

		if dir == def.DirnStop {
			if noOrdersAhead(queue, floor, def.DirnDown) {
				dir = def.DirnUp
			} else if noOrdersAhead(queue, floor, def.DirnUp) {
				dir = def.DirnDown
			}
		}

		if shouldStop(queue, floor, dir) {
			cost += 2
		}
		floor = increment(floor, dir)
		cost += 2
	}
	return cost
}

func noOrdersAhead(queue [def.NumFloors][def.NumButtons]bool, floor, dir int) bool {
	isOrdersAhead := false
	for f := floor; isValidFloor(f); f += dir {
		for b := 0; b < def.NumButtons; b++ {
			if queue[f][b] {
				isOrdersAhead = true
			}
		}
	}
	return !isOrdersAhead
}

func shouldStop(queue [def.NumFloors][def.NumButtons]bool, floor, dir int) bool {
	//fmt.Printf("shouldStop(): floor %v, dir %v\n", floor, dir)
	if queue[floor][def.ButtonCommand] {
		return true
	}
	if dir == def.DirnUp && queue[floor][def.ButtonCallUp] {
		return true
	}
	if dir == def.DirnDown && queue[floor][def.ButtonCallDown] {
		return true
	}
	if floor == 0 && queue[floor][def.ButtonCallUp] {
		return true
	}
	if floor == def.NumFloors-1 && queue[floor][def.ButtonCallDown] {
		return true
	}
	if dir == def.DirnStop {
		for b := 0; b < def.NumButtons; b++ {
			if queue[floor][b] {
				return true
			}
		}
	}
	return false
}

func increment(floor int, dir int) int {
	switch dir {
	case def.DirnDown:
		floor--
	case def.DirnUp:
		floor++
	case def.DirnStop:
		// This is okay.
	default:
		fmt.Println("increment(): error: invalid direction, not incremented")
	}
	return floor
}

func isValidFloor(floor int) bool {
	if floor < 0 {
		return false
	}
	if floor >= def.NumFloors {
		return false
	}
	return true
}
