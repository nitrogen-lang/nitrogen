package parser

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseIfExpression() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseIfExpression")
	}
	expression := &ast.IfExpression{Token: p.curToken}

	if p.curTokenIs(token.LParen) {
		p.nextToken()
		expression.Condition = p.parseGroupedExpression()
		if !p.expectPeek(token.LBrace) {
			return nil
		}
	} else {
		expression.Condition = p.parseGroupedExpressionC(token.LBrace)
	}

	if expression.Condition == nil {
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
	if p.settings.Debug {
		fmt.Println("parseForLoop")
	}
	loop := &ast.ForLoopStatement{}
	expectClosingParen := false

	if p.peekTokenIs(token.LParen) {
		expectClosingParen = true
		p.nextToken()
	}

	if p.peekTokenIs(token.LBrace) || (expectClosingParen && p.peekTokenIs(token.RParen)) {
		loop.Init = nil
		loop.Condition = nil
		loop.Iter = nil
	} else {
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

		loop.Condition = p.parseExpression(priLowest).(ast.Expression)
		p.nextToken()
		if !p.curTokenIs(token.Semicolon) {
			p.addErrorWithPos("expected semicolon, got %s", p.curToken.Type.String())
			return nil
		}
		p.nextToken()

		loop.Iter = p.parseExpression(priLowest)
	}

	if expectClosingParen && !p.expectPeek(token.RParen) {
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

func (p *Parser) parseCompareExpression(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseCompareExpression")
	}
	c := p.curToken
	p.nextToken() // Go over OR, AND

	return &ast.CompareExpression{
		Token: c,
		Left:  left,
		Right: p.parseExpression(priLowest).(ast.Expression),
	}
}

func (p *Parser) parseTryCatch() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseTryCatch")
	}
	if !p.expectPeek(token.LBrace) {
		return nil
	}

	try := p.parseBlockStatements()

	if !p.expectPeek(token.Catch) {
		return nil
	}

	var symbol *ast.Identifier
	if p.peekTokenIs(token.Identifier) {
		p.nextToken()
		symbol = p.parseIdentifier().(*ast.Identifier)
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	catch := p.parseBlockStatements()

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return &ast.TryCatchExpression{
		Try:    try,
		Catch:  catch,
		Symbol: symbol,
	}
}
