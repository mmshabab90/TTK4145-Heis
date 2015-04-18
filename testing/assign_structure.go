package main

import (
	"fmt"
	"time"
)

type cost struct {
	cost int
	lift string
}
type order struct {
	floor int
	button int
}

func main() {
	go func() {
		// This is a map that indexes on a struct and returns a slice:
		assignmentQueue := make(map[order][]cost)

		assignmentQueue[order{2,0}] = append(assignmentQueue[order{2,0}], cost{cost:13, lift:"192.168.0.1"})

		fmt.Println(assignmentQueue)
	}()

	time.Sleep(time.Second)
}