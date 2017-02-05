package main

import "fmt"

func main() {
	var (
		x int     = 1
		y int64   = 14
		z float64 = -1.4
	)

	var ans int = max(x, x+y, z)
	fmt.Println(ans)
}

func max(x int64, ys ...int64) int64 {
	for _, y := range ys {
		if y > x {
			x = y
		}
	}
	return x
}
