package fsm

func stateString(state stateType) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case open:
		return "door open"
	default:
		return "error: bad state"
	}
}

func buttonString(button int) string {
	switch button {
	case def.ButtonUp:
		return "up"
	case def.ButtonDown:
		return "down"
	case def.ButtonIn:
		return "internal"
	default:
		return "error: bad button"
	}
}
