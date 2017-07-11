package token

// TODO: Replace this with an enum int?
type TokenType string

// TODO: Add line, column, and filename to token
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"
	EOL     = "EOL" // End of Line

	// Identifiers & literals
	IDENT  = "IDENT"  // add, foobar, x, y
	INT    = "INT"    // 1343546
	FLOAT  = "FLOAT"  // 12.52
	STRING = "STRING" // "some text"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT     = "<"
	GT     = ">"
	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	// Groups and blocks
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	LSQUARE = "["
	RSQUARE = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	DEF      = "DEF"
	CONST    = "ALWAYS"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	NULL     = "NULL"
)

var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"let":    DEF,
	"always": CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"nil":    NULL,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
