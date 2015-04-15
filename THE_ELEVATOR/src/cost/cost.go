package cost

import (
	"log"
	"../elev"
	"../fsm"
	"../queue"
)

func CalculateCost(targetFloor int, targetButton elev.ButtonType) int {
	var targetDirection elev.DirnType
	switch targetButton {
	case elev.ButtonCallUp:
		targetDirection = elev.DirnUp
	case elev.ButtonCallDown:
		targetDirection = elev.DirnDown
	default:
		log.Println("Error dir cost")
	}

	cost := 0

	floor := elev.GetFloor()
	direction := fsm.GetDirection()
	if floor == -1 {
		cost++
		floor = incrementFloor(floor, direction) // Is this correct?
	}

	// Loop through floors until target found, and accumulate cost:
	for floor != targetFloor && direction != targetDirection {
		if floor <= 0 {
			floor = 1
			direction *= -1
		} else if floor >= elev.NumFloors-1 {
			floor = elev.NumFloors - 2
			direction *= -1
		} else {
			floor = incrementFloor(floor, direction)
			//floor += direction
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
