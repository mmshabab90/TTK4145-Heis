package network

import (
	"time"
)

var costChan = make(chan message) // todo find better place

var onlineLifts = make(map[string]*time.Timer)

func InitMessageHandler(deathChan chan<- string) {
	incoming := make(chan message)
	go handleIncoming(incoming, deathChan)
}

func handleIncoming(incoming <-chan message, deathChan chan<- string) {
	for {
		msg := <-incoming // formerly known as network.ReceiveChan
		switch msg.kind {

		case alive:
			if timer, exist := onlineLifts[msg.addr]; exist {
				timer.Reset(resetTime)
			} else {
				timer := time.NewTimer(resetTime)
				onlineLifts[msg.addr] = timer
				go waitForDeath(deathChan, onlineLifts, msg.addr)
			}

		case newOrder:
			costValue := queue.CalculateCost(msg.floor, msg.button,
				fsm.Floor(), hw.Floor(), fsm.Direction())
			outgoing <- message{
				kind:   cost,
				floor:  msg.floor,
				button: msg.button,
				cost:   costValue}

		case completeOrder:
			queue.RemoveRemoteOrdersAt(msg.floor)

		case cost:
			costChan <- msg
		}
	}
}
