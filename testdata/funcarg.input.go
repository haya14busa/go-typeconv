package testdata

func f() {
	x := 1
	funcarg(x)
	funcarg(x + 1)
	funcarg(int(1) + 1)
	funcarg2(x)
}

func funcarg(x float64) {
}

func funcarg2(x int64) {
}
