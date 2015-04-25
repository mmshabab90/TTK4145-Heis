package liftAssigner

import (
	def "config"
	"fmt"
	"queue"
	"time"
)

type reply struct {
	cost int
	lift string
}
type order struct { //bad name?
	floor   int
	button  int
	timeout bool
	timer   *time.Timer
}

var costTimeoutChan = make(chan order)

// Run collect cost values from all lifts
// decide which lift gets the order when all lifts
// in alive-list have answered or after a timeout
func Run(costChan chan def.Message, numberOfOnlineLifts *int) {
	assignmentQueue := make(map[order][]reply)
	const costReplyTimeout = 10 * time.Second
	var newOrder order
	for {
		select {
		case message := <-costChan:
			newOrder.makeNewOrder(message)
			newReply := getReply(message)

			for oldOrder := range assignmentQueue { //todo: make this more goood?
				if newOrder.isSameOrder(oldOrder) {
					newOrder = oldOrder
				}
			}
			// Check if order in queue
			if value, exist := assignmentQueue[newOrder]; exist {
				// Check if lift in list of that order
				found := false
				for _, e := range value {
					if e == newReply {
						found = true
					}
				}
				// Add it if not found
				if !found {
					assignmentQueue[newOrder] = append(assignmentQueue[newOrder], newReply)
					newOrder.timer.Reset(costReplyTimeout)
				}
			} else {
				// If order not in queue at all, init order list with it
				newOrder.timer = time.NewTimer(costReplyTimeout)
				assignmentQueue[newOrder] = []reply{newReply}
				go costTimer(&newOrder)
			}
			evaluateLists(&assignmentQueue, numberOfOnlineLifts)

		case newOrder := <-costTimeoutChan:
			fmt.Printf("\n COST TIMED OUT!\n\n")
			newOrder.setTimeout(true)
			evaluateLists(&assignmentQueue, numberOfOnlineLifts)
		}
	}
}

// evaluateLists goes through the map of orders with associated costs,
// checks if any orders have received answers from all live lifts or the
// timer has timed out, and finds the best candidate for all such orders.
//The best candidate is added to the shared queue.
// This is very cryptic and ungood. // todo don't admit this
func evaluateLists(que *(map[order][]reply), numberOfOnlineLifts *int) {
	const maxInt = int(^uint(0) >> 1)
	// Loop through all lists
	for order, replyList := range *que {
		// Check if the list is complete or the timer has timed out
		if len(replyList) == *numberOfOnlineLifts || order.timeout {
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
			// Print winner:
			fmt.Printf("Lift %s won order f=%d b=%d\n", lowAddr[12:15], order.floor+1, order.button)

			// Assign order to lift
			queue.AddRemoteOrder(order.floor, order.button, lowAddr)

			// Empty list and stop timer
			order.timer.Stop()
			delete(*que, order)
		}
	}
}

func getReply(m def.Message) reply {
	return reply{cost: m.Cost, lift: m.Addr}
}

func costTimer(newOrder *order) {
	<-newOrder.timer.C
	costTimeoutChan <- *newOrder
}

// --------------- METHODS FOR ORDER TYPE: ---------------

func (o *order) makeNewOrder(msg def.Message) {
	o.floor = msg.Floor
	o.button = msg.Button
}

func (o *order) isSameOrder(other order) bool {
	if other.floor == o.floor && other.button == o.button {
		return true
	} else {
		return false
	}
}

func (o *order) setTimeout(b bool) {
	o.timeout = b
}
