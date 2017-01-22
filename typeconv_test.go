package typeconv

import (
	"bytes"
	"go/printer"
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestRewriteFile_assign(t *testing.T) {
	prog, typeErrs, err := Load(loader.Config{}, []string{"testdata/assign.input.go"})
	if err != nil {
		t.Fatal(err)
	}
	f := prog.InitialPackages()[0].Files[0]
	if err := RewriteFile(prog.Fset, f, typeErrs); err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	if err := printer.Fprint(buf, prog.Fset, f); err != nil {
		t.Fatal(err)
	}
	gf, err := os.Open("testdata/assign.golden.go")
	if err != nil {
		t.Fatal(err)
	}
	defer gf.Close()
	b, err := ioutil.ReadAll(gf)
	if err != nil {
		t.Fatal(err)
	}

	if d := diff.Diff(buf.String(), string(b)); d != "" {
		t.Errorf("diff: (-got +want):\n%s", d)
	}

}
