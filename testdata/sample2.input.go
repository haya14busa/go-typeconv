package main

import "fmt"

func main() {
	var a, b, c int = 1, 2, 3
	var x int64 = int(max(a, int64(b), c))
	fmt.Println(x)
}

func max(a int64, bs ...int64) int64 {
	x := 1
	return x
}
