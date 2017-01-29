package typeconv

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/diff"

	"golang.org/x/tools/go/loader"
)

func TestLoad(t *testing.T) {
	files, err := filepath.Glob("testdata/*.input.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range files {
		prog, typeErrs, err := Load(loader.Config{}, []string{file})
		if err != nil {
			t.Error(err)
		}
		if got := len(prog.InitialPackages()); got != 1 {
			t.Errorf("len(prog.InitialPackages()) == %v, want 1", got)
		}
		pkg := prog.InitialPackages()[0]
		if got := len(pkg.Files); got != 1 {
			t.Errorf("len(pkg.Files) == %v, want 1", got)
		}
		f := pkg.Files[0]
		if got := prog.Fset.File(f.Pos()).Name(); got != file {
			t.Errorf("filename: got %v, want %v", got, file)
		}
		if len(typeErrs) == 0 {
			t.Errorf("len(typeErrs) is empty, expect errors")
		}
	}
}

func TestRewriteProgram(t *testing.T) {
	files, err := filepath.Glob("testdata/*.input.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, fname := range files {
		input := fname
		golden := strings.Replace(input, "input.go", "golden.go", 1)

		prog, typeErrs, err := Load(loader.Config{}, []string{input})
		if err != nil {
			t.Fatalf("%s: %v", fname, err)
		}
		if err := RewriteProgam(prog, typeErrs); err != nil {
			t.Fatalf("%s: %v", fname, err)
		}

		f := prog.InitialPackages()[0].Files[0]
		buf := new(bytes.Buffer)
		if err := format.Node(buf, prog.Fset, f); err != nil {
			t.Fatalf("%s: %v", fname, err)
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

		if d := diff.Diff(buf.String(), string(b)); d != "" {
			t.Errorf("%s: diff: (-got +want):\n%s", fname, d)
		}
	}
}
