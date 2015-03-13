package timer

import (
	"log"
	"time"
)

var _ = log.Fatal // for debugging only

var TimerOut = make(chan bool) 
var ResetTimer = make(chan bool)

/*func Start(sleepytime int) {
	time.Sleep(time.Duration(sleepytime) * time.Second)
	TimerChan <- true
}*/

func Timer() {
	myTimer := time.NewTimer(3 * time.Second)
	myTimer.Stop()
	for{
		select{
		case <- ResetTimer:
			myTimer.Reset(3 * time.Second)
		case <- myTimer.C:
			TimerOut <- true
			myTimer.Stop()
		}
	}
}

