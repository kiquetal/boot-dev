package main

import "fmt"

func main() {

   fmt.Println(getMessageCosts([]string{"hola","queta"}))
}

func getMessageCosts(messages []string) []float64 {

	costs := make([]float64,len(messages))
	for i,_ := range costs {
	  costs[i]=float64(i) * 0.01
	}
	return costs

}
