package main

import (
	def "config"
	"errors"
	"fmt"
	"fsm"
	"hw"
	"log"
	"network"
	"os"
	"os/signal"
	"queue"
	"time"
)

var _ = log.Println
var _ = fmt.Println
var _ = errors.New

var onlineLifts = make(map[string]network.UdpConnection)

var deadChan = make(chan network.UdpConnection)
var costChan = make(chan def.Message)
var costTimeoutChan = make(chan order)

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

func main() {
	var floor int
	var err error
	floor, err = hw.Init()
	if err != nil {
		def.Restart.Run()
		log.Fatal(err)
	}

	e := fsm.Channels{
		NewOrder:     make(chan bool),
		FloorReached: make(chan int),
		MotorDir:     make(chan int, 10),
		FloorLamp:    make(chan int, 10),
		DoorLamp:     make(chan bool, 10),
	}
	fsm.Init(e, floor)

	network.Init()

	//handle ctrl+c
	safeKill() //bad name?

	liftAssigner(e.NewOrder)
	go poll(e)
	queue.Init(e.NewOrder)

	for { //nicer solution?
		time.Sleep(100 * time.Second)
	}
}

func poll(e fsm.Channels) {
	buttonChan := pollButtons()
	floorChan := pollFloors()

	for {
		select {
		case keypress := <-buttonChan:
			switch keypress.Button {
			case def.BtnInside:
				queue.AddLocalOrder(keypress.Floor, keypress.Button)
			case def.BtnUp, def.BtnDown:
				def.OutgoingMsg <- def.Message{
					Category: def.NewOrder,
					Floor:    keypress.Floor,
					Button:   keypress.Button}
			}
		case floor := <-floorChan:
			e.FloorReached <- floor
		case udpMessage := <-network.ReceiveChan:
			handleMessage(network.ParseMessage(udpMessage))
		case connection := <-deadChan:
			handleDeadLift(connection.Addr)
		case order := <-queue.OrderTimeoutChan:
			fmt.Println("Order timeout, I can do it myself!")
			queue.RemoveRemoteOrdersAt(order.Floor)
			queue.AddRemoteOrder(order.Floor, order.Button, def.Laddr)
		case dir := <-e.MotorDir:
			hw.SetMotorDirection(dir)
		case floor := <-e.FloorLamp:
			hw.SetFloorLamp(floor)
		case value := <-e.DoorLamp:
			hw.SetDoorOpenLamp(value)
		}
	}
}

func pollButtons() <-chan def.Keypress {
	c := make(chan def.Keypress)

	go func() {
		var buttonState [def.NumFloors][def.NumButtons]bool

		for {
			for f := 0; f < def.NumFloors; f++ {
				for b := 0; b < def.NumButtons; b++ {
					if (f == 0 && b == def.BtnDown) ||
						(f == def.NumFloors-1 && b == def.BtnUp) {
						continue
					}
					if hw.ReadButton(f, b) {
						if !buttonState[f][b] {
							c <- def.Keypress{Button: b, Floor: f}
						}
						buttonState[f][b] = true
					} else {
						buttonState[f][b] = false
					}
				}
			}
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

func pollFloors() <-chan int {
	c := make(chan int)

	go func() {
		oldFloor := hw.Floor()

		for {
			newFloor := hw.Floor()
			if newFloor != oldFloor && newFloor != -1 {
				c <- newFloor
			}
			oldFloor = newFloor
			time.Sleep(time.Millisecond)
		}
	}()

	return c
}

// consider moving each case into a function
func handleMessage(message def.Message) {
	const aliveTimeout = 2 * time.Second

	switch message.Category {
	case def.Alive:
		if connection, exist := onlineLifts[message.Addr]; exist {
			connection.Timer.Reset(aliveTimeout)
		} else {
			newConnection := network.UdpConnection{message.Addr, time.NewTimer(aliveTimeout)}
			onlineLifts[message.Addr] = newConnection
			go connectionTimer(&newConnection)
		}
	case def.NewOrder:
		// log.Printf("handleMessage(): NewOrder message: f=%d b=%d from lift %s\n",
		//	message.Floor+1, message.Button, message.Addr[12:15])

		cost := queue.CalculateCost(message.Floor, message.Button, fsm.Floor(), hw.Floor(), fsm.Direction())

		costMessage := def.Message{
			Category: def.Cost,
			Floor:    message.Floor,
			Button:   message.Button,
			Cost:     cost}
		// log.Printf("handleMessage(): NewOrder sends cost message: f=%d b=%d (with cost %d) from me\n", costMessage.Floor+1, costMessage.Button, costMessage.Cost)
		def.OutgoingMsg <- costMessage
	case def.CompleteOrder:
		queue.RemoveRemoteOrdersAt(message.Floor)
	case def.Cost:
		// log.Printf("handleMessage(): Cost message: f=%d b=%d with cost %d from lift %s\n", message.Floor+1, message.Button, message.Cost, message.Addr[12:15])
		costChan <- message
	}
}

func handleDeadLift(deadAddr string) {
	fmt.Printf("Connection to IP %s is dead!\n", deadAddr)
	delete(onlineLifts, deadAddr)
	queue.ReassignOrders(deadAddr)
}

func connectionTimer(connection *network.UdpConnection) {
	<-connection.Timer.C
	deadChan <- *connection
}

func costTimer(newOrder *order) {
	<-newOrder.timer.C
	costTimeoutChan <- *newOrder
}

//bad variable names am bad!
func liftAssigner(newOrderChan chan bool) {
	// collect cost values from all lifts
	// decide which lift gets the order when all lifts
	// in alive-list have answered or after a timeout
	// either send the decision on network or pray that all
	// lifts make the same choice every time
	go func() {
		assignmentQueue := make(map[order][]reply)
		const costReplyTimeout = 10 * time.Second
		var newOrder order
		for {
			select {
			case message := <-costChan:
				//newKey, newReply := split(message) //newKey is actually the worst name ever
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
				evaluateLists(&assignmentQueue)

			case newOrder := <-costTimeoutChan:
				fmt.Printf("\n COST TIMED OUT!\n\n")
				newOrder.setTimeout(true)
				evaluateLists(&assignmentQueue)
			}
		}
	}()
}

/*func split(m def.Message) (order, reply) {
	return order{floor: m.Floor, button: m.Button}, reply{cost: m.Cost, lift: m.Addr}
}*/ //isn't used anymore

func getReply(m def.Message) reply {
	return reply{cost: m.Cost, lift: m.Addr}
}

// evaluateLists goes through the map of orders with associated costs, checks
// if any orders have received answers from all live lifts, and finds the
// the best candidate for all such orders. The best candidate is added to the
// shared queue.
// This is very cryptic and ungood. // todo don't admit this
func evaluateLists(que *(map[order][]reply)) {
	const maxInt = int(^uint(0) >> 1)
	// Loop through all lists
	fmt.Printf("Lists: ")
	fmt.Println(*que)
	for order, replyList := range *que {
		// Check if the list is complete
		if len(replyList) == len(onlineLifts) || order.timeout {
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

			// Assign order key to lift
			queue.AddRemoteOrder(order.floor, order.button, lowAddr)
			//queue.PrintQueues()

			// Empty list and stop timer
			order.timer.Stop()
			delete(*que, order)
		}
	}
}

//safeKill gets the motor to stop when the program is killed with ctrl+c
func safeKill() {
	var c = make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		hw.SetMotorDirection(def.DirStop)
		//restart.Run() //vil vi restarte med ctrl+c?
		log.Fatal("[FATAL]\tUser terminated program")
	}()
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
