package timer

import (
	"log"
	"time"
)

var _ = log.Fatal // for debugging only

var TimerChan = make(chan bool) 

func Start(sleepytime int) {
	time.Sleep(time.Duration(sleepytime) * time.Second)
	TimerChan <- true
}

