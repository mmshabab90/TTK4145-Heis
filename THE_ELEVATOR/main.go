package main

import (
	"./src/poller"
	"log"
)

var _ = log.Println

func main() {
	poller.Init()
}
