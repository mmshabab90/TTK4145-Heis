package config

import (
	"net"
	"strings"
	"time"
)

// Global system constants
const NumButtons = 3
const NumFloors = 4

const (
	ButtonUp int = iota
	ButtonDown
	ButtonIn
)

const (
	DirDown int = iota - 1
	DirStop
	DirUp
)

const MaxInt = int(^uint(0) >> 1)

// Move these out
var Outgoing = make(chan Message) // vurder buff
var SyncLightsChan = make(chan bool)

// Move into network module and pass to other modules
var Laddr *net.UDPAddr //Local address
