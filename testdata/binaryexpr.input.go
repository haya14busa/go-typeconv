package testdata

func f() {
	var (
		x int     = 1
		y float64 = 2
	)

	z := x * y
	_ = z
	_ = x == y
	_ = y == x
}
