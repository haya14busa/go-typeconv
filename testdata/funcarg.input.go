package testdata

func f() {
	x := 1
	funcarg(x)
	funcarg(x + 1)
	funcarg(int(1) + 1)
}

func funcarg(x float64) {
}
