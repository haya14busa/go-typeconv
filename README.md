# go-typeconv - Bring implicit type conversion into Go in a explicit way

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
# command-line-arguments
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


## :bird: Author
haya14busa (https://github.com/haya14busa)
