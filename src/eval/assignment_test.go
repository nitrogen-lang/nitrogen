package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestDefStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input, t), tt.expected)
	}
}

func TestAlwaysStatements(t *testing.T) {
	evaled := testEval("always a = 5; a;", t)
	testIntegerObject(t, evaled, 5)

	tests := []struct {
		input    string
		expected string
	}{
		{
			"always a = 5; a = 6;",
			"Assignment to declared constant a",
		},
		{
			"always a = 5; let a = 6;",
			"Assignment to declared constant a",
		},
		{
			"always a = [5, 6, 7];",
			"Constants must be int, float, string, bool or null",
		},
	}

	for _, test := range tests {
		evaled = testEval(test.input, t)
		errObj, ok := evaled.(*object.Exception)

		if !ok {
			t.Fatalf("expected error redeclaring const. got=%T(%+v)", evaled, evaled)
		}

		if errObj.Message != test.expected {
			t.Fatalf("wronte error message. expected %q, got %q", errObj.Message, test.expected)
		}
	}
}

func TestHashAssignEval(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`let m = {"foo": 5}; m["foo"] = 6; m["foo"]`,
			6,
		},
		{
			`let m = [1, 3]; m[1] = 5; m[1]`,
			5,
		},
		{
			`n = 34; n`,
			"Assignment to uninitialized variable n",
		},
		{
			`let n = 0; let p = (n = 34); p`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)

		switch i := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(i))
		case string:
			err, ok := evaluated.(*object.Exception)
			if !ok {
				t.Fatalf("Expected Error, got %T", evaluated)
			}
			if err.Message != i {
				t.Fatalf("Incorrect error, expected %s, got %s", i, err.Message)
			}
		default:
			testNullObject(t, evaluated)
		}
	}
}
