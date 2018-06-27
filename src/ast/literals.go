package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type NullLiteral struct {
	Token token.Token // the token.NULL token
}

func (n *NullLiteral) expressionNode()      {}
func (n *NullLiteral) TokenLiteral() string { return n.Token.Literal }
func (n *NullLiteral) String() string       { return n.Token.Literal }

type IntegerLiteral struct {
	Token token.Token // the token.INT token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

type FloatLiteral struct {
	Token token.Token // the token.FLOAT token
	Value float64
}

func (f *FloatLiteral) expressionNode()      {}
func (f *FloatLiteral) TokenLiteral() string { return f.Token.Literal }
func (f *FloatLiteral) String() string       { return f.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value []rune
}

func (s *StringLiteral) expressionNode()      {}
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }
func (s *StringLiteral) String() string       { return string(s.Value) }

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token // The 'func' token
	Name       string
	FQName     string
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteByte(' ')
	out.WriteString(fl.Name)
	out.WriteByte('(')
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {")
	out.WriteString(fl.Body.String())
	out.WriteByte('}')
	return out.String()
}

type Array struct {
	Token    token.Token // the '[' token
	Elements []Expression
}

func (a *Array) expressionNode()      {}
func (a *Array) TokenLiteral() string { return a.Token.Literal }
func (a *Array) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.String())
	}

	out.WriteByte('[')
	out.WriteString(strings.Join(elements, ", "))
	out.WriteByte(']')

	return out.String()
}

type HashLiteral struct {
	Token token.Token // the '{' token
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expressionNode()      {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}
	for key, value := range h.Pairs {
		pairs = append(pairs, key.String()+": "+value.String())
	}

	out.WriteByte('{')
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteByte('}')
	return out.String()
}

type ClassLiteral struct {
	Name    string
	Parent  string
	Fields  []*DefStatement
	Methods map[string]*FunctionLiteral
}

func (c *ClassLiteral) expressionNode()      {}
func (c *ClassLiteral) TokenLiteral() string { return "class" }
func (c *ClassLiteral) String() string {
	return fmt.Sprintf("class %s ^ %s {...}", c.Name, c.Parent)
}
