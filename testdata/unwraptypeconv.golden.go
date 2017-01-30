package testdata

func f() {
	var x, y int = 1, 2
	var _ int = x
	var _ int = x + y
	var _ int = x + y
	funcarg(x)
	y = x
}

func funcarg(x int) {
}

func returnfn() (float64, int64) {
	var x float64 = 1
	var y int64 = 4
	return x, y
}
