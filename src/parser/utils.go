package parser

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) curTokenIs(t ...token.TokenType) bool {
	for _, tt := range t {
		if p.curToken.Type == tt {
			return true
		}
	}
	return false
}

func (p *Parser) peekTokenIs(t ...token.TokenType) bool {
	for _, tt := range t {
		if p.peekToken.Type == tt {
			return true
		}
	}
	return false
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

func IsIdent(s string) bool {
	return identRegex.MatchString(s)
}

func ImportName(path string) string {
	path = filepath.Base(path)
	dotIndex := strings.Index(path, ".")
	if dotIndex > -1 {
		path = path[:strings.Index(path, ".")]
	}

	if IsIdent(path) {
		return path
	}
	return ""
}
