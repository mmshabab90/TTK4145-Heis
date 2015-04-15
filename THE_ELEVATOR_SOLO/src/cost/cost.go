package cost

import (
	"log"
	"../elev"
	"../fsm"
)

func CalculateCost(targetFloor int, targetDirection elev.DirnType) int {
	cost := 0

	floor := elev.GetFloor()
	if floor == -1 {
		cost++
		floor = incrementFloor(fsm.GetFloor(), direction) // Is this correct?
	}

	// Loop through floors until target found, and accumulate cost:
	direction = fsm.GetDirection()
	for floor != targetFloor && direction != targetDirection {
		if floor <= 0 {
			floor = 1
			direction *= -1
		} else if floor >= NumFloors-1 {
			floor = NumFloors - 2
			direction *= -1
		} else {
			floor += direction
		}

		if queue.ShouldStop(floor, direction) {
			if floor == targetFloor {
				break
			}
			cost += 2
		}
		cost += 2

		floor = incrementFloor(floor, direction)
		}
	}

	return cost
}

func incrementFloor(floor int, direction elev.DirnType) int {
	switch direction {
	case elev.DirnDown:
		floor--
	case elev.DirnUp:
		floor++
	default:
		log.Println("Error: Invalid direction, floor not incremented.")
	}

	return floor
}
