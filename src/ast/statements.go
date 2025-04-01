package ast

import (
	"bytes"
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Filename   string
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
		out.WriteString("const ")
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

type ImportStatement struct {
	Token token.Token // the token.Import token
	Path  *StringLiteral
	Name  *Identifier
}

func (i *ImportStatement) statementNode()       {}
func (i *ImportStatement) TokenLiteral() string { return i.Token.Literal }
func (i *ImportStatement) String() string {
	return fmt.Sprintf("import %s as %s;", i.Path.String(), i.Name.String())
}

type DeleteStatement struct {
	Token token.Token
	Name  string
}

func (d *DeleteStatement) statementNode()       {}
func (d *DeleteStatement) TokenLiteral() string { return d.Token.Literal }
func (d *DeleteStatement) String() string {
	return fmt.Sprintf("delete %s;", d.Name)
}

type AssignStatement struct {
	Token token.Token // the token.DEF token
	Left  Expression
	Value Expression
}

func (a *AssignStatement) statementNode()       {}
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
		if len(str) == 0 || str[len(str)-1] != ';' {
			out.WriteByte(';')
		}

		if i < l {
			out.WriteByte(' ')
		}
	}
	return out.String()
}

type LoopStatement struct {
	Token     token.Token
	Init      *DefStatement
	Condition Expression
	Iter      Node
	Body      *BlockStatement
}

func (fl *LoopStatement) statementNode() {}
func (fl *LoopStatement) TokenLiteral() string {
	if fl.Init == nil {
		return "while"
	}
	return "for"
}
func (fl *LoopStatement) String() string {
	var out bytes.Buffer
	out.WriteString(fl.TokenLiteral())
	out.WriteByte(' ')
	if fl.Init != nil {
		out.WriteString(fl.Init.String())
		out.WriteString("; ")
	}
	out.WriteString(fl.Condition.String())
	if fl.Iter != nil {
		out.WriteString("; ")
		out.WriteString(fl.Iter.String())
	}
	out.WriteString(" { ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}

type IterLoopStatement struct {
	Token token.Token
	Key   *Identifier
	Value *Identifier
	Iter  Expression
	Body  *BlockStatement
}

func (fl *IterLoopStatement) statementNode()       {}
func (fl *IterLoopStatement) TokenLiteral() string { return "for" }
func (fl *IterLoopStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for ")

	if fl.Key != nil {
		out.WriteString(fl.Key.String())
		out.WriteString(", ")
	}

	out.WriteString(fl.Value.String())
	out.WriteString(" in ")
	out.WriteString(fl.Iter.String())

	out.WriteString(" { ")
	out.WriteString(fl.Body.String())
	out.WriteString(" }")
	return out.String()
}

type ContinueStatement struct {
	Token token.Token
}

func (c *ContinueStatement) statementNode()       {}
func (c *ContinueStatement) TokenLiteral() string { return "continue" }
func (c *ContinueStatement) String() string       { return "continue" }

type BreakStatement struct {
	Token token.Token
}

func (b *BreakStatement) statementNode()       {}
func (b *BreakStatement) TokenLiteral() string { return "break" }
func (b *BreakStatement) String() string       { return "break" }

type PassStatement struct {
	Token token.Token
}

func (b *PassStatement) statementNode()       {}
func (b *PassStatement) TokenLiteral() string { return "pass" }
func (b *PassStatement) String() string       { return "pass" }

type BreakpointStatement struct {
	Token token.Token
}

func (b *BreakpointStatement) statementNode()       {}
func (b *BreakpointStatement) TokenLiteral() string { return "breakpoint" }
func (b *BreakpointStatement) String() string       { return "breakpoint" }
