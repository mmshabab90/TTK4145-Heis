package main

import (
	. "fmt"
	"runtime"
)

var i = 0

const million = 1000000

func increase(c chan<- bool) {
	for j := 0; j < million; j++ {
		i++
	}
	c <- true
}

func decrease(c chan<- bool) {
	for j := 0; j < million; j++ {
		i--
	}
	c <- true
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	incrDone := make(chan bool)
	decrDone := make(chan bool)

	go increase(incrDone)
	go decrease(decrDone)

	<-incrDone
	<-decrDone

	Println(i)
}
