package eval

import (
	"testing"

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
		testLiteralErrorObjects(t, testEval(tt.input), tt.expected)
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
		evaled := testEval(tt.input)

		str, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("%d: Expected string, got %#v", i+1, evaled)
		}

		if str.Value != tt.expected {
			t.Fatalf("%d: toString failed. Expected `%s`, got `%s`", i+1, tt.expected, str.Value)
		}
	}
}
