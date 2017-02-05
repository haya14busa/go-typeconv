package main

import "fmt"

func main() {
	var a, b, c int = 1, 2, 3
	var x int64 = max(int64(a), int64(b), int64(c))
	fmt.Println(x)
}

func max(a int64, bs ...int64) int64 {
	x := 1
	return int64(x)
}
