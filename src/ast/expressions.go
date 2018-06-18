package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type PrefixExpression struct {
	Token    token.Token // the prefix token, e.g. !
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteByte(')')

	return out.String()
}

type InfixExpression struct {
	Token    token.Token // the infix token, e.g. +, <
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(i.Left.String())
	out.WriteByte(' ')
	out.WriteString(i.Operator)
	out.WriteByte(' ')
	out.WriteString(i.Right.String())
	out.WriteByte(')')

	return out.String()
}

type CompareExpression struct {
	Token token.Token // and, or
	Left  Expression
	Right Expression
}

func (i *CompareExpression) expressionNode()      {}
func (i *CompareExpression) TokenLiteral() string { return i.Token.Literal }
func (i *CompareExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(i.Left.String())
	out.WriteByte(' ')
	out.WriteString(i.Token.Literal)
	out.WriteByte(' ')
	out.WriteString(i.Right.String())
	out.WriteByte(')')

	return out.String()
}

type IfExpression struct {
	Token       token.Token // The 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if ")
	out.WriteString(ie.Condition.String())
	out.WriteString(" {")
	out.WriteString(ie.Consequence.String())
	out.WriteByte('}')
	if ie.Alternative != nil {
		out.WriteString(" else {")
		out.WriteString(ie.Alternative.String())
		out.WriteByte('}')
	}
	return out.String()
}

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // The '[' token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(i.Left.String())
	out.WriteByte('[')
	out.WriteString(i.Index.String())
	out.WriteString("])")

	return out.String()
}

type AttributeExpression struct {
	Token token.Token
	Left  Expression
	Index *StringLiteral
}

func (i *AttributeExpression) expressionNode()      {}
func (i *AttributeExpression) TokenLiteral() string { return i.Token.Literal }
func (i *AttributeExpression) String() string {
	var out bytes.Buffer

	out.WriteByte('(')
	out.WriteString(i.Left.String())
	out.WriteByte('.')
	out.WriteString(i.Index.String())
	out.WriteByte(')')

	return out.String()
}

type TryCatchExpression struct {
	Try    *BlockStatement
	Catch  *BlockStatement
	Symbol *Identifier
}

func (t *TryCatchExpression) expressionNode()      {}
func (t *TryCatchExpression) TokenLiteral() string { return "try" }
func (t *TryCatchExpression) String() string {
	var out bytes.Buffer
	out.WriteString("try {")
	out.WriteString(t.Try.String())
	out.WriteString("} catch ")
	if t.Symbol != nil {
		out.WriteString(t.Symbol.String())
		out.WriteByte(' ')
	}
	out.WriteByte('{')
	out.WriteString(t.Catch.String())
	out.WriteByte('}')
	return out.String()
}

type MakeInstance struct {
	Class     Expression
	Arguments []Expression
}

func (m *MakeInstance) expressionNode()      {}
func (m *MakeInstance) TokenLiteral() string { return "make" }
func (m *MakeInstance) String() string {
	args := make([]string, len(m.Arguments))
	for i, a := range m.Arguments {
		args[i] = a.String()
	}
	return fmt.Sprintf("make %s(%s)", m.Class, strings.Join(args, ", "))
}
