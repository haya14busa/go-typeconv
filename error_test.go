package typeconv

import (
	"go/types"
	"testing"
)

func TestNewTypeErr(t *testing.T) {
	tests := []struct {
		in      string
		wantTyp typErr
	}{
		{
			in:      "cannot use x (variable of type int) as uint value in variable declaration",
			wantTyp: TypeErrVarDecl,
		},
		{
			in:      "cannot use x * x + y * y (value of type int) as float64 value in variable declaration",
			wantTyp: TypeErrVarDecl,
		},
		{
			in:      "cannot use x (variable of type int) as float64 value in argument to funcarg",
			wantTyp: TypeErrFuncArg,
		},
	}

	for _, tt := range tests {
		terr := NewTypeErr(types.Error{Msg: tt.in})
		if terr == nil {
			t.Errorf("got nil for %v", tt.in)
			continue
		}
		if got := terr.typ(); got != tt.wantTyp {
			t.Errorf("type: got %v, want %v", got, tt.wantTyp)
		}
	}

}
