package ast

import (
	"testing"

	"github.com/lfkeitel/nitrogen/src/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DefStatement{
				Token: token.Token{Type: token.DEF, Literal: "def"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "def myVar = anotherVar;" {
		t.Errorf("program.String() wrong. Got %q", program.String())
	}
}
