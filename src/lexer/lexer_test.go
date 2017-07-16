package lexer

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
always ten = 10;

世界

ident2

6 % 3

let add = func(x, y) {
	x + y;
}

let result = add(five, ten)
!-/5*;
6 < 10 > 7

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10
10 != 9
"foobar"
"foo bar"
[1, 2];
{"foo": "bar"}
nil
# Single line hash comment
// Single line slash comment
/* Multi-line
comment */

12.5
12.5.7

"\n\r\t\v\f\\\"\b\' Hello"
'\n\r\t\v\f\\\"\b\' Hello'

\x2F54a
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedPos     token.Position
	}{
		{token.DEF, "let", makePos(1, 1)},
		{token.IDENT, "five", makePos(1, 5)},
		{token.ASSIGN, "=", makePos(1, 10)},
		{token.INT, "5", makePos(1, 12)},
		{token.SEMICOLON, ";", makePos(1, 13)},

		{token.CONST, "always", makePos(2, 1)},
		{token.IDENT, "ten", makePos(2, 8)},
		{token.ASSIGN, "=", makePos(2, 12)},
		{token.INT, "10", makePos(2, 14)},
		{token.SEMICOLON, ";", makePos(2, 16)},

		{token.IDENT, "世界", makePos(4, 1)},
		{token.SEMICOLON, ";", makePos(4, 3)},

		{token.IDENT, "ident2", makePos(6, 1)},
		{token.SEMICOLON, ";", makePos(6, 7)},

		{token.INT, "6", makePos(8, 1)},
		{token.MODULUS, "%", makePos(8, 3)},
		{token.INT, "3", makePos(8, 5)},
		{token.SEMICOLON, ";", makePos(8, 6)},

		{token.DEF, "let", makePos(10, 1)},
		{token.IDENT, "add", makePos(10, 5)},
		{token.ASSIGN, "=", makePos(10, 9)},
		{token.FUNCTION, "func", makePos(10, 11)},
		{token.LPAREN, "(", makePos(10, 15)},
		{token.IDENT, "x", makePos(10, 16)},
		{token.COMMA, ",", makePos(10, 17)},
		{token.IDENT, "y", makePos(10, 19)},
		{token.RPAREN, ")", makePos(10, 20)},
		{token.LBRACE, "{", makePos(10, 22)},
		{token.IDENT, "x", makePos(11, 2)},
		{token.PLUS, "+", makePos(11, 4)},
		{token.IDENT, "y", makePos(11, 6)},
		{token.SEMICOLON, ";", makePos(11, 7)},
		{token.RBRACE, "}", makePos(12, 1)},
		{token.SEMICOLON, ";", makePos(12, 2)},

		{token.DEF, "let", makePos(14, 1)},
		{token.IDENT, "result", makePos(14, 5)},
		{token.ASSIGN, "=", makePos(14, 12)},
		{token.IDENT, "add", makePos(14, 14)},
		{token.LPAREN, "(", makePos(14, 17)},
		{token.IDENT, "five", makePos(14, 18)},
		{token.COMMA, ",", makePos(14, 22)},
		{token.IDENT, "ten", makePos(14, 24)},
		{token.RPAREN, ")", makePos(14, 27)},
		{token.SEMICOLON, ";", makePos(14, 28)},

		{token.BANG, "!", makePos(15, 1)},
		{token.MINUS, "-", makePos(15, 2)},
		{token.SLASH, "/", makePos(15, 3)},
		{token.INT, "5", makePos(15, 4)},
		{token.ASTERISK, "*", makePos(15, 5)},
		{token.SEMICOLON, ";", makePos(15, 6)},

		{token.INT, "6", makePos(16, 1)},
		{token.LT, "<", makePos(16, 3)},
		{token.INT, "10", makePos(16, 5)},
		{token.GT, ">", makePos(16, 8)},
		{token.INT, "7", makePos(16, 10)},
		{token.SEMICOLON, ";", makePos(16, 11)},

		{token.IF, "if", makePos(18, 1)},
		{token.LPAREN, "(", makePos(18, 4)},
		{token.INT, "5", makePos(18, 5)},
		{token.LT, "<", makePos(18, 7)},
		{token.INT, "10", makePos(18, 9)},
		{token.RPAREN, ")", makePos(18, 11)},
		{token.LBRACE, "{", makePos(18, 13)},
		{token.RETURN, "return", makePos(19, 2)},
		{token.TRUE, "true", makePos(19, 9)},
		{token.SEMICOLON, ";", makePos(19, 13)},
		{token.RBRACE, "}", makePos(20, 1)},
		{token.ELSE, "else", makePos(20, 3)},
		{token.LBRACE, "{", makePos(20, 8)},
		{token.RETURN, "return", makePos(21, 2)},
		{token.FALSE, "false", makePos(21, 9)},
		{token.SEMICOLON, ";", makePos(21, 14)},
		{token.RBRACE, "}", makePos(22, 1)},
		{token.SEMICOLON, ";", makePos(22, 2)},

		{token.INT, "10", makePos(24, 1)},
		{token.EQ, "==", makePos(24, 4)},
		{token.INT, "10", makePos(24, 7)},
		{token.SEMICOLON, ";", makePos(24, 9)},

		{token.INT, "10", makePos(25, 1)},
		{token.NOT_EQ, "!=", makePos(25, 4)},
		{token.INT, "9", makePos(25, 7)},
		{token.SEMICOLON, ";", makePos(25, 8)},

		{token.STRING, "foobar", makePos(26, 1)},
		{token.SEMICOLON, ";", makePos(26, 9)},

		{token.STRING, "foo bar", makePos(27, 1)},
		{token.SEMICOLON, ";", makePos(27, 10)},

		{token.LSQUARE, "[", makePos(28, 1)},
		{token.INT, "1", makePos(28, 2)},
		{token.COMMA, ",", makePos(28, 3)},
		{token.INT, "2", makePos(28, 5)},
		{token.RSQUARE, "]", makePos(28, 6)},
		{token.SEMICOLON, ";", makePos(28, 7)},

		{token.LBRACE, "{", makePos(29, 1)},
		{token.STRING, "foo", makePos(29, 2)},
		{token.COLON, ":", makePos(29, 7)},
		{token.STRING, "bar", makePos(29, 9)},
		{token.RBRACE, "}", makePos(29, 14)},
		{token.SEMICOLON, ";", makePos(29, 15)},

		{token.NULL, "nil", makePos(30, 1)},
		{token.SEMICOLON, ";", makePos(30, 4)},

		{token.COMMENT, "Single line hash comment", makePos(31, 1)},
		{token.COMMENT, "Single line slash comment", makePos(32, 1)},
		{token.COMMENT, " Multi-line\ncomment ", makePos(33, 1)},

		{token.FLOAT, "12.5", makePos(36, 1)},
		{token.SEMICOLON, ";", makePos(36, 5)},

		{token.ILLEGAL, "Invalid float literal", makePos(37, 1)},
		{token.ILLEGAL, "Invalid float literal", makePos(37, 5)},
		{token.INT, "7", makePos(37, 6)},
		{token.SEMICOLON, ";", makePos(37, 7)},

		{token.STRING, "\n\r\t\v\f\\\"\b\\' Hello", makePos(39, 1)},
		{token.SEMICOLON, ";", makePos(39, 27)},

		{token.STRING, "\\n\\r\\t\\v\\f\\\\\\\"\\b' Hello", makePos(40, 1)},
		{token.SEMICOLON, ";", makePos(40, 27)},

		{token.INT, "0x2F54a", makePos(42, 1)},
		{token.SEMICOLON, ";", makePos(42, 7)},

		{token.EOF, "", makePos(43, 0)},
	}

	l := NewString(input)
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, %q, got=%q, %q",
				i, tt.expectedType, tt.expectedLiteral, tok.Type, tok.Literal)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}

		if tok.Pos.Line != tt.expectedPos.Line {
			t.Fatalf("tests[%d] - line wrong. Expected %d, got %d <%s,%s>",
				i, tt.expectedPos.Line, tok.Pos.Line, tok.Type, tok.Literal)
		}

		if tok.Pos.Col != tt.expectedPos.Col {
			t.Fatalf("tests[%d] - col wrong. Expected %d, got %d <%s,%s>",
				i, tt.expectedPos.Col, tok.Pos.Col, tok.Type, tok.Literal)
		}
	}
}
