package config

import (
	"time"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	BtnUp int = iota
	BtnDown
	BtnInside
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

// Local IP address
var Laddr string

//constants setting the time the elevators wait before, the cost and the order times out
const CostTime = 10 * time.Second // todo rename

//message kind constants
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

type Keypress struct {
	Button int
	Floor  int
}

// Generic network message. No other messages are ever sent on the network.
type Message struct {
	Category int
	Floor    int
	Button   int
	Cost     int
	Addr     string `json:"-"`
}

var OutgoingMsg = make(chan Message, 10) // vurder buff //todo: rename outgoing/outbox (and move?)
var SyncLightsChan = make(chan bool)     // todo move!
