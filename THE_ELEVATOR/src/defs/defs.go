package defs

import "net"

const NumButtons = 3
const NumFloors = 4
const (
	ButtonCallUp int = iota
	ButtonCallDown
	ButtonCommand
)
const (
	DirnDown int = iota - 1
	DirnStop
	DirnUp
)
const (
	Alive int = iota
	NewOrder
	CompleteOrder
	Cost
)

type Message struct {
	Kind   int
	Floor  int
	Button int
	Cost   int
	Addr   string `json:"-"`
}

var MessageChan = make(chan Message)

var Laddr *net.UDPAddr //Local address
