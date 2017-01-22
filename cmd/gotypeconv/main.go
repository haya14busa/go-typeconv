package main

import (
	"flag"
	"fmt"
	"go/printer"
	"io"
	"os"

	typeconv "github.com/haya14busa/go-typeconv"

	"golang.org/x/tools/go/loader"
)

func main() {
	flag.Parse()
	if err := run(os.Stderr, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(stderr io.Writer, args []string) error {
	prog, typeErrs, err := typeconv.Load(loader.Config{}, args)
	if err != nil {
		return err
	}
	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			if err := typeconv.RewriteFile(prog.Fset, f, typeErrs); err != nil {
				return err
			}
			if err := printer.Fprint(os.Stdout, prog.Fset, f); err != nil {
				return err
			}
		}
	}
	return nil
}
