package main

import "fmt"

func main() {
	var x int
	var a, b, c int = 1, 2, 3
	x = max(a, b, c)
	fmt.Println(x)
}

func max(xs ...int64) int64 {
	x := 1
	return x
}
