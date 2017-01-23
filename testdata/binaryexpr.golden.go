package testdata

func f() {
	var (
		x int     = 1
		y float64 = 2
	)

	z := float64(x) * y
	_ = z
	_ = float64(x) == y
	_ = y == float64(x)
}
