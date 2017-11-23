package parser

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseArrayLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseArrayLiteral")
	}
	array := &ast.Array{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RSquare)
	return array
}

func (p *Parser) parseHashLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseHashLiteral")
	}
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
			p.addErrorWithPos("Invalid hash literal")
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseIndexExpression")
	}
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	if p.curTokenIs(token.Arrow) {
		p.nextToken()
		i := p.parseExpression(priAssign)
		// Only strings and identifiers are valid with the arrow syntax
		if _, ok := i.(*ast.StringLiteral); !ok {
			ident, ok := i.(*ast.Identifier)
			if !ok {
				p.addErrorWithPos("Index operator requires an identifier or string")
				return nil
			}
			// Convert identifier into a string for later lookup
			i = &ast.StringLiteral{
				Token: token.Token{
					Type:    token.String,
					Literal: ident.Value,
					Pos:     p.curToken.Pos,
				},
				Value: ident.Value,
			}
		}
		exp.Index = i
	} else {
		p.nextToken()
		exp.Index = p.parseExpression(priLowest)

		if !p.expectPeek(token.RSquare) {
			return nil
		}
	}
	return exp
}
