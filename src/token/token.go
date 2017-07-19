package token

import "strconv"

type TokenType int

// Position represents the line and column number where a token starts
// in a source file.
type Position struct {
	Line, Col int
}

// TODO: Add filename to token
type Token struct {
	Type    TokenType
	Literal string
	Pos     Position
}

// All tokens in Nitrogen
const (
	Illegal TokenType = iota + 1
	EOF
	Comment
	EOL

	// Identifiers & literals
	Identifier
	Integer
	Float
	String

	// Operators
	Assign
	Plus
	Dash
	Bang
	Asterisk
	Slash
	Modulo

	PlusAssign
	MinusAssign
	TimesAssign
	SlashAssign

	LessThan
	GreaterThan
	LessThanEq
	GreaterThanEq
	Equal
	NotEqual

	// Delimiters
	Comma
	Semicolon
	Colon

	// Groups and blocks
	LParen
	RParen
	LBrace
	RBrace
	LSquare
	RSquare

	// Keywords
	keywordBeg
	LAnd
	LOr

	Function
	Let
	Always
	True
	False
	If
	Else
	Return
	Nil
	For
	Continue
	Break
	keywordEnd
)

var tokens = [...]string{
	Illegal: "ILLEGAL",
	EOF:     "EOF",
	Comment: "COMMENT",
	EOL:     "EOL",

	// Identifiers & literals
	Identifier: "IDENT",
	Integer:    "INT",
	Float:      "FLOAT",
	String:     "STRING",

	// Operators
	Assign:   "=",
	Plus:     "+",
	Dash:     "-",
	Bang:     "!",
	Asterisk: "*",
	Slash:    "/",
	Modulo:   "%",

	PlusAssign:  "+=",
	MinusAssign: "-=",
	TimesAssign: "*=",
	SlashAssign: "/=",

	LessThan:      "<",
	GreaterThan:   ">",
	LessThanEq:    "<=",
	GreaterThanEq: ">=",
	Equal:         "==",
	NotEqual:      "!=",

	// Delimiters
	Comma:     ",",
	Semicolon: ";",
	Colon:     ":",

	// Groups and blocks
	LParen:  "(",
	RParen:  ")",
	LBrace:  "{",
	RBrace:  "}",
	LSquare: "[",
	RSquare: "]",

	// Keywords
	Function: "func",
	Let:      "let",
	Always:   "always",
	True:     "true",
	False:    "false",
	If:       "if",
	Else:     "else",
	Return:   "return",
	Nil:      "nil",
	LAnd:     "and",
	LOr:      "or",
	For:      "for",
	Continue: "continue",
	Break:    "break",
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Identifier
}

func (t TokenType) String() string {
	if 0 <= t && t < TokenType(len(tokens)) {
		return tokens[t]
	}
	return "token(" + strconv.Itoa(int(t)) + ")"
}
