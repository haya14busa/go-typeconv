package testdata

func f() {
	x := 1
	funcarg(float64(x))
	funcarg(float64(x + 1))
}

func funcarg(x float64) {
}
