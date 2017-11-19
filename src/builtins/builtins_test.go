package builtins

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestBuiltinsCantBeOverridden(t *testing.T) {
	input := `let len = func(x) { x }`
	evaled := moduleutils.TestEval(input)
	if evaled == nil {
		t.Fatal("no error object returned")
	}

	err, ok := evaled.(*object.Exception)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaled, moduleutils.ShowError(evaled))
	}

	if err.Message != "Attempted redeclaration of builtin function 'len'" {
		t.Errorf("Error has wrong message. got=%q", err.Message)
	}
}
