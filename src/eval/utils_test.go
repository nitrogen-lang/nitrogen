package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

func testEval(input string) object.Object {
	l := lexer.NewString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. expected=%t, got=%T (%+v)",
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

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj.Type() != object.NULL_OBJ {
		t.Errorf("object is not Null. got=%T (%+v)", obj, showError(obj))
		return false
	}
	return true
}

func testStringObject(t *testing.T, got object.Object, expected string) {
	str, ok := got.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", got, showError(got))
	}

	if str.Value != expected {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
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
