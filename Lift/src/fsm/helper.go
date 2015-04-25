package fsm

import (
	def "config"
	"hw"
	"queue"
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
	case def.BtnUp:
		return "up"
	case def.BtnDown:
		return "down"
	case def.BtnInside:
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
				if (b == def.BtnUp && f == def.NumFloors-1) ||
					(b == def.BtnDown && f == 0) {
					continue
				} else {
					switch b {
					case def.BtnInside:
						hw.SetButtonLamp(f, b, queue.IsLocalOrder(f, b))
					case def.BtnUp, def.BtnDown:
						hw.SetButtonLamp(f, b, queue.IsRemoteOrder(f, b))
					}
				}
			}
		}
	}
}
