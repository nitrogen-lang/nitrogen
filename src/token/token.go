package token

import "strconv"

type TokenType int

// Position represents the line and column number where a token starts
// in a source file.
type Position struct {
	Line, Col uint
	Filename  string
}

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
	ByteString

	// Operators
	Assign
	Plus
	Dash
	Bang
	Asterisk
	Slash
	Modulo
	Dot

	PlusAssign
	MinusAssign
	TimesAssign
	SlashAssign
	ModAssign

	LessThan
	GreaterThan
	LessThanEq
	GreaterThanEq
	Equal
	NotEqual

	BitwiseAnd
	BitwiseOr
	BitwiseAndNot

	ShiftLeft
	ShiftRight

	// Delimiters
	Comma
	Semicolon
	Colon
	Carrot
	Fatarrow
	Underscore

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
	Const
	True
	False
	If
	Elif
	Else
	Return
	Nil
	For
	While
	Loop
	Continue
	Break
	Recover
	Class
	New
	Pass
	Import
	As
	Delete
	Use
	Native
	Do
	In
	Interface
	Implements
	Breakpoint
	Match
	Export
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
	ByteString: "BYTES",

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
	ModAssign:   "%=",

	LessThan:      "<",
	GreaterThan:   ">",
	LessThanEq:    "<=",
	GreaterThanEq: ">=",
	Equal:         "==",
	NotEqual:      "!=",

	BitwiseAnd:    "&",
	BitwiseOr:     "|",
	BitwiseAndNot: "&^",

	ShiftLeft:  "<<",
	ShiftRight: ">>",

	// Delimiters
	Comma:      ",",
	Semicolon:  ";",
	Colon:      ":",
	Carrot:     "^",
	Fatarrow:   "=>",
	Underscore: "_",

	// Groups and blocks
	LParen:  "(",
	RParen:  ")",
	LBrace:  "{",
	RBrace:  "}",
	LSquare: "[",
	RSquare: "]",

	// Keywords
	Function:   "fn",
	Let:        "let",
	Const:      "const",
	True:       "true",
	False:      "false",
	If:         "if",
	Elif:       "elif",
	Else:       "else",
	Return:     "return",
	Nil:        "nil",
	LAnd:       "and",
	LOr:        "or",
	For:        "for",
	While:      "while",
	Loop:       "loop",
	Continue:   "continue",
	Break:      "break",
	Recover:    "recover",
	Class:      "class",
	New:        "new",
	Pass:       "pass",
	Import:     "import",
	As:         "as",
	Delete:     "delete",
	Use:        "use",
	Native:     "native",
	Do:         "do",
	In:         "in",
	Interface:  "interface",
	Implements: "implements",
	Breakpoint: "breakpoint",
	Match:      "match",
	Export:     "export",
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
