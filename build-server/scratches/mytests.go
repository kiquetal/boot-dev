package main

import (
	"fmt"
)

func main(){

	bulkSend(10)

}

func bulkSend(numMessages int) float64 {
	var total float64
	for i:=0;i <numMessages;i++ {
		result := float64(i) / 1000
		fmt.Printf("Hello, World! %v\n", result)
	//	total += 1.0 + float64(i/100)
	}
	return total
}
