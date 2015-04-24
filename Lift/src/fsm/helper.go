package fsm

import (
	"../defs"
	"../hw"
	"../queue"
	//"time"
)

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
	case defs.ButtonUp:
		return "up"
	case defs.ButtonDown:
		return "down"
	case defs.ButtonCommand:
		return "command"
	default:
		return "error: bad button"
	}
}

func syncLights() {
	for {
		<-defs.SyncLightsChan

		for f := 0; f < defs.NumFloors; f++ {
			for b := 0; b < defs.NumButtons; b++ {
				if (b == defs.ButtonUp && f == defs.NumFloors-1) ||
					(b == defs.ButtonDown && f == 0) {
					continue
				} else {
					switch b {
					case defs.ButtonCommand:
						hw.SetButtonLamp(f, b, queue.IsOrder(f, b))
					case defs.ButtonUp, defs.ButtonDown:
						hw.SetButtonLamp(f, b, queue.IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
