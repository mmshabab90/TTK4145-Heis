package queue

import (
	def "../config"
)

func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if remote.isActiveOrder(f, b) {
					if b != def.ButtonIn && remote.Q[f][b].Addr == def.Laddr.String() {
						local.setOrder(f, b, orderStatus{true, ""})
					}
				}
			}
		}
	}
}

func syncLights() {
	for f := 0; f < def.NumFloors; f++ {
		for b := 0; b < def.NumButtons; b++ {
			switch b {
			case def.ButtonUp:
				if f != def.NumFloors-1 && remote.isActiveOrder(f, b) {
					setButtonLamp <- keypress{f, b}
				}
			case def.ButtonDown:
				if f != 0 && remote.isActiveOrder(f, b) {
					setButtonLamp <- keypress{f, b}
				}
			case def.ButtonIn:
				if local.isActiveOrder(f, b) {
					setButtonLamp <- keypress{f, b}
				}
			}
		}
	}
}
