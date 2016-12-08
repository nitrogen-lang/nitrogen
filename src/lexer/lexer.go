package lexer

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

// TODO: Support Unicode by default
type Lexer struct {
	input  *bufio.Reader
	curCh  byte // current char under examination
	peekCh byte // peek character
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(reader)}
	// Populate both current and peek char
	l.readChar()
	l.readChar()
	return l
}

func NewString(input string) *Lexer {
	return New(strings.NewReader(input))
}

func (l *Lexer) readChar() {
	l.curCh = l.peekCh

	var err error
	l.peekCh, err = l.input.ReadByte()
	if err != nil {
		l.peekCh = 0
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.devourWhitespace()

	switch l.curCh {
	// Operators
	case '+':
		tok = newToken(token.PLUS, l.curCh)
	case '-':
		tok = newToken(token.MINUS, l.curCh)
	case '*':
		tok = newToken(token.ASTERISK, l.curCh)
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			tok = token.Token{
				Type:    token.COMMENT,
				Literal: l.readSingleLineComment(),
			}
		} else if l.peekChar() == '*' {
			l.readChar()
			tok = token.Token{
				Type:    token.COMMENT,
				Literal: l.readMultiLineComment(),
			}
		} else {
			tok = newToken(token.SLASH, l.curCh)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: "!=",
			}
		} else {
			tok = newToken(token.BANG, l.curCh)
		}

	// Equality
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{
				Type:    token.EQ,
				Literal: "==",
			}
		} else {
			tok = newToken(token.ASSIGN, l.curCh)
		}
	case '<':
		tok = newToken(token.LT, l.curCh)
	case '>':
		tok = newToken(token.GT, l.curCh)

	// Control characters
	case ',':
		tok = newToken(token.COMMA, l.curCh)
	case ';':
		tok = newToken(token.SEMICOLON, l.curCh)
	case ':':
		tok = newToken(token.COLON, l.curCh)

	// Groupings
	case '(':
		tok = newToken(token.LPAREN, l.curCh)
	case ')':
		tok = newToken(token.RPAREN, l.curCh)
	case '{':
		tok = newToken(token.LBRACE, l.curCh)
	case '}':
		tok = newToken(token.RBRACE, l.curCh)
	case '[':
		tok = newToken(token.LSQUARE, l.curCh)
	case ']':
		tok = newToken(token.RSQUARE, l.curCh)

	case '"':
		tok.Literal = l.readString()
		tok.Type = token.STRING
	case '#':
		tok.Literal = l.readSingleLineComment()
		tok.Type = token.COMMENT
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.curCh) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			return tok
		}

		tok = newToken(token.ILLEGAL, l.curCh)
	}

	l.readChar()
	return tok
}

func (l *Lexer) peekChar() byte {
	return l.peekCh
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isLetter(l.curCh) {
		ident.WriteByte(l.curCh)
		l.readChar()
	}
	return ident.String()
}

// TODO: Support escape sequences, standard Go should be fine, or PHP.
func (l *Lexer) readString() string {
	var ident bytes.Buffer
	l.readChar() // Go past the starting double quote

	for l.curCh != '"' {
		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return ident.String()
}

// TODO: Support arbitrary based numbers [base]b[number]
// base defaults to 10. 8b10 = 8 in octal 16b10 = 16 in hex
func (l *Lexer) readNumber() token.Token {
	var ident bytes.Buffer
	tokenType := token.INT

	for isDigit(l.curCh) {
		// The parser will handle bad floats
		if l.curCh == '.' && tokenType == token.INT {
			tokenType = token.FLOAT
		}

		ident.WriteByte(l.curCh)
		l.readChar()
	}

	return token.Token{
		Type:    token.TokenType(tokenType),
		Literal: ident.String(),
	}
}

func (l *Lexer) readSingleLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over # or / characters

	for l.curCh != '\n' {
		com.WriteByte(l.curCh)
		l.readChar()
	}
	return strings.TrimSpace(com.String())
}

func (l *Lexer) readMultiLineComment() string {
	var com bytes.Buffer
	l.readChar() // Go over * character

	for l.curCh != 0 {
		if l.curCh == '*' && l.peekChar() == '/' {
			l.readChar() // Skip *
			break
		}

		com.WriteByte(l.curCh)
		l.readChar()
	}
	return com.String()
}

func (l *Lexer) devourWhitespace() {
	for isWhitespace(l.curCh) {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ('0' <= ch && ch <= '9') || ch == '.'
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
