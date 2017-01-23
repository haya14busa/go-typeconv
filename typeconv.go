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
func RewriteFile(fset *token.FileSet, f *ast.File, pkg *loader.PackageInfo, typeErrs []types.Error) error {
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
			if err := rewriteErrVarDecl(path, pkg, terr); err != nil {
				return err
			}
		case *ErrFuncArg:
			if err := rewriteErrFuncArg(path, pkg, terr); err != nil {
				return err
			}
		case *ErrAssign:
			if err := rewriteErrAssign(path, pkg, terr); err != nil {
				return err
			}
		case *ErrMismatched:
			if err := rewriteErrMismatched(path, pkg, terr); err != nil {
				return err
			}
		case *ErrReturn:
			if err := rewriteErrReturn(path, pkg, terr); err != nil {
				return err
			}
		}
	}

	return nil
}

func rewriteErrVarDecl(path []ast.Node, pkg *loader.PackageInfo, terr *ErrVarDecl) error {
	for i := range path {
		if i+1 >= len(path) {
			break
		}
		child, parent := path[i], path[i+1]
		if valuespec, ok := parent.(*ast.ValueSpec); ok {
			if ok := checkConvertibleErrVarDecl(terr, valuespec, child, pkg.Info); !ok {
				continue
			}
			idx := -1
			for i, value := range valuespec.Values {
				if value == child {
					idx = i
					break
				}
			}
			if idx == -1 {
				return nil
			}
			valuespec.Values[idx] = &ast.CallExpr{
				Fun:  ast.NewIdent(terr.NameType),
				Args: []ast.Expr{valuespec.Values[idx]},
			}
			break
		}
	}
	return nil
}

// checkConvertibleErrVarDecl checks child type is convertible to parent type.
// In fact, type error message seemes already covers this check... but leave it
// for just in case.
// e.g. `cannot convert "string" (untyped string constant) to int`
func checkConvertibleErrVarDecl(terr *ErrVarDecl, parent *ast.ValueSpec, child ast.Node, typeinfo types.Info) bool {
	parentExpr, ok := parent.Type.(ast.Expr)
	if !ok {
		return false
	}
	parentType := typeinfo.Types[parentExpr].Type
	if parentType.String() != terr.NameType {
		return false
	}
	childExpr, ok := child.(ast.Expr)
	if !ok {
		return false
	}
	childType := typeinfo.Types[childExpr].Type
	if childType.String() != terr.ValueType {
		return false
	}
	return types.ConvertibleTo(childType, parentType)
}

func rewriteErrFuncArg(path []ast.Node, pkg *loader.PackageInfo, terr *ErrFuncArg) error {
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
				continue
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

func rewriteErrAssign(path []ast.Node, pkg *loader.PackageInfo, terr *ErrAssign) error {
	for i := range path {
		if i+1 >= len(path) {
			break
		}
		child, parent := path[i], path[i+1]
		if assign, ok := parent.(*ast.AssignStmt); ok {
			idx := -1
			for i, r := range assign.Rhs {
				if r == child {
					idx = i
					break
				}
			}
			if idx == -1 {
				continue
			}
			left, right := assign.Lhs[idx], assign.Rhs[idx]
			if !types.ConvertibleTo(pkg.TypeOf(right), pkg.TypeOf(left)) {
				continue
			}
			assign.Rhs[idx] = &ast.CallExpr{
				Fun:  ast.NewIdent(terr.LeftType),
				Args: []ast.Expr{assign.Rhs[idx]},
			}
		}
	}

	return nil
}

func rewriteErrMismatched(path []ast.Node, pkg *loader.PackageInfo, terr *ErrMismatched) error {
	for i := range path {
		if i+1 >= len(path) {
			break
		}
		child, parent := path[i], path[i+1]
		if binaryexpr, ok := parent.(*ast.BinaryExpr); ok {
			if !(child == binaryexpr.X || child == binaryexpr.Y) {
				continue
			}

			ltyp := pkg.Info.TypeOf(binaryexpr.X)
			rtyp := pkg.Info.TypeOf(binaryexpr.Y)

			// TODO(haya14busa): DefaultRule is global variable.
			r2l, r2lOk := DefaultRule.ConvertibleTo(rtyp.String(), ltyp.String())
			r2lOk = r2lOk && types.ConvertibleTo(rtyp, ltyp)
			l2r, l2rOk := DefaultRule.ConvertibleTo(ltyp.String(), rtyp.String())
			l2rOk = l2rOk && types.ConvertibleTo(ltyp, rtyp)

			switch {
			case (r2lOk && !l2rOk) || (r2lOk && l2rOk && r2l > l2r): // right to left
				binaryexpr.Y = &ast.CallExpr{
					Fun:  ast.NewIdent(ltyp.String()),
					Args: []ast.Expr{binaryexpr.Y},
				}
			case (!r2lOk && l2rOk) || (r2lOk && l2rOk && r2l <= l2r): // left to right
				binaryexpr.X = &ast.CallExpr{
					Fun:  ast.NewIdent(rtyp.String()),
					Args: []ast.Expr{binaryexpr.X},
				}
			default:
				return nil
			}
		}
	}
	return nil
}

func rewriteErrReturn(path []ast.Node, pkg *loader.PackageInfo, terr *ErrReturn) error {
	for i := range path {
		if i+3 >= len(path) {
			break
		}
		child, parent, funcDeclNode := path[i], path[i+1], path[i+3]
		returnStmt, ok := parent.(*ast.ReturnStmt)
		if !ok {
			continue
		}
		funcDecl, ok := funcDeclNode.(*ast.FuncDecl)
		if !ok {
			continue
		}
		idx := -1
		for i, r := range returnStmt.Results {
			if r == child {
				idx = i
			}
		}
		if idx == -1 {
			continue
		}
		gotType := pkg.Info.TypeOf(returnStmt.Results[idx])
		wantType := pkg.Info.TypeOf(funcDecl.Type.Results.List[idx].Type)
		if types.ConvertibleTo(gotType, wantType) {
			returnStmt.Results[idx] = &ast.CallExpr{
				Fun:  ast.NewIdent(wantType.String()),
				Args: []ast.Expr{returnStmt.Results[idx]},
			}
			return nil
		}
	}
	return nil
}
