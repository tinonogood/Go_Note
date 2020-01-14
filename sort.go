package main

import (
	"fmt"
)


func main() {
	strArr := []string{"a", "bb", "ccc", "eeeee", "dddd", "ffffff"}	
	numArr := []int{-1, 0 ,100, 10, -10000}
	
	sortStrArr(strArr)
	fmt.Println(strArr)
	
	sortNumArr(numArr)
	fmt.Println(numArr)
	
}

func sortStrArr(strArr []string){
	var(
		n = len(strArr)
		sorted = false
	)
	for !sorted{
		swapped := false
		for i:=0; i<n-1; i++{
			if len(strArr[i]) < len(strArr[i+1]){
				strArr[i+1], strArr[i] = strArr[i], strArr[i+1]
				swapped = true
			}
		}
		if !swapped {
			sorted=true
		}
		n = n-1
	}
}
	
func sortNumArr(numArr []int) {
    var n = len(numArr)
    for i := 1; i < n; i++ {
        j := i
        for j > 0 {
            if numArr[j-1] > numArr[j] {
                numArr[j-1], numArr[j] = numArr[j], numArr[j-1]
            }
            j = j - 1
        }
    }
}
