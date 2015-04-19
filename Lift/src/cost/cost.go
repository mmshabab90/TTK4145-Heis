package cost

import (
	"../defs"
)

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
func Calculate(queue [defs.NumFloors][defs.NumButtons]bool, targetFloor, targetDir, prevFloor, currFloor, currDir int) int {
	addToQueue(queue, targetFloor, targetDir)
	floor := prevFloor
	dir := currDir
	cost := 0

	if currFloor == -1 {
		cost++
		floor = increment(floor, dir)
	} else if currDir != defs.DirnStop {
		cost += 2
		floor = increment(floor, dir)
	}

	for !( (floor == targetFloor) && ( (dir == targetDir) || noOrdersAhead(queue, floor, dir) ) ) {
		if (floor == 0) || (floor >= defs.NumFloors) {
			dir *= -1
		}
		if noOrdersAhead(queue, floor, dir) {
			dir *= -1
		}
		if shouldStop(queue, floor, dir) {
			cost += 2
		}
		floor = increment(floor, dir)
	}
	return cost
}

func addToQueue(queue [defs.NumFloors][defs.NumButtons]bool, floor, dir int) {
	switch dir {
	case defs.ButtonCallDown:
		queue[floor][defs.ButtonCallDown] = true
	case defs.ButtonCallUp:
		queue[floor][defs.ButtonCallUp] = true
	default:
		// error, nothing added
	}
}

func noOrdersAhead(queue [defs.NumFloors][defs.NumButtons]bool, floor, dir int) bool {
	isOrdersAhead := false
	for f := floor; f >= 0 && f < defs.NumFloors; f += dir {
		for b := 0; b < defs.NumButtons; b++ {
			if queue[f][b] {
				isOrdersAhead = true
			}
		}
	}
	return !isOrdersAhead
}

func shouldStop(queue [defs.NumFloors][defs.NumButtons]bool, floor, dir int) bool {
	if queue[floor][defs.ButtonCommand] {
		return true
	}
	if dir == defs.DirnUp && queue[floor][defs.ButtonCallUp] {
		return true
	}
	if dir == defs.DirnDown && queue[floor][defs.ButtonCallDown] {
		return true
	}
	return false
}

func increment(floor int, dir int) int {
	switch dir {
		case defs.DirnDown:
			floor--
		case defs.DirnUp:
			floor++
		default:
			// error; no change to floor
	}
	return floor
}
