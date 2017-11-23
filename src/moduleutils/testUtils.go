package moduleutils

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/eval"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

var testInterpreter = eval.NewInterpreter()

func TestEval(input string) object.Object {
	l := lexer.NewString(input)
	p := parser.New(l, nil)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return testInterpreter.Eval(program, env)
}

// Verification functions
func TestIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. expected=%d, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
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

func TestFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. expected=%g, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
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

func TestBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Bool. expected=%t, got=%T (%+v)",
			expected,
			obj,
			ShowError(obj),
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

func ShowError(obj object.Object) string {
	if obj == nil {
		return "nil"
	}

	if obj, ok := obj.(*object.Exception); ok {
		return obj.Inspect()
	}
	return obj.Type().String()
}

func TestLiteralErrorObjects(t *testing.T, got object.Object, expected interface{}) {
	switch expected := expected.(type) {
	case int:
		TestIntegerObject(t, got, int64(expected))
	case float64:
		TestFloatObject(t, got, expected)
	case bool:
		TestBoolObject(t, got, expected)
	case string:
		errObj, ok := got.(*object.Exception)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, ShowError(got))
			return
		}

		if errObj.Message != expected {
			t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
		}
	}
}
