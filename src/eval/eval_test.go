package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestNullEval(t *testing.T) {
	input := `nil`
	evaluated := testEval(input, t)

	if evaluated != object.NULL {
		t.Fatalf("object is not Null. got=%T (%+v)", evaluated, showError(evaluated))
	}
}

// Actual test cases
func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5.5", 5.5},
		{"10.2", 10.2},
		{"-5.6", -5.6},
		{"-10.8", -10.8},
		{"5.5 + 5.5 + 5.2 + 5.3 - 10.1", 11.4},
		{"2.5 * 2.2 * 2.3 * 2.1 * 2.4", 63.75599999999999},
		{"-50.2 + 100.0 + -50.1", -0.30000000000000426},
		{"5.2 * 2.3 + 10.1", 22.06},
		{"50.2 / 2.1 * 2.3 + 10.5", 65.48095238095237},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"6 % 3", 0},
		{"6 % 3 + 4", 4},
		{"6 % 4 * 4", 8},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{`
            if (10 > 1) {
                if (10 > 1) {
                    return 10;
                }
                return 1;
            }`, 10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`if (10 > 1) {
                if (10 > 1) {
                    return true + false;
                }
                return 1;
            }
            `,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[func(x) { x }];`,
			"Invalid map key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input, t)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, showError(evaluated))
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input, t)
	str, ok := evaluated.(*object.String)

	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, showError(evaluated))
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringEquality(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{`"foobar" == "foobar"`, true},
		{`"foobar" == "foo bar"`, false},
		{`"foobar" != "foo bar"`, true},
		{`"foobar" != "foobar"`, false},
	}

	for _, tt := range tests {
		testBooleanObject(t, testEval(tt.input, t), tt.expected)
	}
}
