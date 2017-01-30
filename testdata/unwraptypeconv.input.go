package testdata

func f() {
	var x, y int = 1, 2
	var _ int = int64(x)
	var _ int = int64(x) + y
	var _ int = x + int64(y)
	funcarg(int64(x))
	y = int64(x)
}

func funcarg(x int) {
}

func returnfn() (float64, int64) {
	var x float64 = 1
	var y int64 = 4
	return int(x), float64(y)
}
