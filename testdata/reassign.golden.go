package testdata

func f() {
	var (
		x float64 = 1
		y int     = 2
	)
	x = float64(y)
	_ = x
	y, x, x = 1, float64(y), float64(y)
}
