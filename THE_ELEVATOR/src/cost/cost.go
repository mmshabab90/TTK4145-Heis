package cost

// this doesn't belong here:
type queueEntry struct {
	isOrder bool
	ipAddr string
}
var sharedQueue [elev.NumFloors][2]queueEntry
// end does not belong

func calculateCost(targetFloor int, targetDirection elev.DirnType) int {
	cost := 0
	// Add 1 cost if lift between floors:
	floor := elev.GetFloor()
	if floor == -1 {
		cost++
		floor = fsm.GetFloor()

		// Find next floor:
		switch fsm.GetDirection() {
		case elev.DirnDown:
			floor--
		case elev.DirnUp:
			floor++
		default:
			// Error!
		}
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
			cost += 2
		}
		cost += 2
	}

	return cost
}
