package main
import "fmt"

func multiply(a,b int) int {

  return a*b

}

func sum(a,b int ) int {
return a+b
}

func apply(f func(a,b int) int) func(a,b int ) int { 
 return func(a,b int ) int {
    return  f(a,b)
 }
 
}
func main() {
fmt.Println(apply(multiply)(4,3))
}
