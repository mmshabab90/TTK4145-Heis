package network

import (
	"time"
)

var costChan = make(chan Message) // todo find better place

var onlineLifts = make(map[string]*time.Timer)

func InitMessageHandler(deathChan chan<- string) {
	incoming := make(chan Message)
	go handleIncoming(incoming, deathChan)
}

func messageForwarder() {
}

func handleIncoming(incoming <-chan Message, deathChan chan<- string) {
	for {
		msg := <-incoming
		switch msg.Kind {

		case Alive:
			if timer, exist := onlineLifts[msg.Addr]; exist {
				timer.Reset(resetTime)
			} else {
				timer := time.NewTimer(resetTime)
				onlineLifts[msg.Addr] = timer
				go waitForDeath(deathChan, onlineLifts, msg.Addr)
			}

		case NewOrder:
			costValue := queue.CalculateCost(msg.Floor, msg.Button,
				fsm.Floor(), hw.Floor(), fsm.Direction())
			outgoing <- Message{
				Kind:   Cost,
				Floor:  msg.Floor,
				Button: msg.Button,
				Cost:   costValue}

		case CompleteOrder:
			queue.RemoveRemoteOrdersAt(msg.Floor)

		case Cost:
			costChan <- msg
		}
	}
}
