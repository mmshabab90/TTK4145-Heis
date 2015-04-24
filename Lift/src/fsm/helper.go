package fsm

import (
	def "../config"
	"../hw"
	"../queue"
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
	case def.ButtonCommand:
		return "command"
	default:
		return "error: bad button"
	}
}

func syncLights() { // todo: probably move to queue
	for {
		<-def.SyncLightsChan

		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if (b == def.ButtonUp && f == def.NumFloors-1) ||
					(b == def.ButtonDown && f == 0) {
					continue
				} else {
					switch b {
					case def.ButtonCommand:
						hw.SetButtonLamp(f, b, queue.IsLocalOrder(f, b))
					case def.ButtonUp, def.ButtonDown:
						hw.SetButtonLamp(f, b, queue.IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
