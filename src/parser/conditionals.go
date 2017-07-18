package parser

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	expression.Condition = p.parseGroupedExpression()
	if expression.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatements()

	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if !p.expectPeek(token.LBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatements()
	}

	return expression
}

func (p *Parser) parseForLoop() ast.Statement {
	loop := &ast.ForLoopStatement{}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	if !p.peekTokenIs(token.Identifier) {
		p.peekError(token.Identifier)
		return nil
	}

	p.insertToken(token.Token{Type: token.Let, Literal: "let"})
	p.nextToken()

	loop.Init = p.parseDefStatement().(*ast.DefStatement)
	if !p.curTokenIs(token.Semicolon) {
		p.addErrorWithPos("expected semicolon, got %s", p.curToken.Type.String())
		return nil
	}
	p.nextToken()

	loop.Condition = p.parseExpression(priLowest)
	p.nextToken()
	if !p.curTokenIs(token.Semicolon) {
		p.addErrorWithPos("expected semicolon, got %s", p.curToken.Type.String())
		return nil
	}
	p.nextToken()

	loop.Iter = p.parseExpression(priLowest)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	if !p.peekTokenIs(token.LBrace) {
		p.peekError(token.LBrace)
		return nil
	}

	p.nextToken()
	loop.Body = p.parseBlockStatements()
	p.nextToken()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return loop
}

func (p *Parser) parseCompareExpression(left ast.Expression) ast.Expression {
	c := p.curToken
	p.nextToken() // Go over OR, AND

	return &ast.CompareExpression{
		Token: c,
		Left:  left,
		Right: p.parseExpression(priLowest),
	}
}
