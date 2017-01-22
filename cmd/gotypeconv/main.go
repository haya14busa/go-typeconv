package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	typeconv "github.com/haya14busa/go-typeconv"

	"golang.org/x/tools/go/loader"
)

type option struct {
	write  bool
	doDiff bool
}

func main() {
	opt := &option{}
	flag.BoolVar(&opt.write, "w", false, "write result to (source) file instead of stdout")
	flag.BoolVar(&opt.doDiff, "d", false, "display diffs instead of rewriting files")
	flag.Parse()
	if err := run(os.Stdout, flag.Args(), opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(w io.Writer, args []string, opt *option) error {
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
			buf := new(bytes.Buffer)
			if err := format.Node(buf, prog.Fset, f); err != nil {
				return err
			}
			res := buf.Bytes()
			in, err := os.Open(filename)
			if err != nil {
				return err
			}
			src, err := ioutil.ReadAll(in)
			if err != nil {
				return err
			}
			if !bytes.Equal(src, res) {
				if opt.write {
					fh, err := os.Create(filename)
					if err != nil {
						return err
					}
					fh.Write(res)
					fh.Close()
				}
				if opt.doDiff {
					data, err := diff(src, res)
					if err != nil {
						return fmt.Errorf("computing diff: %s", err)
					}
					fmt.Fprintf(w, "diff %s gotypeconv/%s\n", filename, filename)
					w.Write(data)
				}
			}
			if !opt.write && !opt.doDiff {
				w.Write(res)
			}
		}
	}
	return nil
}

// copied and modified from $GOPATH/src/github.com/golang/go/src/cmd/gofmt/gofmt.go
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "gotypeconv")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "gotypeconv")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return

}
