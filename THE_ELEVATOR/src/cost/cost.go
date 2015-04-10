package cost

// this doesn't belong here:
type queueEntry struct {
	isOrder bool
	ipAddr string
}
var sharedQueue [elev.NumFloors][2]queueEntry
// end does not belong

func calculateCost(targetFloor int, targetDirection elev.DirnType) {
	cost := 0

	// legg til 1 cost om heisen er mellom etasjer:
	floor := elev.GetFloor()
	if floor == -1 {
		cost++
		floor = fsm.GetFloor()

		switch fsm.GetDirection() {
		case elev.DirnDown:
			floor--
		case elev.DirnUp:
			floor++
		default:
			// Error!
		}
	}

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

		cost += 2
		if shouldStop() {
			cost += 2
		}
	}

	return cost

	// finn antall stopp før ankomst til newOrder:
	// finn antall etasjer å kjøre før ankomst til newOrder:
	// (pass på at det tas høyde for retningen heisen evt. kjører)
}

func shouldStop() {
	
}