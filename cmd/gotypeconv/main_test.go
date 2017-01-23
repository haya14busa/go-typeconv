package main

import (
	"bytes"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	opt := &option{}
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		if err := run(buf, []string{"github.com/haya14busa/go-typeconv"}, opt); err != nil {
			b.Fatal(err)
		}
	}
}
