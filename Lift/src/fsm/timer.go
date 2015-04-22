package fsm

import (
	"time"
)

const doorOpenTime = 1*time.Second

func startTimer(doorReset chan bool, eventDoorTimeout chan bool) {
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

