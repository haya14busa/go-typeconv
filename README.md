## go-typeconv - Bring implicit type conversion into Go in a explicit way

[![Build Status](https://travis-ci.org/haya14busa/go-typeconv.svg?branch=master)](https://travis-ci.org/haya14busa/go-typeconv)
[![Go Report Card](https://goreportcard.com/badge/github.com/haya14busa/go-typeconv)](https://goreportcard.com/report/github.com/haya14busa/go-typeconv)
[![Coverage](https://codecov.io/gh/haya14busa/go-typeconv/branch/master/graph/badge.svg)](https://codecov.io/gh/haya14busa/go-typeconv)
[![LICENSE](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go doesn't have implicit type conversion.

> Unlike in C, in Go assignment between items of different type requires an explicit conversion.
> -- Type conversions https://tour.golang.org/basics/13

I like this design. Explicit is better than implicit.
In Go, almost all things are expressed explicitly.

However, sometimes... it's too tedious to fix type coversion errors by hand.
If a required type is `int64` and got `int` type, why not converting it automatically?
I'm tired of wrapping expressions with `int64()` or something here and there.

Here comes gotypeconv! gotypeconv takes source code, detects the type conversion errors and fixes them automatically by rewriting AST.

gotypeconv is like gofmt (it actually formats code as well), but it also fixes type conversions errors.

### Installation

```
go get -u github.com/haya14busa/go-typeconv/cmd/gotypeconv
```

### Usage example

#### ./testdata/tour.input.go

```go
// https://tour.golang.org/basics/13
package main

import (
	"fmt"
	"math"
)

func main() {
	var x, y int = 3, 4
	var f float64 = math.Sqrt(x*x + y*y)
	var z uint = f
	fmt.Println(x, y, z)
}
```

Above code has type conversion errors as follow.

```
$ go build testdata/tour.input.go
testdata/tour.input.go:11: cannot use x * x + y * y (type int) as type float64 in argument to math.Sqrt
testdata/tour.input.go:12: cannot use f (type float64) as type uint in assignment
```

gotypeconv can fix them automatically!

```
$ gotypeconv ./testdata/tour.input.go
// https://tour.golang.org/basics/13
package main

import (
        "fmt"
        "math"
)

func main() {
        var x, y int = 3, 4
        var f float64 = math.Sqrt(float64(x*x + y*y))
        var z uint = uint(f)
        fmt.Println(x, y, z)
}
```

gotypeconv also supports displaying diff (`-d` flag) and rewriting files in-place (`-w` flag) same as gofmt.

### More example

Go doesn't have overloading. https://golang.org/doc/faq#overloading
I like this design too.

However, sometimes... it's inconvenient.
For example, when you want `max` utility function, you may write something like this `func max(x int64, ys ...int64) int64`.
It works, but when you want to calculate max of given `int`s, you cannot ues this function unless wrapping them with `int64()`.

You also may start to write `func max(x int, ys ...int) int {`, and change type to int64 later.
Then, you need to wrap expressions with `int64()` here and there in this case as well.

Here comes gotypeconv, again!

```go
package main

import "fmt"

func main() {
	var (
		x int     = 1
		y int64   = 14
		z float64 = -1.4
	)

	var ans int = max(x, x+y, z)
	fmt.Println(ans)
}

func max(x int64, ys ...int64) int64 {
	for _, y := range ys {
		if y > x {
			x = y
		}
	}
	return x
}
```

Above code can be fixed gotypeconv. (`$ gotypeconv -d testdata/max.input.go`)

```diff
@@ -9,7 +9,7 @@
                z float64 = -1.4
        )

-       var ans int = max(x, x+y, z)
+       var ans int = int(max(int64(x), int64(x)+y, int64(z)))
        fmt.Println(ans)
 }
```

(I miss generics in this case... but gotypeconv can also solve the problem!)

#### Hou to Use in Vim

Use https://github.com/haya14busa/vim-gofmt with following sample config.

```vim
let g:gofmt_formatters = [
\   { 'cmd': 'gofmt', 'args': ['-s', '-w'] },
\   { 'cmd': 'goimports', 'args': ['-w'] },
\   { 'cmd': 'gotypeconv', 'args': ['-w'] },
\ ]
```

## :bird: Author
haya14busa (https://github.com/haya14busa)
