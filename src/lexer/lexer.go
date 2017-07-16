package lexer

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type Lexer struct {
	input     *bufio.Reader
	curCh     rune // current char under examination
	peekCh    rune // peek character
	lastToken token.Token

	fileList  []string
	line, col int
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(reader)}
	// Populate both current and peek char
	l.readRune()
	l.readRune()
	l.col = 1
	l.line = 1
	return l
}

func NewFileList(files []string) (*Lexer, error) {
	l := &Lexer{fileList: files}
	if err := l.loadFile(); err != nil {
		return nil, err
	}
	l.readRune()
	l.readRune()
	l.col = 1
	l.line = 1
	return l, nil
}

func NewString(input string) *Lexer {
	return New(strings.NewReader(input))
}

var errFailedLoadingFile = errors.New("Failed to open file")

func (l *Lexer) loadFile() error {
	if len(l.fileList) == 0 {
		panic("No more files to load")
	}

	nextFile := l.fileList[0]
	l.fileList = l.fileList[1:]

	file, err := os.Open(nextFile)
	if err != nil {
		return errFailedLoadingFile
	}

	l.input = bufio.NewReader(file)
	return nil
}

func (l *Lexer) readRune() {
	oldPeek := l.peekCh

	newPeek, _, err := l.input.ReadRune()
	if err != nil {
		if len(l.fileList) > 0 {
			if err := l.loadFile(); err != nil {
				panic("Failed opening file")
			}
			l.readRune()
			return
		}
		l.peekCh = 0
		l.curCh = oldPeek
		return
	}
	l.curCh = oldPeek
	l.peekCh = newPeek
	l.col++
}

func makePos(line, col int) token.Position {
	return token.Position{
		Line: line,
		Col:  col,
	}
}

func (l *Lexer) curPosition() token.Position {
	return makePos(l.line, l.col)
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.devourWhitespaceNotNewLine()

	switch l.curCh {
	case '\n':
		if l.needSemicolon() {
			tok = l.newToken(token.SEMICOLON, ';')
			l.resetPos()
		} else {
			l.devourWhitespace()
			return l.NextToken()
		}

	// Operators
	case '+':
		tok = l.newToken(token.PLUS, l.curCh)
	case '-':
		tok = l.newToken(token.MINUS, l.curCh)
	case '*':
		tok = l.newToken(token.ASTERISK, l.curCh)
	case '%':
		tok = l.newToken(token.MODULUS, l.curCh)
	case '/':
		if l.peekChar() == '/' {
			l.readRune()
			tok = l.readSingleLineComment()
			l.resetPos()
		} else if l.peekChar() == '*' {
			l.readRune()
			tok = l.readMultiLineComment()
		} else {
			tok = l.newToken(token.SLASH, l.curCh)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readRune()
			tok = token.Token{
				Type:    token.NOT_EQ,
				Literal: "!=",
				Pos:     l.curPosition(),
			}
			tok.Pos.Col -= 1
		} else {
			tok = l.newToken(token.BANG, l.curCh)
		}

	// Equality
	case '=':
		if l.peekChar() == '=' {
			l.readRune()
			tok = token.Token{
				Type:    token.EQ,
				Literal: "==",
				Pos:     l.curPosition(),
			}
			tok.Pos.Col -= 1
		} else {
			tok = l.newToken(token.ASSIGN, l.curCh)
		}
	case '<':
		tok = l.newToken(token.LT, l.curCh)
	case '>':
		tok = l.newToken(token.GT, l.curCh)

	// Control characters
	case ',':
		tok = l.newToken(token.COMMA, l.curCh)
	case ';':
		tok = l.newToken(token.SEMICOLON, l.curCh)
	case ':':
		tok = l.newToken(token.COLON, l.curCh)

	// Groupings
	case '(':
		tok = l.newToken(token.LPAREN, l.curCh)
	case ')':
		tok = l.newToken(token.RPAREN, l.curCh)
	case '{':
		tok = l.newToken(token.LBRACE, l.curCh)
	case '}':
		tok = l.newToken(token.RBRACE, l.curCh)
	case '[':
		tok = l.newToken(token.LSQUARE, l.curCh)
	case ']':
		tok = l.newToken(token.RSQUARE, l.curCh)

	case '"':
		tok = l.readString()
	case '\'':
		tok = l.readRawString()
	case '#':
		tok = l.readSingleLineComment()
		l.resetPos()
	case '\\':
		if l.peekCh == 'x' {
			l.readRune()
			tok = l.readNumber()
		} else {
			tok = l.newToken(token.ILLEGAL, l.curCh)
		}
	case 0:
		if l.needSemicolon() {
			tok = l.newToken(token.SEMICOLON, ';')
			l.resetPos()
		} else {
			tok.Literal = ""
			tok.Type = token.EOF
			tok.Pos = l.curPosition()
		}

	default:
		if isLetter(l.curCh) {
			tok.Pos = l.curPosition()
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			l.lastToken = tok
			return tok
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			l.lastToken = tok
			return tok
		}

		tok = l.newToken(token.ILLEGAL, l.curCh)
	}

	l.readRune()
	l.lastToken = tok
	return tok
}

func (l *Lexer) peekChar() rune {
	return l.peekCh
}

