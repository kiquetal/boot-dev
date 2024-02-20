package main
import "fmt"
type cost struct {
	day   int
	value float64
}

func main() {
	
	getCostsByDay2([]cost{
			 {0, 4.0},
    			 {1, 2.1},
    			 {5, 2.5},
    			{1, 3.1},
	})

}

func getCostsByDay2(costs []cost) []float64 {
	var maxDay int
	fmt.Printf("%v\n",costs)
	var saveIndx [] int
	for _, v := range costs {

		day,_ := v.day, v.value
	
		if day > maxDay {
			maxDay = day
		}
		found := false
		for _,v := range saveIndx {
                      if day == v {
				found = true
				break;
		      }

		}
		if (!found) {
			saveIndx=append(saveIndx,day)
		}
	}
        fmt.Printf("%v\n",saveIndx)	
	slicedArray :=make([]cost,(maxDay+1))	
	for i:=0;i<len(saveIndx);i++ {
		fmt.Printf("mi indice %v\n---\n",saveIndx[i])
              for _,v2 :=range costs {
			if saveIndx[i]== v2.day {
				fmt.Printf("value is %v\n",v2.value)
				slicedArray[saveIndx[i]]=cost{
					day:saveIndx[i],
					value:float64(v2.value) + float64(slicedArray[saveIndx[i]].value),
				}
			}
		  }


	      }
      	fmt.Printf("%v\n",slicedArray)
        toReturnCost:=make([]float64,maxDay+1)
	for i,c := range slicedArray {
	      toReturnCost[i]=float64(c.value)
	}
	fmt.Printf("result %v",toReturnCost)
	return nil
}

