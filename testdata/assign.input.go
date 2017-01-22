package testdata

func f() {
	var x, y int = 3, 4
	var _ float64 = x*x + y*y
	var _ uint = x
	var _ uint = int(1)
	var _ uint = int(1) + int(1)
}
