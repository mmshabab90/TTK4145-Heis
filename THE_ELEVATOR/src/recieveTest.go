package main

import (
	"./network"
	//"encoding/json"
	//"./defs"
	//"fmt"
)

func main() {
	network.Init()
	for {
		m := <- network.ReceiveChan
		//var a defs.Message
		
		network.PrintMessage(network.ParseMessage(m))
		/*if err := json.Unmarshal(m.Data, &a); err != nil {
			fmt.Printf("json.Unmarshal error: %s\n", err)
		}
		fmt.Printf("\n----------Message start----------\n")
		fmt.Printf("Kind: %d\n", a.Kind)
		fmt.Printf("Floor: %d\n", a.Floor)
		fmt.Printf("Button: %d\n", a.Button)
		fmt.Printf("Cost: %d\n", a.Cost)
		fmt.Printf("Addr: %s\n", a.Addr)
		fmt.Println("----------Message   end----------")
		*/
	}
}
