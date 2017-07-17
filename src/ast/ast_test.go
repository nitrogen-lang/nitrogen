package ast

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&DefStatement{
				Token: token.Token{Type: token.Let, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.Identifier, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. Got %q", program.String())
	}
}
