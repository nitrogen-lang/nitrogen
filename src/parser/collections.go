package parser

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.Array{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RSquare)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(priLowest)

		if !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(priLowest)

		hash.Pairs[key] = value

		if p.peekToken.Type == token.Semicolon {
			p.addErrorWithPos("Hash pairs must end with a comma")
			return nil
		}

		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			p.addError("Invalid hash literal")
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()

	exp.Index = p.parseExpression(priLowest)

	if !p.expectPeek(token.RSquare) {
		return nil
	}
	return exp
}
