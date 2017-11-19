package typing

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/moduleutils"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestBuiltinIntConvFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`toInt(23.5)`, 23},
		{`toInt(1)`, 1},
		{`toInt("hello world")`, "Argument to `toInt` must be FLOAT or INT, got STRING"},
		{`toInt([])`, "Argument to `toInt` must be FLOAT or INT, got ARRAY"},
	}

	for _, tt := range tests {
		moduleutils.TestLiteralErrorObjects(t, moduleutils.TestEval(tt.input), tt.expected)
	}
}

func TestBuiltinFloatConvFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`toFloat(23.5)`, 23.5},
		{`toFloat(1)`, 1.0},
		{`toFloat("hello world")`, "Argument to `toFloat` must be FLOAT or INT, got STRING"},
		{`toFloat([])`, "Argument to `toFloat` must be FLOAT or INT, got ARRAY"},
	}

	for _, tt := range tests {
		moduleutils.TestLiteralErrorObjects(t, moduleutils.TestEval(tt.input), tt.expected)
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
		moduleutils.TestLiteralErrorObjects(t, moduleutils.TestEval(tt.input), tt.expected)
	}
}

func TestToStringBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`toString(5)`, "5"},
		{`toString(5.6)`, "5.6"},
		{`toString("Hello")`, "Hello"},
		{`toString(true)`, "true"},
		{`toString(false)`, "false"},
		{`toString(nil)`, "nil"},
		{`toString([1, 2, 3])`, ""},
		{`toString({"key": "value"})`, ""},
	}

	for i, tt := range tests {
		evaled := moduleutils.TestEval(tt.input)

		str, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("%d: Expected string, got %#v", i+1, evaled)
		}

		if str.Value != tt.expected {
			t.Fatalf("%d: toString failed. Expected `%s`, got `%s`", i+1, tt.expected, str.Value)
		}
	}
}

func TestParseIntBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`parseInt("5")`, 5},
		{`parseInt("-5")`, -5},
		{`parseInt("5.6")`, 0},
		{`parseInt("Hello")`, 0},
		{`parseInt(true)`, -1},
	}

	for i, tt := range tests {
		evaled := moduleutils.TestEval(tt.input)

		if tt.expected == 0 {
			if evaled != object.NullConst {
				t.Fatalf("Expected nil, got %#v", evaled)
			}
			continue
		}

		if tt.expected == -1 {
			if evaled.Type() != object.ExceptionObj {
				t.Fatalf("Expected error, got %#v", evaled)
			}
			continue
		}

		io, ok := evaled.(*object.Integer)
		if !ok {
			t.Fatalf("%d: Expected int, got %#v", i+1, evaled)
		}

		if io.Value != tt.expected {
			t.Fatalf("%d: toInt failed. Expected `%d`, got `%d`", i+1, tt.expected, io.Value)
		}
	}
}

func TestParseFloatBuiltin(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{`parseFloat("5")`, 5.0},
		{`parseFloat("-5")`, -5.0},
		{`parseFloat("5.6")`, 5.6},
		{`parseFloat("Hello")`, 0},
		{`parseFloat(true)`, -1},
	}

	for i, tt := range tests {
		evaled := moduleutils.TestEval(tt.input)

		if tt.expected == 0 {
			if evaled != object.NullConst {
				t.Fatalf("Expected nil, got %#v", evaled)
			}
			continue
		}

		if tt.expected == -1 {
			if evaled.Type() != object.ExceptionObj {
				t.Fatalf("Expected error, got %#v", evaled)
			}
			continue
		}

		io, ok := evaled.(*object.Float)
		if !ok {
			t.Fatalf("%d: Expected float, got %#v", i+1, evaled)
		}

		if io.Value != tt.expected {
			t.Fatalf("%d: toFloat failed. Expected `%f`, got `%f`", i+1, tt.expected, io.Value)
		}
	}
}
