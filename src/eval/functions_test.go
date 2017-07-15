package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated := testEval(input, t)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%#v)", evaluated, showError(evaluated))
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2);"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { x * 2; }; double(5);", 10},
		{"let add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input, t), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
        let newAdder = func(x) {
            func(y) { x + y };
        };

        let addTwo = newAdder(2);
        addTwo(2);`

	testIntegerObject(t, testEval(input, t), 4)
}

func TestExtraArgs(t *testing.T) {
	input := `func extra(a) { args[0]; } extra(1, 2)`
	testIntegerObject(t, testEval(input, t), 2)
}

func TestExtraArgsError(t *testing.T) {
	input := `func extra(a) { args = 5; } extra(1, 2)`
	evaled := testEval(input, t)
	errObj, ok := evaled.(*object.Error)
	if !ok {
		t.Fatalf("Expected error, got %#v", evaled)
	}

	if errObj.Message != "Assignment to declared constant args" {
		t.Fatalf("Incorrect error message. Got '%s'", errObj.Message)
	}
}
