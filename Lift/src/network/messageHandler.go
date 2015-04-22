package network

import (
	"time"
)

var costChan = make(chan message) // todo find better place

var onlineLifts = make(map[string]*time.Timer)

func init() {
	incoming := make(chan message)
	go handleIncoming(incoming)
}

func handleIncoming(incoming <-chan message) {
	for {
		msg := <-incoming // formerly known as network.ReceiveChan
		switch msg.kind {

		case alive:
			if timer, exist := onlineLifts[msg.addr]; exist {
				timer.Reset(resetTime)
			} else {
				timer := time.NewTimer(resetTime)
				onlineLifts[msg.addr] = timer
				go waitForDeath(onlineLifts, msg.addr)
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
