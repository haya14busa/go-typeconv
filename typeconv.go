// Package typeconv provides missing implicit type conversion in Go by
// rewriting AST.
package typeconv

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/loader"
)

// Load creates the initial packages specified by conf along with slice of
// types.Error.
func Load(conf loader.Config, args []string) (*loader.Program, []types.Error, error) {
	conf.AllowErrors = true
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
	conf.FromArgs(args, true)
	prog, err := conf.Load()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load program: %v", err)
	}
	return prog, typeErrs, nil
}
