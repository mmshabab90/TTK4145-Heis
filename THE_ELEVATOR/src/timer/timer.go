package timer

import (
	"log"
	"time"
)

var _ = log.Fatal // For debugging only, remove when done

var TimerOut = make(chan bool) 
var ResetTimer = make(chan bool)
var doorOpenTime = 3 * time.Second

func Init() {
	timer := time.NewTimer(doorOpenTime)
	timer.Stop()
	
	go func() {
		for{
			select{
			case <- ResetTimer:
				timer.Reset(doorOpenTime)
			case <- timer.C:
				TimerOut <- true
				timer.Stop()
			}
		}
	}()
}

