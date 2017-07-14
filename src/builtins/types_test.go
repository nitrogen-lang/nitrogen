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
