package testdata

func f() {
	var (
		x float64 = 1
		y int     = 2
	)
	x = y
	_ = x
	y, x, x = 1, y, y
}
