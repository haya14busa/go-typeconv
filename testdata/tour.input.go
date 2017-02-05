// https://tour.golang.org/basics/13
package main

import (
	"fmt"
	"math"
)

func main() {
	var x, y int = 3, 4
	var f float64 = math.Sqrt(x*x + y*y)
	var z uint = f
	fmt.Println(x, y, z)
}
