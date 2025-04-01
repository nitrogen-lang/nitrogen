package collections

import (
	"testing"

	mut "github.com/nitrogen-lang/nitrogen/src/moduleutils_test"
	"github.com/nitrogen-lang/nitrogen/src/object"
)

func TestBuiltinLenFunction(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "len(): Unsupported type INTEGER"},
		{`len("one", "two")`, "Incorrect number of arguments. Got 2, expected 1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`len(nil)`, 0},
	}

	for _, tt := range tests {
		mut.TestLiteralErrorObjects(t, mut.TestEval(tt.input), tt.expected)
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
		mut.TestLiteralErrorObjects(t, mut.TestEval(tt.input), tt.expected)
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
		mut.TestLiteralErrorObjects(t, mut.TestEval(tt.input), tt.expected)
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
		got := mut.TestEval(tt.input)

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

		errObj, ok := got.(*object.Exception)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, mut.ShowError(got))
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
		got := mut.TestEval(tt.input)

		if arrObj, ok := got.(*object.Array); ok {
			if arrObj.Inspect() != tt.expected {
				t.Errorf("Incorrect array. Expected=%s, got=%s",
					tt.expected, arrObj.Inspect())
			}
			continue
		}

		errObj, ok := got.(*object.Exception)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, mut.ShowError(got))
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestBuiltinHashMerge(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`hashMerge({"key2": "value"}, {"key2": "value2"})`, `{key2: "value2"}`},       // Test overwrite
		{`hashMerge({"key2": "value"}, {"key2": "value2"}, false)`, `{key2: "value"}`}, // Test no overwrite
		{`hashMerge()`, "hashMerge requires at least 2 arguments. Got 0"},
		{`hashMerge({"key": "value"}, 10)`, "First two arguments must be maps"},
		{`hashMerge(10, {"key": "value"})`, "First two arguments must be maps"},
	}

	for _, tt := range tests {
		got := mut.TestEval(tt.input)

		if hashObj, ok := got.(*object.Hash); ok {
			inspect := hashObj.Inspect()
			if inspect != tt.expected {
				t.Errorf("Incorrect hash map. Expected='%s', got='%s'",
					tt.expected, inspect)
			}
			continue
		}

		errObj, ok := got.(*object.Exception)
		if !ok {
			t.Errorf("object is not Error. got=%T (%+v)", got, mut.ShowError(got))
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestBuiltinHashMergeSpecial(t *testing.T) {
	evaled := mut.TestEval(`hashMerge({"key": "value"}, {"key2": "value2"})`)
	hashObj, ok := evaled.(*object.Hash)
	if !ok {
		t.Fatalf("Got error during hashMerge: %#v", evaled)
	}

	expected := map[object.HashKey]string{
		(object.MakeStringObj("key")).HashKey():  "value",
		(object.MakeStringObj("key2")).HashKey(): "value2",
	}

	for k, v := range expected {
		val, exists := hashObj.Pairs[k]
		if !exists {
			t.Fatalf("Map missing key %q", k)
		}

		valStr := val.Value.(*object.String).String()
		if valStr != v {
			t.Fatalf("Incorrect map value for key %q: %q", k, v)
		}
	}
}
