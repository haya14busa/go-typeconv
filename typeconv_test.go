package typeconv

import (
	"path/filepath"
	"testing"

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
