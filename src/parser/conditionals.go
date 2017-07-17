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

func (p *Parser) parseCompareExpression(left ast.Expression) ast.Expression {
	c := p.curToken
	p.nextToken() // Go over OR, AND

	return &ast.CompareExpression{
		Token: c,
		Left:  left,
		Right: p.parseExpression(priLowest),
	}
}
