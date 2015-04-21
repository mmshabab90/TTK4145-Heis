package network

import (
	def "../config"
	"time"
)

var costChan = make(chan Message) // find better place

func init() {
	aliveLifts := make(map[string]*time.Timer)
	incoming := make(chan Message)
	go handleIncoming(incoming, aliveLifts)
}

func handleIncoming(incoming <-chan Message, aliveLifts map[string]*time.Timer) {
	for {
		msg := <-incoming // formerly known as network.ReceiveChan
		switch msg.Kind {

		case Alive:
			if timer, exist := aliveLifts[msg.Addr]; exist {
				timer.Timer.Reset(resetTime)
			} else {
				timer := time.NewTimer(resetTime)
				onlineLifts[msg.Addr] = timer
				go waitForDeath(aliveLifts, &connection, msg.Addr)
			}

		case NewOrder:
			cost := queue.CalculateCost(msg.Floor, msg.Button,
				fsm.Floor(), hw.Floor(), fsm.Direction())
			def.Outgoing <- Message{
				Kind:   Cost,
				Floor:  msg.Floor,
				Button: msg.Button,
				Cost:   cost}

		case CompleteOrder:
			queue.RemoveRemoteOrdersAt(msg.Floor)

		case Cost:
			costChan <- msg
		}
	}
}
