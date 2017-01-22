package main

import (
	"flag"
	"fmt"
	"go/format"
	"go/printer"
	"io"
	"os"

	typeconv "github.com/haya14busa/go-typeconv"

	"golang.org/x/tools/go/loader"
)

type option struct {
	write bool
}

func main() {
	opt := &option{}
	flag.BoolVar(&opt.write, "w", false, "write result to (source) file instead of stdout")
	flag.Parse()
	if err := run(os.Stderr, flag.Args(), opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(stderr io.Writer, args []string, opt *option) error {
	prog, typeErrs, err := typeconv.Load(loader.Config{}, args)
	if err != nil {
		return err
	}
	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			filename := prog.Fset.File(f.Pos()).Name()
			if err := typeconv.RewriteFile(prog.Fset, f, typeErrs); err != nil {
				return err
			}
			if opt.write {
				_ = format.Node
				fh, err := os.Create(filename)
				if err != nil {
					return err
				}
				format.Node(fh, prog.Fset, f)
				fh.Close()
			} else {
				if err := printer.Fprint(os.Stdout, prog.Fset, f); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
