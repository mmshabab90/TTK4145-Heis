package temp

type Keypress struct {
	Button int
	Floor int
}

func Init() {
	if !elev.Init() {
		log.Fatalln("Io_init() failed!")
	}

	elev.SetMotorDirection(elev.DirnDown)
	for elev.GetFloor() == -1 {}
	elev.SetMotorDirection(elev.DirnStop)

	fsm.Init()
	// Add some error handling here.
}

func PollKeypresses(c chan keypress) {
	var buttonState = [elev.NumFloors][elev.NumButtons]bool{
		{false, false, false},
		{false, false, false},
		{false, false, false},
		{false, false, false}, // find a better way to do this
	}

	for {
		for f := 0; f < elev.NumFloors; f++ {
			for b := 0; b < elev.NumButtons; b++ {
				if (f == 0 && b == elev.ButtonCallDown) ||
				(f == elev.NumFloors-1 && b == elev.ButtonCallUp) {
					continue
				}
				if elev.GetButton(f, b) {
					if !buttonState[f][b] {
						c <- Keypress{Button: b, Floor: f}
					}
					buttonState[f][b] = true
				} else {
					buttonState[f][b] = false
				}
			}
		}
	}
}

func PollFloor(c chan int) {
	oldFloor := elev.GetFloor()

	for {
		newFloor := elev.GetFloor()
		if newFloor != oldFloor {
			if newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
		}
	}
}
