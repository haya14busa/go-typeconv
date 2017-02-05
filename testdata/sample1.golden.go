package main

import "fmt"

func main() {
	var x int
	var a, b, c int = 1, 2, 3
	x = int(max(int64(a), int64(b), int64(c)))
	fmt.Println(x)
}

func max(xs ...int64) int64 {
	x := 1
	return int64(x)
}
