package main

import(
	. "fmt"
	"runtime"
	"time"
)

var(
	i int = 0
) 

func increase() {
	for j := 0; j < 1000000; j++ {
		i++
	}
}

func decrease() {
	for j := 0; j < 1000000; j++ {
		i--
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go increase()
	go decrease()

	time.Sleep(100*time.Millisecond)
	
	Println(i)
}