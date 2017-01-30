package testdata

func f() {
	x := 1
	funcarg(float64(x))
	funcarg(float64(x + 1))
	funcarg(float64(int(1) + 1))
	funcarg2(int64(x))
}

func funcarg(x float64) {
}

func funcarg2(x int64) {
}
