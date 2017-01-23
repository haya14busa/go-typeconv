// Package typeconv provides missing implicit type conversion in Go by
// rewriting AST.
package typeconv

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/loader"
)

// Load creates the initial packages specified by conf along with slice of
// types.Error.
func Load(conf loader.Config, args []string) (*loader.Program, []types.Error, error) {
	conf.AllowErrors = true
	conf.ParserMode = parser.ParseComments
	var typeErrs []types.Error
	typeErrFn := conf.TypeChecker.Error
	conf.TypeChecker.Error = func(err error) {
		if err, ok := err.(types.Error); ok {
			typeErrs = append(typeErrs, err)
		}
		if typeErrFn != nil {
			typeErrFn(err)
		}
	}
	if _, err := conf.FromArgs(args, true); err != nil {
		return nil, nil, err
	}
	prog, err := conf.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load program: %v", err)
	}
	return prog, typeErrs, nil
}

// RewriteFile rewrites ast.File to fix type conversion errors.
func RewriteFile(fset *token.FileSet, f *ast.File, typeErrs []types.Error) error {
	filename := fset.File(f.Pos()).Name()

	for _, e := range typeErrs {
		// fmt.Println(e) // debug
		if filename != e.Fset.File(e.Pos).Name() {
			continue
		}

		path, exact := astutil.PathEnclosingInterval(f, e.Pos, e.Pos)
		if !exact {
			return fmt.Errorf("cannot get exact node position for type error: %v", e)
		}

		terr := NewTypeErr(e)
		if terr == nil {
			continue
		}

		switch terr := terr.(type) {
		case *ErrVarDecl:
			if err := rewriteErrVarDecl(path, terr); err != nil {
				return err
			}
		case *ErrFuncArg:
			if err := rewriteErrFuncArg(path, terr); err != nil {
				return err
			}
		}
	}

	return nil
}

func rewriteErrVarDecl(path []ast.Node, terr *ErrVarDecl) error {
	for i := range path {
		if i+1 >= len(path) {
			break
		}
		child, parent := path[i], path[i+1]
		if valuespec, ok := parent.(*ast.ValueSpec); ok {
			idx := -1
			for i, value := range valuespec.Values {
				if value == child {
					idx = i
					break
				}
			}
			if idx == -1 {
				return fmt.Errorf("cannot find expected value: %v", child)
			}
			// TODO(haya14busa): check terr.ValueType is convertible to terr.NameType
			valuespec.Values[idx] = &ast.CallExpr{
				Fun:  ast.NewIdent(terr.NameType),
				Args: []ast.Expr{valuespec.Values[idx]},
			}
		}
	}
	return nil
}

func rewriteErrFuncArg(path []ast.Node, terr *ErrFuncArg) error {
	for i := range path {
		if i+1 >= len(path) {
			break
		}
		child, parent := path[i], path[i+1]
		if call, ok := parent.(*ast.CallExpr); ok {
			idx := -1
			for i, arg := range call.Args {
				if arg == child {
					idx = i
				}
			}
			if idx == -1 {
				return fmt.Errorf("cannot find expected value: %v", child)
			}
			// TODO(haya14busa): check terr.ArgType is convertible to terr.ParamType
			call.Args[idx] = &ast.CallExpr{
				Fun:  ast.NewIdent(terr.ParamType),
				Args: []ast.Expr{call.Args[idx]},
			}
		}
	}
	return nil
}
