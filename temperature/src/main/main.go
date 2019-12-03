package main

import (
    "fmt"
    "strings"
    "temperature"
)

var value float64
var unit string 
 
func main(){
    fmt.Println("Temperature? (seperate Value and Unit by space): ")
    fmt.Scanln(&value, &unit)
    
    unit = strings.ToUpper(unit[:1])
    
    switch unit {
        case "C":
            tempF := temperature.ConvertCtoF(value)
            fmt.Printf("%.2f F\n", tempF)
        case "F":
            c := temperature.ConvertFtoC(value)
            fmt.Printf("%.2f C\n", c)
        default:
            fmt.Println("conver temperature fail...")
    }
}
