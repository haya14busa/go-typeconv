package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	typeconv "github.com/haya14busa/go-typeconv"

	"golang.org/x/tools/go/loader"
)

type option struct {
	write  bool
	doDiff bool
	rules  strslice
}

func main() {
	opt := &option{}
	flag.BoolVar(&opt.write, "w", false, "write result to (source) file instead of stdout")
	flag.BoolVar(&opt.doDiff, "d", false, "display diffs instead of rewriting files")
	flag.Var(&opt.rules, "r", "type conversion rules currently just for type conversion of binary expression (e.g., 'int -> uint32')")
	flag.Parse()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if err := run(out, flag.Args(), opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string, opt *option) error {
	if err := addRules(opt.rules); err != nil {
		return err
	}
	prog, typeErrs, err := typeconv.Load(loader.Config{}, args)
	if err != nil {
		return err
	}
	if err := typeconv.RewriteProgam(prog, typeErrs); err != nil {
		return err
	}
	for _, pkg := range prog.InitialPackages() {
		for _, f := range pkg.Files {
			printFile(w, opt, prog, f)
		}
	}
	return nil
}

func printFile(w io.Writer, opt *option, prog *loader.Program, f *ast.File) error {
	filename := prog.Fset.File(f.Pos()).Name()
	buf := new(bytes.Buffer)
	if err := format.Node(buf, prog.Fset, f); err != nil {
		return err
	}
	res := buf.Bytes()
	in, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer in.Close()
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
	return nil
}

func addRules(rules []string) error {
	for _, r := range rules {
		f := strings.Split(r, "->")
		if len(f) != 2 {
			return fmt.Errorf("type conversion rule must be the form 'from -> to': %v", r)
		}
		from, to := strings.TrimSpace(f[0]), strings.TrimSpace(f[1])
		typeconv.DefaultRule.Add(from, to)
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

type strslice []string

func (ss *strslice) String() string {
	return fmt.Sprintf("%v", *ss)
}

func (ss *strslice) Set(value string) error {
	*ss = append(*ss, value)
	return nil
}
