package queue

import (
	def "../config"
)

func updateLocalQueue() {
	for {
		<-updateLocal
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				if remote.isOrder(f, b) {
					if b != def.ButtonIn && remote.Q[f][b].Addr == def.Laddr.String() {
						local.setOrder(f, b, orderStatus{true, ""})
					}
				}
			}
		}
	}
}

func syncLights(setButtonLamp chan<- def.Keypress) {
	for {
		<- syncChan
		for f := 0; f < def.NumFloors; f++ {
			for b := 0; b < def.NumButtons; b++ {
				switch b {
				case def.ButtonUp:
					if f != def.NumFloors-1 && remote.isOrder(f, b) {
						setButtonLamp <- def.Keypress{f, b}
					}
				case def.ButtonDown:
					if f != 0 && remote.isOrder(f, b) {
						setButtonLamp <- def.Keypress{f, b}
					}
				case def.ButtonIn:
					if local.isOrder(f, b) {
						setButtonLamp <- def.Keypress{f, b}
					}
				}
			}
		}
	}
}
