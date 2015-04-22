package fsm

import (
	"time"
)

var doorReset = make(chan bool)

const doorOpenTime = 3*time.Second

func startTimer() {
	timer := time.NewTimer(0)
	timer.Stop()
	for {
		select {
		case <-doorReset:
			timer.Reset(doorOpenTime)
		case <-timer.C:
			timer.Stop()
			eventDoorTimeout <- true
		}
	}
}
