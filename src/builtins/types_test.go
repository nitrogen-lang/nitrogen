package eval

import (
	"testing"
)

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

func TestIsBuiltins(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`isFloat(3.14159)`, true},
		{`isFloat(3)`, false},
		{`isInt(3)`, true},
		{`isInt(3.14159)`, false},
		{`isBool(true)`, true},
		{`isBool(false)`, true},
		{`isNull(nil)`, true},
		{`isNull("nil")`, false},
		{`isFunc(func() { 10; })`, true},
		{`isFunc(10)`, false},
		{`isString("Hello")`, true},
		{`isString(10)`, false},
		{`isArray([10, "true", false])`, true},
		{`isArray("array")`, false},
		{`isMap({"key": "value"})`, true},
		{`isMap("array")`, false},
	}

	for _, tt := range tests {
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
	}
}
