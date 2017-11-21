package parser

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseClassLiteral() ast.Expression {
	c := &ast.ClassLiteral{
		Fields:  make([]*ast.DefStatement, 0),
		Methods: make(map[string]*ast.FunctionLiteral),
	}

	if p.peekTokenIs(token.Carrot) {
		p.nextToken()
		if !p.expectPeek(token.Identifier) {
			return nil
		}
		c.Parent = p.curToken.Literal
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	body := p.parseBlockStatements()

	for _, statement := range body.Statements {
		def, ok := statement.(*ast.DefStatement)
		if !ok {
			p.addErrorWithPos("Only function and variable statements are allowed in a class definition")
			return nil
		}

		switch s := def.Value.(type) {
		case *ast.FunctionLiteral:
			c.Methods[s.Name] = s
		default:
			c.Fields = append(c.Fields, def)
		}
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return c
}

func (p *Parser) parseMakeExpression() ast.Expression {
	m := &ast.MakeInstance{}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	m.Class = p.curToken.Literal
	p.nextToken()
	m.Arguments = p.parseExpressionList(token.RParen)

	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return m
}
