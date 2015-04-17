package cost

// Consider moving this into queue package.

import (
	"log"
	"../hw"
	"../fsm"
	"../queue"
)

func CalculateCost(targetFloor int, targetButton hw.ButtonType) int {
	switch targetButton {
	case hw.ButtonCallUp:
	case hw.ButtonCallDown:
	case hw.ButtonCommand:
		log.Println("CalculateCost() called with internal order!")
		return -1 // return something else
	default:
		log.Printf("CalculateCost() called with invalid order: %d\n", targetButton)
		return -1 // Ditto
	}
	
	cost := 0
	floor := fsm.GetFloor()
	direction := fsm.GetDirection()

	if hw.GetFloor() == -1 {
		cost += 1
		floor = incrementFloor(floor, direction)
	}

	for !(floor == targetFloor && queue.ShouldStop(floor, direction)) {
		log.Printf("Floor: %d, direction: %d\n", floor, direction)
		if queue.ShouldStop(floor, direction) {
			cost += 2
		}
		direction = queue.ChooseDirection(floor, direction)
		floor = incrementFloor(floor, direction)
		cost += 2
	}
	
	return cost
}

func incrementFloor(floor int, direction hw.DirnType) int {
	switch direction {
	case hw.DirnDown:
		floor--
	case hw.DirnUp:
		floor++
	default:
		log.Println("Error: Invalid direction, floor not incremented.")
	}

	return floor
}
