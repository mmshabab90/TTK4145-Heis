package cost

import (
	"../defs"
	"fmt"
	"time"
)

var _ = time.Sleep

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
func Calculate(queue [defs.NumFloors][defs.NumButtons]bool, targetFloor, targetButton, prevFloor, currFloor, currDir int) int {
	/*
	fmt.Println("Calculate called with")
	fmt.Printf("   targetFloor %v\n", targetFloor)
	fmt.Printf("   targetButton %v\n", targetButton)
	fmt.Printf("   prevFloor %v\n", prevFloor)
	fmt.Printf("   currFloor %v\n", currFloor)
	fmt.Printf("   currDir %v\n", currDir)
	*/
	floor := prevFloor
	dir := currDir

	queue[targetFloor][targetButton] = true
	
	/*
	var targetDir int
	if targetButton == defs.ButtonCallDown {
		targetDir = defs.DirnDown
	} else if targetButton == defs.ButtonCallUp {
		targetDir = defs.DirnUp
	}*/
	
	//fmt.Printf("  (targetDir %v)\n", targetDir)

	//fmt.Printf("queue: %v\n", queue)

	cost := 0

	if currFloor == -1 {
		cost++
		floor = increment(floor, dir)
	} else if currDir != defs.DirnStop {
		cost += 2
		floor = increment(floor, dir)
	}
	
	//fmt.Printf("floor = %v, dir = %v\n", floor, dir)
	for !(shouldStop(queue, floor, dir) && floor == targetFloor) {
		if floor <= 0 {
			dir = defs.DirnUp
		} else if floor >= defs.NumFloors-1 {
			dir = defs.DirnDown
		}
		
		if dir == defs.DirnStop {
			if noOrdersAhead(queue, floor, defs.DirnDown) {
				dir = defs.DirnUp
			} else if noOrdersAhead(queue, floor, defs.DirnUp) {
				dir = defs.DirnDown
			}
		}
		
		// if noOrdersAhead(queue, floor, dir) {
		// 	dir *= -1
		// }
		
		if shouldStop(queue, floor, dir) {
			cost += 2
		}
		floor = increment(floor, dir)
		cost += 2
		//fmt.Printf("floor = %v, dir = %v\n", floor, dir)
		//time.Sleep(500*time.Millisecond)
	}
	return cost
}

func noOrdersAhead(queue [defs.NumFloors][defs.NumButtons]bool, floor, dir int) bool {
	//fmt.Printf("noOrdersAhead() running with floor %v and dir %v\n", floor, dir)
	isOrdersAhead := false
	for f := floor; isValidFloor(f); f += dir {
		for b := 0; b < defs.NumButtons; b++ {
			if queue[f][b] {
				isOrdersAhead = true
			}
		}
	}
	return !isOrdersAhead
}

func shouldStop(queue [defs.NumFloors][defs.NumButtons]bool, floor, dir int) bool {
	//fmt.Printf("shouldStop(): floor %v, dir %v\n", floor, dir)
	if queue[floor][defs.ButtonCommand] {
		return true
	}
	if dir == defs.DirnUp && queue[floor][defs.ButtonCallUp] {
		return true
	}
	if dir == defs.DirnDown && queue[floor][defs.ButtonCallDown] {
		return true
	}
	if floor == 0 && queue[floor][defs.ButtonCallUp] {
		return true
	}
	if floor == defs.NumFloors - 1 && queue[floor][defs.ButtonCallDown] {
		return true
	}
	if dir == defs.DirnStop {
		for b := 0; b < defs.NumButtons; b++ {
			if queue[floor][b] {
				return true
			}
		}
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
			fmt.Println("increment(): error: invalid direction, not incremented")
	}
	return floor
}

func isValidFloor(floor int) bool {
	if floor < 0 {
		return false
	}
	if floor >= defs.NumFloors {
		return false
	}
	return true
}
