package main

import(
	. "fmt"
	"runtime"
	// "time"
)

var(
	i int = 0
) 

func increase(channel chan int, doneChannel chan bool){
	for j := 0; j < 1000001; j++{
		i = <- channel
		i++
		channel <- i
	}
	doneChannel <- true
}

func decrease(channel chan int, doneChannel chan bool){
	for j := 0; j < 1000000; j++{
		i = <- channel
		i--
		channel <- i
	}
	doneChannel <- true
}

func main() {
	
	channel := make(chan int, 1);
	doneChannel := make(chan bool, 1);

	channel <- i

	runtime.GOMAXPROCS(runtime.NumCPU())

	go increase(channel, doneChannel)

	go decrease(channel, doneChannel)

	<- doneChannel
	<- doneChannel
	
	Println(i)
}
