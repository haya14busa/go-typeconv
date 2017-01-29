package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	ddiff "github.com/kylelemons/godebug/diff"
)

func TestRun_package(t *testing.T) {
	opt := &option{}
	buf := new(bytes.Buffer)
	if err := run(buf, []string{"github.com/haya14busa/go-typeconv"}, opt); err != nil {
		t.Fatal(err)
	}
	if buf.Len() == 0 {
		t.Error("run: output is empty")
	}
}

func TestRun_testdata(t *testing.T) {
	opt := &option{}
	files, err := filepath.Glob("../../testdata/*.input.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, fname := range files {
		input := fname
		golden := strings.Replace(input, "input.go", "golden.go", 1)
		buf := new(bytes.Buffer)
		if err := run(buf, []string{input}, opt); err != nil {
			t.Fatal(err)
		}
		gf, err := os.Open(golden)
		if err != nil {
			t.Fatalf("%s: %v", fname, err)
		}
		defer gf.Close()
		b, err := ioutil.ReadAll(gf)
		if err != nil {
			t.Fatalf("%s: %v", fname, err)
		}

		if d := ddiff.Diff(buf.String(), string(b)); d != "" {
			t.Errorf("%s: diff: (-got +want):\n%s", fname, d)
		}
	}
}

func BenchmarkRun(b *testing.B) {
	opt := &option{}
	for i := 0; i < b.N; i++ {
		buf := new(bytes.Buffer)
		if err := run(buf, []string{"github.com/haya14busa/go-typeconv"}, opt); err != nil {
			b.Fatal(err)
		}
	}
}
