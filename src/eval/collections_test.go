package eval

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input, t)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, showError(evaluated))
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			3,
		},
		{
			"[1, 2, 3][-3]",
			1,
		},
		{
			"[1, 2, 3][-4]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"[1, 2, 3] + [4, 5]", "[1, 2, 3, 4, 5]"},
		{"[] + [1, 2]", "[1, 2]"},
		{"[1, 2] + []", "[1, 2]"},
		{"[] + []", "[]"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		arr, ok := evaluated.(*object.Array)

		if !ok {
			t.Fatalf("object is not String. got=%T (%+v)", evaluated, showError(evaluated))
		}

		if arr.Inspect() != tt.expected {
			t.Errorf("Array has wrong value. expected=%q, got=%q",
				tt.expected, arr.Inspect())
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
		{
			"one": 10 - 9,
			two: 1 + 1,
			"thr" + "ee": 6 / 2,
			4: 4,
		}`

	evaluated := testEval(input, t)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, showError(evaluated))
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}
		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestStringIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`"Hello, world"[0]`,
			"H",
		},
		{
			`"Hello, world"[1]`,
			"e",
		},
		{
			`"Hello, world"[2]`,
			"l",
		},
		{
			`let i = 0; "Hello, world"[i];`,
			"H",
		},
		{
			`"Hello, world"[1 + 1];`,
			"l",
		},
		{
			`let myArray = "Hello, world"; myArray[2];`,
			"l",
		},
		{
			`let myArray = "Hello, world"; myArray[0] + myArray[1] + myArray[2];`,
			"Hel",
		},
		{
			`"Hello, world"[12]`,
			nil,
		},
		{
			`"Hello, world"[-1]`,
			"d",
		},
		{
			`"Hello, world"[-3]`,
			"r",
		},
		{
			`"Hello, world"[-13]`,
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input, t)
		str, ok := tt.expected.(string)
		if ok {
			testStringObject(t, evaluated, str)
		} else {
			testNullObject(t, evaluated)
		}
	}
}