func (l *Lexer) resetPos() {
	l.line++
	l.col = 0
}

func (l *Lexer) needSemicolon() bool {
	return l.lastTokenWas(
		token.IDENT,
		token.INT,
		token.FLOAT,
		token.STRING,
		token.NULL,
		token.RETURN,
		token.RPAREN,
		token.RSQUARE,
		token.RBRACE) && !l.lastTokenWas(token.SEMICOLON)
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isIdent(l.curCh) || isDigit(l.curCh) {
		ident.WriteRune(l.curCh)
		l.readRune()
	}
	return ident.String()
}

func (l *Lexer) readString() token.Token {
	var ident bytes.Buffer
	pos := l.curPosition()
	l.readRune() // Go past the starting double quote

	for l.curCh != '"' {
		if l.curCh == '\n' {
			return token.Token{
				Literal: "Newline not allowed in string",
				Type:    token.ILLEGAL,
				Pos:     l.curPosition(),
			}
		}

		if l.curCh == '\\' {
			l.readRune()
			switch l.curCh {
			case 'b': // backspace
				ident.WriteRune('\b')
			case 'n': // newline
				ident.WriteRune('\n')
			case 'r': // carriage return
				ident.WriteRune('\r')
			case 't': // horizontal tab
				ident.WriteRune('\t')
			case 'v': // vertical tab
				ident.WriteRune('\v')
			case 'f': // form feed
				ident.WriteRune('\f')
			case '\\': // back slash
				ident.WriteRune('\\')
			case '"': // double quote
				ident.WriteRune('"')
			default:
				ident.WriteByte('\\')
				ident.WriteRune(l.curCh)
			}
			l.readRune()
			continue
		}
		ident.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Literal: ident.String(),
		Type:    token.STRING,
		Pos:     pos,
	}
}

func (l *Lexer) readRawString() token.Token {
	var ident bytes.Buffer
	pos := l.curPosition()
	l.readRune() // Go past the starting double quote

	for l.curCh != '\'' {
		if l.curCh == '\\' && l.peekCh == '\'' {
			l.readRune() // Go past backslash so the next line will write a single quote
		}
		ident.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Literal: ident.String(),
		Type:    token.STRING,
		Pos:     pos,
	}
}

func (l *Lexer) readNumber() token.Token {
	var number bytes.Buffer
	pos := l.curPosition()
	base := ""
	tokenType := token.INT

	if l.curCh == 'x' {
		base = "0x"
		pos.Col-- // Correct for initial \
		l.readRune()
	}

	if l.curCh == '.' {
		l.readRune()
		return token.Token{
			Type:    token.ILLEGAL,
			Literal: "Invalid float literal",
			Pos:     pos,
		}
	}

	for isDigit(l.curCh) || isHexDigit(l.curCh) {
		if l.curCh == '.' {
			if tokenType != token.INT {
				return token.Token{
					Type:    token.ILLEGAL,
					Literal: "Invalid float literal",
					Pos:     pos,
				}
			}
			tokenType = token.FLOAT
		}

		number.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Type:    token.TokenType(tokenType),
		Literal: base + number.String(),
		Pos:     pos,
	}
}

func (l *Lexer) readSingleLineComment() token.Token {
	var com bytes.Buffer
	pos := l.curPosition()
	if l.curCh == '/' {
		pos.Col-- // Correct column for inital /
	}
	l.readRune() // Go over # or / characters

	for l.curCh != '\n' {
		com.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Literal: strings.TrimSpace(com.String()),
		Type:    token.COMMENT,
		Pos:     pos,
	}
}

func (l *Lexer) readMultiLineComment() token.Token {
	var com bytes.Buffer
	pos := l.curPosition()
	pos.Col--    // Correct column for initial /
	l.readRune() // Go over * character

	for l.curCh != 0 {
		if l.curCh == '*' && l.peekChar() == '/' {
			l.readRune() // Skip *
			break
		}

		if l.curCh == '\n' {
			l.resetPos()
		}

		com.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Literal: com.String(),
		Type:    token.COMMENT,
		Pos:     pos,
	}
}

func (l *Lexer) devourWhitespace() {
	for l.isWhitespace(l.curCh) {
		if l.curCh == '\n' {
			l.resetPos()
		}
		l.readRune()
	}
}

func (l *Lexer) devourWhitespaceNotNewLine() {
	for l.curCh != '\n' && l.isWhitespace(l.curCh) {
		l.readRune()
	}
}

func (l *Lexer) newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
		Pos:     l.curPosition(),
	}
}

// Identifiers must start with a letter
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || unicode.IsLetter(ch)
}

// After the first letter, an ident can be a letter or number
func isIdent(ch rune) bool {
	return isLetter(ch) || (ch != '.' && isDigit(ch)) // A period is not a valid identifier name
}

// Only Latin numbers
func isDigit(ch rune) bool {
	return ('0' <= ch && ch <= '9') || ch == '.'
}

func isHexDigit(ch rune) bool {
	return ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func (l *Lexer) isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) lastTokenWas(types ...token.TokenType) bool {
	for _, t := range types {
		if l.lastToken.Type == t {
			return true
		}
	}
	return false
}
