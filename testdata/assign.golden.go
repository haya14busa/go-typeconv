package testdata

func f() {
	var x, y int = 3, 4
	var _ float64 = float64(x*x + y*y)
	var _ uint = uint(x)
}
