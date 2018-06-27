package parser

import (
	"fmt"
	"strconv"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseIdentifier() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseIdentifier")
	}
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseIntegerLiteral")
	}
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.addErrorWithPos("Invalid integer: %q", p.curToken.Literal)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseFloatLiteral")
	}
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.addErrorWithPos("Invalid float: %q", p.curToken.Literal)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseNullLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseNullLiteral")
	}
	return &ast.NullLiteral{Token: p.curToken}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseStringLiteral")
	}
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: []rune(p.curToken.Literal),
	}
}

func (p *Parser) parseBoolean() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseBoolean")
	}
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.True)}
}
