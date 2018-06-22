package parser

import (
	"regexp"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if !p.peekTokenIs(t) {
		p.peekError(t)
		return false
	}

	p.nextToken()
	return true
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return priLowest
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return priLowest
}

func createKeywordToken(keyword string) token.Token {
	return token.Token{
		Type:    token.LookupIdent(keyword),
		Literal: keyword,
	}
}

var identRegex = regexp.MustCompile(`^[\p{L}_][\p{L}\d_]*$`)

func isIdent(s string) bool {
	return identRegex.MatchString(s)
}
