// Shitty file name
package network

import (
	def "../config"
	"fmt"
	"time"
)

type reply struct {
	cost int
	lift string
}
type order struct {
	floor  int
	button int
	// timeout bool
	// timer   *time.Timer
}

type RemoteOrder struct {
	Floor int
	Button int
	Addr string
}

var deadLift = make(chan string)

// liftAssigner collects cost values from all lifts, decides which lift gets
// the order when all lifts in alive-list have answered or after a timeout.
func liftAssigner(addRemoteOrder chan<- RemoteOrder, costMessage <- chan Message, onlineLifts map[string]*time.Timer) {
	assignmentQueue := make(map[order][]reply)
	for {
		message := <-costMessage

		newOrder, newReply := split(message)
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
				//newOrder.timer.Reset(10 * time.Second)
			}
		} else {
			// If order not in queue at all, init order list with it
			assignmentQueue[newOrder] = []reply{newReply}
			// newOrder.timer = time.NewTimer(10 * time.Second)
			// go orderTimer(&newOrder)
		}
		evaluateLists(assignmentQueue, onlineLifts, addRemoteOrder)
		/*case newOrder := <-orderTimeoutChan:
		fmt.Printf("\n\n ORDER TIMED OUT!\n")
		// newOrder.timeout = true
		evaluateLists(assignmentQueue)*/
	}
}

func split(m Message) (order, reply) {
	return order{floor: m.Floor, button: m.Button}, reply{cost: m.Cost, lift: m.Addr}
}

// evaluateLists goes through the map of orders with associated costs, checks
// if any orders have received answers from all live lifts, and finds the
// the best candidate for all such orders. The best candidate is added to the
// shared queue.
// This is very cryptic and ungood.
func evaluateLists(que map[order][]reply, onlineLifts map[string]*time.Timer, addRemoteOrder chan<- RemoteOrder) {
	// Loop thru all lists
	fmt.Printf("Lists: ")
	fmt.Println(que)
	for key, replyList := range que {
		// Check if the list is complete
		if len(replyList) == len(onlineLifts) /*|| key.timeout*/ {
			fmt.Printf("Laddr = %v\n", def.Laddr)
			var (
				lowCost = def.MaxInt
				lowAddr string
			)
			// Loop thru costs in each complete list
			for _, reply := range replyList {
				// hvis ny bedre enn gammel: best = ny
				// hvis ny og gammel like bra og best: ny = lavest ip
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
			fmt.Printf("Lift %s won order f=%d b=%d\n", lowAddr[12:15], key.floor+1, key.button)
			// Assign order key to lift
			addRemoteOrder <- RemoteOrder{key.floor, key.button, lowAddr}
			delete(que, key)
		}
	}
}
