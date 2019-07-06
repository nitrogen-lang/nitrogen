package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

type Lexer struct {
	input     *bufio.Reader
	curCh     rune // current char under examination
	peekCh    rune // peek character
	lastToken token.Token

	fileList    []string
	line, col   uint
	currentFile string
}

func New(reader io.Reader) *Lexer {
	l := &Lexer{input: bufio.NewReader(reader)}
	// Populate both current and peek char
	l.readRune()
	l.readRune()
	l.col = 1
	l.line = 1
	l.currentFile = "anonymous"
	return l
}

func NewFile(file string) (*Lexer, error) {
	return NewFileList([]string{file})
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

func (l *Lexer) loadFile() error {
	if len(l.fileList) == 0 {
		panic("No more files to load")
	}

	nextFile := l.fileList[0]
	l.fileList = l.fileList[1:]

	file, err := os.Open(nextFile)
	if err != nil {
		return fmt.Errorf("Failed to open file %s", nextFile)
	}

	l.input = bufio.NewReader(file)
	l.currentFile, _ = filepath.Abs(nextFile)
	return nil
}

func (l *Lexer) readRune() {
	oldPeek := l.peekCh

	newPeek, _, err := l.input.ReadRune()
	if err != nil {
		if len(l.fileList) > 0 {
			if err := l.loadFile(); err != nil {
				panic(err)
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

func makePos(line, col uint) token.Position {
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
			tok = l.newToken(token.Semicolon, ';')
			l.resetPos()
		} else {
			l.devourWhitespace()
			return l.NextToken()
		}

	// Operators
	case '+':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.PlusAssign,
				Literal:  "+=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Plus, l.curCh)
		}
	case '-':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.MinusAssign,
				Literal:  "-=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Dash, l.curCh)
		}
	case '*':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.TimesAssign,
				Literal:  "*=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Asterisk, l.curCh)
		}
	case '%':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.ModAssign,
				Literal:  "%=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Modulo, l.curCh)
		}
	case '/':
		switch l.peekChar() {
		case '/':
			l.readRune()
			tok = l.readSingleLineComment()
			l.resetPos()
		case '*':
			l.readRune()
			tok = l.readMultiLineComment()
		case '=':
			tok = token.Token{
				Type:     token.SlashAssign,
				Literal:  "/=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		default:
			tok = l.newToken(token.Slash, l.curCh)
		}
	case '!':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.NotEqual,
				Literal:  "!=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Bang, l.curCh)
		}

	// Equality
	case '=':
		if l.peekChar() == '=' {
			tok = token.Token{
				Type:     token.Equal,
				Literal:  "==",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.Assign, l.curCh)
		}
	case '<':
		switch l.peekChar() {
		case '=':
			tok = token.Token{
				Type:     token.LessThanEq,
				Literal:  "<=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		case '<':
			tok = token.Token{
				Type:     token.ShiftLeft,
				Literal:  "<<",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		default:
			tok = l.newToken(token.LessThan, l.curCh)
		}
	case '>':
		switch l.peekChar() {
		case '=':
			tok = token.Token{
				Type:     token.GreaterThanEq,
				Literal:  ">=",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		case '>':
			tok = token.Token{
				Type:     token.ShiftRight,
				Literal:  ">>",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		default:
			tok = l.newToken(token.GreaterThan, l.curCh)
		}

	case '&':
		if l.peekChar() == '^' {
			tok = token.Token{
				Type:     token.BitwiseAndNot,
				Literal:  "&^",
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
			l.readRune()
		} else {
			tok = l.newToken(token.BitwiseAnd, l.curCh)
		}
	case '|':
		tok = l.newToken(token.BitwiseOr, l.curCh)

	// Control characters
	case ',':
		tok = l.newToken(token.Comma, l.curCh)
	case ';':
		tok = l.newToken(token.Semicolon, l.curCh)
	case ':':
		tok = l.newToken(token.Colon, l.curCh)
	case '.':
		tok = l.newToken(token.Dot, l.curCh)
	case '^':
		tok = l.newToken(token.Carrot, l.curCh)

	// Groupings
	case '(':
		tok = l.newToken(token.LParen, l.curCh)
	case ')':
		tok = l.newToken(token.RParen, l.curCh)
	case '{':
		tok = l.newToken(token.LBrace, l.curCh)
	case '}':
		tok = l.newToken(token.RBrace, l.curCh)
	case '[':
		tok = l.newToken(token.LSquare, l.curCh)
	case ']':
		tok = l.newToken(token.RSquare, l.curCh)

	case '"':
		tok = l.readString()
	case '\'':
		tok = l.readRawString()
	case '#':
		tok = l.readSingleLineComment()
		l.resetPos()
	case 0:
		if l.needSemicolon() {
			tok = l.newToken(token.Semicolon, ';')
			l.resetPos()
		} else {
			tok.Literal = ""
			tok.Type = token.EOF
			tok.Pos = l.curPosition()
			tok.Filename = l.currentFile
		}

	default:
		if isLetter(l.curCh) {
			tok.Pos = l.curPosition()
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			tok.Filename = l.currentFile
			l.lastToken = tok
			return tok
		} else if isDigit(l.curCh) {
			tok = l.readNumber()
			l.lastToken = tok
			return tok
		}

		tok = l.newToken(token.Illegal, l.curCh)
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
		token.Identifier,
		token.Integer,
		token.Float,
		token.String,
		token.True,
		token.False,
		token.Nil,
		token.Return,
		token.Break,
		token.Continue,
		token.RParen,
		token.RSquare,
		token.RBrace) && !l.lastTokenWas(token.Semicolon)
}

func (l *Lexer) readIdentifier() string {
	var ident bytes.Buffer
	for isIdent(l.curCh) {
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
				Literal:  "Newline not allowed in string",
				Type:     token.Illegal,
				Pos:      l.curPosition(),
				Filename: l.currentFile,
			}
		}

		if l.curCh == '\\' {
			l.readRune()
			switch l.curCh {
			case '0': // null
				ident.WriteRune(0)
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
			case 'e': // escape
				ident.WriteRune('\033')
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
		Literal:  ident.String(),
		Type:     token.String,
		Pos:      pos,
		Filename: l.currentFile,
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
		Literal:  ident.String(),
		Type:     token.String,
		Pos:      pos,
		Filename: l.currentFile,
	}
}

func (l *Lexer) readNumber() token.Token {
	var number bytes.Buffer
	pos := l.curPosition()
	base := ""
	tokenType := token.Integer

	if l.curCh == '0' {
		switch l.peekChar() {
		case 'x':
			base = "0x"
			l.readRune()
			l.readRune()
		case 'b':
			base = "0b"
			l.readRune()
			l.readRune()
		case 'o':
			base = "0o"
			l.readRune()
			l.readRune()
		}
	}

	if l.curCh == '.' {
		l.readRune()
		return token.Token{
			Type:     token.Illegal,
			Literal:  "Invalid float literal",
			Pos:      pos,
			Filename: l.currentFile,
		}
	}

	for isDigit(l.curCh) || isHexDigit(l.curCh) || l.curCh == '_' {
		if l.curCh == '_' {
			l.readRune()
			continue
		} else if l.curCh == '.' {
			if tokenType != token.Integer {
				return token.Token{
					Type:     token.Illegal,
					Literal:  "Invalid float literal",
					Pos:      pos,
					Filename: l.currentFile,
				}
			}
			tokenType = token.Float
		}

		number.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Type:     token.TokenType(tokenType),
		Literal:  base + number.String(),
		Pos:      pos,
		Filename: l.currentFile,
	}
}

func (l *Lexer) readSingleLineComment() token.Token {
	var com bytes.Buffer
	pos := l.curPosition()
	if l.curCh == '/' {
		pos.Col-- // Correct column for inital /
	}
	l.readRune() // Go over # or / characters

	for l.curCh != '\n' && l.curCh != 0 {
		com.WriteRune(l.curCh)
		l.readRune()
	}

	return token.Token{
		Literal:  strings.TrimSpace(com.String()),
		Type:     token.Comment,
		Pos:      pos,
		Filename: l.currentFile,
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
		Literal:  com.String(),
		Type:     token.Comment,
		Pos:      pos,
		Filename: l.currentFile,
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
		Type:     tokenType,
		Literal:  string(ch),
		Pos:      l.curPosition(),
		Filename: l.currentFile,
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
