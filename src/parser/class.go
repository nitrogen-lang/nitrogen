package parser

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseClassLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseClassLiteral")
	}

	classToken := p.curToken
	c := &ast.ClassLiteral{
		Token:   classToken,
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
			p.addErrorWithPos(classToken.Pos, "Only function and variable statements are allowed in a class definition")
			return nil
		}

		switch s := def.Value.(type) {
		case *ast.FunctionLiteral:
			c.Methods[s.Name] = s
		default:
			c.Fields = append(c.Fields, def)
		}
	}

	errored := false
	for _, f := range c.Fields {
		if fn, exists := c.Methods[f.Name.Value]; exists {
			p.addErrorWithPos(fn.Token.Pos, "Duplicate named function %s", fn.Name)
			errored = true
		}
	}

	if errored {
		return nil
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return c
}

func (p *Parser) parseMakeExpression() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseMakeExpression")
	}
	m := &ast.NewInstance{
		Token: p.curToken,
	}

	p.nextToken()
	cExpression := p.parseExpression(priLowest)

	call, ok := cExpression.(*ast.CallExpression)
	if !ok {
		p.addErrorWithCurPos("Invalid object creation")
		return nil
	}

	m.Class = call.Function
	m.Arguments = call.Arguments

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return m
}

func (p *Parser) parseInterfaceLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseInterfaceLiteral")
	}

	interfaceToken := p.curToken
	iface := &ast.InterfaceLiteral{
		Token:   interfaceToken,
		Methods: make(map[string]*ast.IfaceMethodDef),
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}
	p.nextToken()

	for !p.curTokenIs(token.RBrace) {
		if !p.curTokenIs(token.Identifier) {
			return nil
		}

		ifaceMeth := &ast.IfaceMethodDef{
			Name: p.curToken.Literal,
		}

		if !p.expectPeek(token.LParen) {
			return nil
		}

		params := p.parseFunctionParameters()
		ifaceMeth.Params = make([]string, len(params))
		for i, p := range params {
			ifaceMeth.Params[i] = p.String()
		}

		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}

		iface.Methods[ifaceMeth.Name] = ifaceMeth
		p.nextToken()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return iface
}
