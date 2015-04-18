package cost

import (
	"../defs"
	"../queue"
	"errors"
	"fmt"
	"log"
)

// --------------- PUBLIC: ---------------

func CalculateCost(targetFloor, targetButton, fsmFloor, fsmDir, currFloor int) (int, error) {
	switch targetButton {
	case defs.ButtonCallUp, defs.ButtonCallDown:
	case defs.ButtonCommand:
		return 0, errors.New("CalculateCost() called with internal order!")
	default:
		return 0, fmt.Errorf("CalculateCost() called with invalid order: %d\n", targetButton)
	}

	cost := 0

	if currFloor == -1 {
		cost += 1
		fsmFloor = incrementFloor(fsmFloor, fsmDir)
	}

	for !(fsmFloor == targetFloor && queue.ShouldStop(fsmFloor, fsmDir)) {
		log.Printf("Floor: %d, direction: %d\n", fsmFloor, fsmDir)
		if queue.ShouldStop(fsmFloor, fsmDir) {
			cost += 2
		}
		fsmDir = queue.ChooseDirection(fsmFloor, fsmDir)
		fsmFloor = incrementFloor(fsmFloor, fsmDir)
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
