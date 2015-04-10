package cost

func calculateCost(newOrder orderType) {
	cost := 0

	// legg til 1 cost om heisen er mellom etasjer:
	floor := elev.GetFloor()
	if floor == -1 {
		cost++
		floor = "første etasje heisen vil nå"
	}

	direction = "current direction"
	for floor != newOrder.floor && direction != newOrder.direction {
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