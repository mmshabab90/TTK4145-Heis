package cost

// Consider moving this into queue package.

import (
	"log"
	"fmt"
	"errors"
	"../hw"
	"../queue"
	"../defs"
)

// --------------- PUBLIC: ---------------

func CalculateCost(targetFloor, targetButton, currentFloor, currentDirection int) (int, error) {
	switch targetButton {
	case defs.ButtonCallUp:
	case defs.ButtonCallDown:
	case defs.ButtonCommand:
		return 0, errors.New("CalculateCost() called with internal order!")
	default:
		return 0, fmt.Errorf("CalculateCost() called with invalid order: %d\n", targetButton)
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
	
	return cost, nil
}

// --------------- PRIVATE: ---------------

func incrementFloor(floor int, direction int) int {
	switch direction {
	case defs.DirnDown:
		floor--
	case defs.DirnUp:
		floor++
	default:
		log.Println("Error: Invalid direction, floor not incremented.")
	}

	return floor
}
