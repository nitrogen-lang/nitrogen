package parser

import "testing"

var testCases = []struct {
	i string
	e bool
}{
	{i: `hello`, e: true},
	{i: `1hello`, e: false},
	{i: `hello2`, e: true},
	{i: `_hello2`, e: true},
	{i: `_hello`, e: true},
	{i: `hel_lo`, e: true},
	{i: `hello_`, e: true},
	{i: `@hello_`, e: false},
	{i: `世界`, e: true},
}

func TestIdentMatch(t *testing.T) {
	for _, tcase := range testCases {
		res := isIdent(tcase.i)
		if res != tcase.e {
			t.Fatalf("isIdent failed for %s. Expected %t, got %t", tcase.i, tcase.e, res)
		}
	}
}
