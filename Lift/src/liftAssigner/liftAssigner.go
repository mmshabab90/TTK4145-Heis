// Package liftAssigner gathers the cost values of the lifts on the network,
// and assigns the best candidate to each order.
package liftAssigner

import (
	def "config"
	"log"
	"queue"
	"time"
)

type reply struct {
	cost int
	lift string
}
type order struct { //bad name?
	floor  int
	button int
	timer  *time.Timer
}

// Run collects cost values from all lifts for each new order, and attempts
// to find the best lift for each order, when either all online lifts have
// replied or after a timeout.
func Run(costReply <-chan def.Message, numberOfOnlineLifts *int) {
	assignmentQueue := make(map[order][]reply)

	var timeout = make(chan *order)
	const timeoutDuration = 10 * time.Second

	for {
		select {
		case message := <-costReply:
			newOrder := order{floor: message.Floor, button: message.Button}
			newReply := reply{cost: message.Cost, lift: message.Addr}

			for oldOrder := range assignmentQueue {
				if equal(oldOrder, newOrder) {
					newOrder = oldOrder
				}
			}

			// Check if order in queue
			if replyList, exist := assignmentQueue[newOrder]; exist {
				// Check if newReply is already registered
				found := false
				for _, reply := range replyList {
					if reply == newReply {
						found = true
					}
				}
				// Add it if it wasn't
				if !found {
					assignmentQueue[newOrder] = append(assignmentQueue[newOrder], newReply)
					newOrder.timer.Reset(timeoutDuration)
				}
			} else {
				// If order not in queue at all, init order list with it
				newOrder.timer = time.NewTimer(timeoutDuration)
				assignmentQueue[newOrder] = []reply{newReply}
				go costTimer(&newOrder, timeout)
			}
			chooseBestLift(&assignmentQueue, numberOfOnlineLifts, false)

		case <-timeout:
			log.Println(def.ColR, "COST TIMED OUT!", def.ColN)
			chooseBestLift(&assignmentQueue, numberOfOnlineLifts, true)
		}
	}
}

// chooseBestLift checks if any of the orders waiting for a lift assignment
// have collected enough information to have a lift assigned. For all orders
// that have, it selects a lift, and adds it to the queue.
// It assumes that all lifts always make the same decision, but if they do not,
// a timer for each order assured that this never gives unhandled orders.
func chooseBestLift(que *(map[order][]reply), numberOfOnlineLifts *int, orderTimedOut bool) {
	const maxInt = int(^uint(0) >> 1)
	// Loop through all lists
	for order, replyList := range *que {
		// Check if the list is complete or the timer has timed out
		if len(replyList) == *numberOfOnlineLifts || orderTimedOut {

			if orderTimedOut {
				log.Println("order.timeout is very true")
			}

			var lowCost = maxInt
			var lowAddr string

			// Loop through costs in each complete list
			for _, reply := range replyList {
				if reply.cost < lowCost {
					lowCost = reply.cost
					lowAddr = reply.lift
				} else if reply.cost == lowCost {
					if reply.lift < lowAddr {
						lowCost = reply.cost
						lowAddr = reply.lift
					}
				}
			}

			// Assign order to lift
			queue.AddRemoteOrder(order.floor, order.button, lowAddr)

			// Empty list and stop timer
			order.timer.Stop()
			delete(*que, order)
		}
	}
}

func costTimer(newOrder *order, timeout chan *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

func equal(o1, o2 order) bool {
	return o1.floor == o2.floor && o1.button == o2.button
}
