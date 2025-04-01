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
		keyExp, ok := key.(ast.Expression)

		if !ok || !p.expectPeek(token.Colon) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(priLowest)
		valueExp, ok := value.(ast.Expression)
		if !ok {
			return nil
		}

		hash.Pairs[keyExp] = valueExp

		if p.peekToken.Type == token.Semicolon {
			p.addErrorWithCurPos("Hash pairs must end with a comma")
			return nil
		}

		if !p.peekTokenIs(token.RBrace) && !p.expectPeek(token.Comma) {
			p.addErrorWithCurPos("Invalid hash literal")
			return nil
		}
	}

	if !p.expectPeek(token.RBrace) {
		return nil
	}

	return hash
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseIndexExpression")
	}
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(priLowest).(ast.Expression)

	if !p.expectPeek(token.RSquare) {
		return nil
	}
	return exp
}

func (p *Parser) parseAttributeExpression(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseAttributeExpression")
	}
	exp := &ast.AttributeExpression{Token: p.curToken, Left: left}

	p.nextToken()
	i := p.parseExpression(priAssign)

	ident, ok := i.(*ast.Identifier)
	if !ok {
		p.addErrorWithCurPos("Attribute operator requires an identifier")
		return nil
	}
	// Convert identifier into a string for later lookup
	exp.Index = &ast.StringLiteral{
		Token: token.Token{
			Type:    token.String,
			Literal: ident.Value,
			Pos:     p.curToken.Pos,
		},
		Value: []rune(ident.Value),
	}

	return exp
}
