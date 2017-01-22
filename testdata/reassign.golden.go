package testdata

func f() {
	var (
		x float64 = 1
		y int     = 2
	)
	x = float64(y)
	_ = x
}
