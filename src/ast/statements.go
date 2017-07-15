package ast

import (
	"bytes"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type DefStatement struct {
	Token token.Token // the token.DEF token
	Const bool
	Name  *Identifier
	Value Expression
}

func (d *DefStatement) statementNode()       {}
func (d *DefStatement) TokenLiteral() string { return d.Token.Literal }
func (d *DefStatement) String() string {
	if _, ok := d.Value.(*FunctionLiteral); ok {
		return d.Value.String()
	}

	var out bytes.Buffer

	if d.Const {
		out.WriteString("always ")
	} else {
		out.WriteString("let ")
	}
	out.WriteString(d.Name.String())
	out.WriteString(" = ")
	if d.Value != nil {
		out.WriteString(d.Value.String())
	}
	out.WriteByte(';')

	return out.String()
}

type AssignStatement struct {
	Token token.Token // the token.DEF token
	Left  Expression
	Value Expression
}

func (a *AssignStatement) expressionNode()      {}
func (a *AssignStatement) TokenLiteral() string { return a.Token.Literal }
func (a *AssignStatement) String() string {
	var out bytes.Buffer

	out.WriteString(a.Left.String())
	out.WriteString(" = ")
	if a.Value != nil {
		out.WriteString(a.Value.String())
	}
	out.WriteByte(';')

	return out.String()
}

type ReturnStatement struct {
	Token token.Token // the 'return' token
	Value Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString("return ")
	if r.Value != nil {
		out.WriteString(r.Value.String())
	}
	out.WriteByte(';')

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	l := len(bs.Statements) - 1
	for i, s := range bs.Statements {
		str := s.String()
		out.WriteString(str)
		if str[len(str)-1] != ';' {
			out.WriteByte(';')
		}

		if i < l {
			out.WriteByte(' ')
		}
	}
	return out.String()
}
