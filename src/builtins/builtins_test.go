package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

func testEval(input string) object.Object {
	l := lexer.NewString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return eval.Eval(program, env)
}

// Verification functions
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. expected=%d, got=%T (%+v)",
			expected,
			obj,
			showError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. expected=%g, got=%T (%+v)",
			expected,
			obj,
			showError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%g, want=%g",
			result.Value, expected)
		return false
	}
	return true
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Bool. expected=%t, got=%T (%+v)",
			expected,
			obj,
			showError(obj),
		)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func showError(obj object.Object) string {
	if obj == nil {
		return "nil"
	}

	if obj, ok := obj.(*object.Error); ok {
		return obj.Inspect()
	}
	return obj.Type().String()
}

func TestBuiltinsCantBeOverridden(t *testing.T) {
	input := `let len = func(x) { x }`
	evaled := testEval(input)
	if evaled == nil {
		t.Fatal("no error object returned")
	}

	err, ok := evaled.(*object.Error)
	if !ok {
		t.Fatalf("object is not Error. got=%T (%+v)", evaled, showError(evaled))
	}

	if err.Message != "Attempted redeclaration of builtin function 'len'" {
		t.Errorf("Error has wrong message. got=%q", err.Message)
	}
}

func testLiteralErrorObjects(t *testing.T, got object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		testIntegerObject(t, got, int64(expected))
	case float64:
		testFloatObject(t, got, expected)
	case bool:
		testBoolObject(t, got, expected)
	case string:
		errObj, ok := got.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, showError(got))
			return
		}

		if errObj.Message != expected {
			t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
		}
	}
}
