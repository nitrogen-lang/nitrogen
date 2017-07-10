package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

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

func TestBuiltinLenFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "Unsupported type INTEGER"},
		{`len("one", "two")`, "Incorrect number of arguments. Got 2, expected 1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinFirstFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first("four")`, "Argument to `first` must be ARRAY, got STRING"},
		{`first()`, "Incorrect number of arguments. Got 0, expected 1"},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinLastFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last("four")`, "Argument to `last` must be ARRAY, got STRING"},
		{`last()`, "Incorrect number of arguments. Got 0, expected 1"},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinRestFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`rest([1])`, `[]`},
		{`rest([1, 2, 3])`, `[2, 3]`},
		{`rest([])`, ""},
		{`rest("four")`, "Argument to `rest` must be ARRAY, got STRING"},
		{`rest()`, "Incorrect number of arguments. Got 0, expected 1"},
	}

	for _, tt := range tests {
		got := testEval(tt.input)

		if _, ok := got.(*object.Null); ok {
			if tt.expected != "" {
				t.Errorf("Incorrect value. Expected=%s, got=NULL", tt.expected)
			}
			continue
		}

		if arrObj, ok := got.(*object.Array); ok {
			if arrObj.Inspect() != tt.expected {
				t.Errorf("Incorrect array. Expected=%s, got=%s",
					tt.expected, arrObj.Inspect())
			}
			continue
		}

		errObj, ok := got.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, showError(got))
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestBuiltinPushFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`push([1], 2)`, `[1, 2]`},
		{`push([1, 2, 3], 4)`, `[1, 2, 3, 4]`},
		{`push([], 1)`, `[1]`},
		{`push("four", "five")`, "Argument to `push` must be ARRAY, got STRING"},
		{`push()`, "Incorrect number of arguments. Got 0, expected 2"},
		{`push([1])`, "Incorrect number of arguments. Got 1, expected 2"},
	}

	for _, tt := range tests {
		got := testEval(tt.input)

		if arrObj, ok := got.(*object.Array); ok {
			if arrObj.Inspect() != tt.expected {
				t.Errorf("Incorrect array. Expected=%s, got=%s",
					tt.expected, arrObj.Inspect())
			}
			continue
		}

		errObj, ok := got.(*object.Error)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, showError(got))
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestBuiltinIntConvFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`to_int(23.5)`, 23},
		{`to_int(1)`, 1},
		{`to_int("hello world")`, "Argument to `to_int` must be FLOAT or INT, got STRING"},
		{`to_int([])`, "Argument to `to_int` must be FLOAT or INT, got ARRAY"},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}

func TestBuiltinFloatConvFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`to_float(23.5)`, 23.5},
		{`to_float(1)`, 1.0},
		{`to_float("hello world")`, "Argument to `to_float` must be FLOAT or INT, got STRING"},
		{`to_float([])`, "Argument to `to_float` must be FLOAT or INT, got ARRAY"},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}
