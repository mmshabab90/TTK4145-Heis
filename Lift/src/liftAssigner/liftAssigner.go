package liftAssigner

import (
	def "config"
	"fmt"
	"log"
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

// Run collects cost values from all lifts for each new order, and assigns
// a lift to each order when either all online lifts have replied or after
// a timeout.
func Run(costReply <-chan def.Message, numberOfOnlineLifts *int) {
	assignmentQueue := make(map[order][]reply)

	var timeout = make(chan *order)
	//const timeoutDuration = 10 * time.Second
	const timeoutDuration = 1 * time.Nanosecond

	var newOrder order
	for {
		select {
		case message := <-costReply:
			newOrder = order{floor: message.Floor, button: message.Button}
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
					newOrder.timer.Reset(timeoutDuration)
				}
			} else {
				// If order not in queue at all, init order list with it
				newOrder.timer = time.NewTimer(timeoutDuration)
				assignmentQueue[newOrder] = []reply{newReply}
				go costTimer(&newOrder, timeout)
			}
			evaluateLists(&assignmentQueue, numberOfOnlineLifts, false)

		case <-timeout:
			log.Println(def.ClrR, "COST TIMED OUT!", def.ClrN)
			evaluateLists(&assignmentQueue, numberOfOnlineLifts, true)
		}
	}
}

// evaluateLists goes through the map of orders with associated costs,
// checks if any orders have received answers from all live lifts or the
// timer has timed out, and finds the best candidate for all such orders.
//The best candidate is added to the shared queue.
// This is very cryptic and ungood. // todo don't admit this
func evaluateLists(que *(map[order][]reply), numberOfOnlineLifts *int, orderTimedOut bool) {
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

func costTimer(newOrder *order, timeout chan *order) {
	<-newOrder.timer.C
	timeout <- newOrder
}

// --------------- METHODS FOR ORDER TYPE: ---------------

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
