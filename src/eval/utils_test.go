package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/object"
	"github.com/nitrogen-lang/nitrogen/src/parser"
)

var testInterpreter = NewInterpreter()

func testEval(input string, t *testing.T) object.Object {
	l := lexer.NewString(input)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) > 1 {
		t.Fatal(p.Errors()[0])
	}
	env := object.NewEnvironment()

	return testInterpreter.Eval(program, env)
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
	if obj.Type() != object.NullObj {
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

	if obj, ok := obj.(*object.Exception); ok {
		return obj.Inspect()
	}
	return obj.Type().String()
}

func TestStringStack(t *testing.T) {
	stack := newStringStack()

	if stack.len() != 0 {
		t.Fatalf("Stack has none 0 length: %d", stack.len())
	}

	if stack.getFront() != "" {
		t.Fatalf("Stack front should be empty string: %s", stack.getFront())
	}

	stack.push("item1")
	if stack.len() != 1 {
		t.Fatalf("Stack length should be 1: %d", stack.len())
	}

	if stack.getFront() != "item1" {
		t.Fatalf("Stack front should be `item1`: %s", stack.getFront())
	}

	stack.push("item2")
	if stack.len() != 2 {
		t.Fatalf("Stack length should be 2: %d", stack.len())
	}

	if stack.getFront() != "item2" {
		t.Fatalf("Stack front should be `item2`: %s", stack.getFront())
	}

	stack.pop()
	if stack.len() != 1 {
		t.Fatalf("Stack length should be 1: %d", stack.len())
	}

	if stack.getFront() != "item1" {
		t.Fatalf("Stack front should be `item1`: %s", stack.getFront())
	}

	stack.pop()
	if stack.len() != 0 {
		t.Fatalf("Stack has none 0 length: %d", stack.len())
	}
}
