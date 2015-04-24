package config

import (
	"time"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	ButtonUp int = iota
	ButtonDown
	ButtonIn //todo: Rename to ButtonInternal or something
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

// Constants for sending aliveMsg, and detecting deaths
const SpamInterval = 400 * time.Millisecond
const OnlineLiftResetTime = 2 * time.Second //todo: name?

// Constant setting the time the elevators wait befor an order times out
const OrderTime = 10 * time.Second //todo: name? //set this to 30 Seconds after done debuging

// Message kind constants
const (
	Alive int = iota + 1
	NewOrder
	CompleteOrder
	Cost
)

// Generic network message. No other messages are ever sent on the network.
type Message struct {
	Description int
	Floor       int
	Button      int
	Cost        int
	Addr        string `json:"-"`
}

type Keypress struct {
	Button int
	Floor  int
}

var OutgoingMsg = make(chan Message) // vurder buff
var SyncLightsChan = make(chan bool)

// Local address
var Laddr string
