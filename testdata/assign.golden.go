package testdata

func f() {
	var x, y int = 3, 4
	var _ float64 = float64(x*x + y*y)
	var _ uint = uint(x)
	var _ uint = uint(int(1))
	var _ uint = uint(int(1) + int(1))
	// cannot convert
	var _ int = "string"
}
