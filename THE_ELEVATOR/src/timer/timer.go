package timer //rename to doortimer?

import (
	"log"
	"time"
)

var _ = log.Fatal // For debugging only, remove when done

var TimerOut = make(chan bool)
var ResetTimer = make(chan bool)
const doorOpenTime = 3 * time.Second

func Init() {
	timer := time.NewTimer(0)
	timer.Stop()

	go func() {
		for {
			select {
			case <-ResetTimer:
				timer.Reset(doorOpenTime)
			case <-timer.C:
				TimerOut <- true
				timer.Stop()
			}
		}
	}()
}
