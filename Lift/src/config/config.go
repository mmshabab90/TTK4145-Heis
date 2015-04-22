package config

import (
	"net"
	"strings"
	"time"
)

// Global system constants
const NumFloors = 4
const NumButtons = 3

const (
	ButtonUp int = iota
	ButtonDown
	ButtonIn
)

const MaxInt = int(^uint(0) >> 1)

// Move these out
var Outgoing = make(chan Message) // vurder buff
var SyncLightsChan = make(chan bool)

// Move into network module and pass to other modules
var Laddr *net.UDPAddr //Local address
