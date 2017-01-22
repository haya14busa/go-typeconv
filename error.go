package typeconv

import (
	"go/types"
	"regexp"
)

// TypeError represents type error.
type TypeError interface {
	typ() typErr
}

// ErrVarDecl represents type error of variable declaration.
//
// Example:
//	var x, y int = 3, 4
//	var _ float64 = x*x + y*y
//	var _ uint = x
//	var _ uint = int(1)
//	var _ uint = int(1) + int(1)
//
// Error:
//	cannot use x * x + y * y (value of type int) as float64 value in variable declaration
//	cannot use x (variable of type int) as uint value in variable declaration
//	cannot use int(1) (constant 1 of type int) as uint value in variable declaration
//	cannot use int(1) + int(1) (constant 2 of type int) as uint value in variable declaration
type ErrVarDecl struct {
	NameType  string
	ValueType string
}

func (*ErrVarDecl) typ() typErr {
	return TypeErrVarDecl
}

type typErr int

const (
	// cannot use x (variable of type int) as uint value in variable declaration
	TypeErrVarDecl typErr = iota

	// cannot use x (variable of type int) as float64 value in argument to funcarg
	TypeErrFuncArg

	// cannot use y (type int) as type float64 in assignment
	TypeErrAssign

	// invalid operation: mismatched types int and float64
	TypeErrMismatched
)

var regexps = [...]*regexp.Regexp{
	TypeErrVarDecl: regexp.MustCompile(`\((constant .+|variable|value) of type (?P<got>.+)\) as (?P<want>.+) value in variable declaration$`),
}

// NewTypeErr creates TypeError from types.Error.
func NewTypeErr(err types.Error) TypeError {
	for i, re := range regexps {
		ms := re.FindStringSubmatch(err.Msg)
		if len(ms) == 0 {
			continue
		}
		names := re.SubexpNames()
		switch typErr(i) {
		case TypeErrVarDecl:
			return newErrVarDecl(ms, names)
		}
	}
	return nil
}

func newErrVarDecl(matches, names []string) *ErrVarDecl {
	err := &ErrVarDecl{}
	for i, name := range names {
		if i == 0 {
			continue
		}
		m := matches[i]
		switch name {
		case "got":
			err.ValueType = m
		case "want":
			err.NameType = m
		}
	}
	return err
}
