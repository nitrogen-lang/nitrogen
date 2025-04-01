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
		expression.Condition = p.parseGroupedExpressionE()
		if !p.expectPeek(token.RParen) {
			return nil
		}
	} else {
		expression.Condition = p.parseGroupedExpressionE()
	}

	if p.peekTokenIs(token.Colon) {
		p.nextToken()
		expression.Consequence = p.parseSingleStmtBlock()
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return expression
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatements()

	if p.peekTokenIs(token.Elif) {
		p.peekToken.Type = token.If
		expression.Alternative = p.parseSingleStmtBlock()
	} else if p.peekTokenIs(token.Else) {
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
	loop := &ast.LoopStatement{
		Token: p.curToken,
	}
	expectClosingParen := false

	if p.peekTokenIs(token.LParen) {
		expectClosingParen = true
		p.nextToken()
	}

	if !p.peekTokenIs(token.Identifier) {
		p.peekError(token.Identifier)
		return nil
	}

	peekTok := p.peekToken
	p.nextToken()

	p.insertToken(token.Token{Type: token.Let, Literal: "let"})
	p.nextToken()

	if p.peekTokenIs(token.Comma) {
		p.curToken = peekTok
		loop := &ast.IterLoopStatement{Token: p.curToken}

		if !p.curTokenIs(token.Identifier) {
			p.addErrorWithCurPos("expected an ident, got %s", p.curToken.Type.String())
			return nil
		}

		loop.Key = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		p.nextToken() // Skip comma
		p.nextToken()

		if !p.curTokenIs(token.Identifier) {
			p.addErrorWithCurPos("expected an ident, got %s", p.curToken.Type.String())
			return nil
		}

		loop.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(token.In) {
			return nil
		}
		p.nextToken()

		val, ok := p.parseExpression(priLowest).(ast.Expression)
		if !ok {
			return nil
		}
		loop.Iter = val

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
	} else if p.peekTokenIs(token.In) {
		p.curToken = peekTok
		loop := &ast.IterLoopStatement{Token: p.curToken}

		if !p.curTokenIs(token.Identifier) {
			p.addErrorWithCurPos("expected an ident, got %s", p.curToken.Type.String())
			return nil
		}

		loop.Value = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

		if !p.expectPeek(token.In) {
			return nil
		}
		p.nextToken()

		val, ok := p.parseExpression(priLowest).(ast.Expression)
		if !ok {
			return nil
		}
		loop.Iter = val

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
	} else {
		p.insertToken(peekTok)

		loop.Init = p.parseDefStatement().(*ast.DefStatement)
		if !p.curTokenIs(token.Semicolon) {
			p.addErrorWithCurPos("expected semicolon, got %s", p.curToken.Type.String())
			return nil
		}
		p.nextToken()

		loop.Condition = p.parseExpression(priLowest).(ast.Expression)
		p.nextToken()
		if !p.curTokenIs(token.Semicolon) {
			p.addErrorWithCurPos("expected semicolon, got %s", p.curToken.Type.String())
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

func (p *Parser) parseWhileLoop() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseWhileLoop")
	}
	loop := &ast.LoopStatement{
		Token: p.curToken,
	}
	expectClosingParen := false

	if p.peekTokenIs(token.LParen) {
		expectClosingParen = true
		p.nextToken()
	}

	p.nextToken()
	condExp, ok := p.parseExpression(priLowest).(ast.Expression)
	if !ok {
		return nil
	}
	loop.Condition = condExp

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

func (p *Parser) parseInfiniteLoop() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseInfiniteLoop")
	}

	loop := &ast.LoopStatement{
		Token:     p.curToken,
		Init:      nil,
		Condition: nil,
		Iter:      nil,
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

func (p *Parser) parseMatchExpression() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseMatchExpression")
	}

	var matchExp ast.Expression

	if p.curTokenIs(token.LParen) {
		matchExp = p.parseGroupedExpressionE()
		if !p.expectPeek(token.RParen) {
			return nil
		}
	} else {
		matchExp = p.parseGroupedExpressionE()
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}
	p.nextToken()

	cases := make([]*ast.IfExpression, 0, 5)

	for {
		var caseLiteral ast.Expression
		if !p.curTokenIs(token.Underscore) {
			caseLiteral = p.parseBaseLiteral()
			if caseLiteral == nil {
				p.addErrorWithCurPos("expected a literal, got %s", p.curToken.Type.String())
				return nil
			}
		}

		if !p.expectPeek(token.Fatarrow) {
			return nil
		}
		p.nextToken()

		caseConsequence := p.parseSingleOrBlockStatements()
		if caseConsequence == nil {
			p.addErrorWithCurPos("expected a statement, got %s", p.curToken.Type.String())
			return nil
		}

		if !p.expectPeek(token.Comma) {
			return nil
		}

		theCase := &ast.IfExpression{
			Token:       p.curToken,
			Consequence: caseConsequence,
		}

		if caseLiteral != nil {
			theCase.Condition = &ast.InfixExpression{
				Token:    token.Token{Type: token.Equal, Literal: "=="},
				Left:     matchExp,
				Operator: "==",
				Right:    caseLiteral,
			}
		}

		cases = append(cases, theCase)

		if p.peekTokenIs(token.RBrace) {
			p.nextToken()
			break
		}
		p.nextToken()
	}

	var defaultCase *ast.BlockStatement

	caseChain := cases[0]
	currCase := caseChain

	for _, c := range cases[1:] {
		if c.Condition == nil {
			defaultCase = c.Consequence
			continue
		}

		currCase.Alternative = &ast.BlockStatement{
			Token:      c.Token,
			Statements: []ast.Statement{&ast.ExpressionStatement{Token: c.Token, Expression: c}},
		}
		currCase = c
	}

	if defaultCase != nil {
		currCase.Alternative = defaultCase
	}

	return caseChain
}
