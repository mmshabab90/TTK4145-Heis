package cost

// Consider moving this into queue package.

import (
	"log"
	"../hw"
	"../queue"
)

func CalculateCost(targetFloor, targetButton, currentFloor, currentDirection int) int {
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

	if hw.GetFloor() == -1 {
		cost += 1
		currentFloor = incrementFloor(currentFloor, currentDirection)
	}

	for !(currentFloor == targetFloor && queue.ShouldStop(currentFloor, currentDirection)) {
		log.Printf("Floor: %d, direction: %d\n", currentFloor, currentDirection)
		if queue.ShouldStop(currentFloor, currentDirection) {
			cost += 2
		}
		currentDirection = queue.ChooseDirection(currentFloor, currentDirection)
		currentFloor = incrementFloor(currentFloor, currentDirection)
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
