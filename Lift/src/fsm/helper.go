package fsm

import (
	def "../config"
)

func Direction() int {
	return dir
}

func Floor() int {
	return floor
}

func stateString(state int) string {
	switch state {
	case idle:
		return "idle"
	case moving:
		return "moving"
	case doorOpen:
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
		return "command"
	default:
		return "error: bad button"
	}
}
