package cost

import (
	"log"
	"../elev"
	"../fsm"
	"../queue"
)

func CalculateCost(targetFloor int, targetButton elev.ButtonType) int {
	// Really need some good error taking care of here
	var targetDirection elev.DirnType

	switch targetButton {
	case elev.ButtonCallUp:
		targetDirection = elev.DirnUp
	case elev.ButtonCallDown:
		targetDirection = elev.DirnDown
	case elev.ButtonCommand:
		log.Println("Should not calculate cost on internal command!")
		return -1
	default:
		log.Fatalln("Error direction in cost")
	}

	cost := 0

	floor := elev.GetFloor()
	direction := fsm.GetDirection()
	log.Printf("Floor: %d, direction: %d\n", floor, direction)
	if floor == -1 {
		cost += 1
		floor = incrementFloor(fsm.GetFloor(), direction) // Is this correct?
	}

	// Loop through floors until target found, and accumulate cost:
	for floor != targetFloor && direction != targetDirection {
		// Assert something

		// Handle top/bottom floors:
		if floor <= 0 {
			floor = 1
			direction = elev.DirnUp
		} else if floor >= elev.NumFloors - 1 {
			floor = elev.NumFloors - 2
			direction = elev.DirnDown
		} else {
			floor = incrementFloor(floor, direction)
			direction = queue.ChooseDirection(floor, direction)
		}
		cost += 2

		log.Printf("Floor: %d, direction: %d\n", floor, direction)

		if queue.ShouldStop(floor, direction) {
			if floor == targetFloor {
				break
			}
			cost += 2
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
