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

"\n\r\t\v\f\\\"\b\' Hello"
'\n\r\t\v\f\\\"\b\' Hello'

\x2F54a

or

and

for

->
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
		expectedPos     token.Position
	}{
		{token.Let, "let", makePos(1, 1)},
		{token.Identifier, "five", makePos(1, 5)},
		{token.Assign, "=", makePos(1, 10)},
		{token.Integer, "5", makePos(1, 12)},
		{token.Semicolon, ";", makePos(1, 13)},

		{token.Always, "always", makePos(2, 1)},
		{token.Identifier, "ten", makePos(2, 8)},
		{token.Assign, "=", makePos(2, 12)},
		{token.Integer, "10", makePos(2, 14)},
		{token.Semicolon, ";", makePos(2, 16)},

		{token.Identifier, "世界", makePos(4, 1)},
		{token.Semicolon, ";", makePos(4, 3)},

		{token.Identifier, "ident2", makePos(6, 1)},
		{token.Semicolon, ";", makePos(6, 7)},

		{token.Integer, "6", makePos(8, 1)},
		{token.Modulo, "%", makePos(8, 3)},
		{token.Integer, "3", makePos(8, 5)},
		{token.Semicolon, ";", makePos(8, 6)},

		{token.Let, "let", makePos(10, 1)},
		{token.Identifier, "add", makePos(10, 5)},
		{token.Assign, "=", makePos(10, 9)},
		{token.Function, "func", makePos(10, 11)},
		{token.LParen, "(", makePos(10, 15)},
		{token.Identifier, "x", makePos(10, 16)},
		{token.Comma, ",", makePos(10, 17)},
		{token.Identifier, "y", makePos(10, 19)},
		{token.RParen, ")", makePos(10, 20)},
		{token.LBrace, "{", makePos(10, 22)},
		{token.Identifier, "x", makePos(11, 2)},
		{token.Plus, "+", makePos(11, 4)},
		{token.Identifier, "y", makePos(11, 6)},
		{token.Semicolon, ";", makePos(11, 7)},
		{token.RBrace, "}", makePos(12, 1)},
		{token.Semicolon, ";", makePos(12, 2)},

		{token.Let, "let", makePos(14, 1)},
		{token.Identifier, "result", makePos(14, 5)},
		{token.Assign, "=", makePos(14, 12)},
		{token.Identifier, "add", makePos(14, 14)},
		{token.LParen, "(", makePos(14, 17)},
		{token.Identifier, "five", makePos(14, 18)},
		{token.Comma, ",", makePos(14, 22)},
		{token.Identifier, "ten", makePos(14, 24)},
		{token.RParen, ")", makePos(14, 27)},
		{token.Semicolon, ";", makePos(14, 28)},

		{token.Bang, "!", makePos(15, 1)},
		{token.Dash, "-", makePos(15, 2)},
		{token.Slash, "/", makePos(15, 3)},
		{token.Integer, "5", makePos(15, 4)},
		{token.Asterisk, "*", makePos(15, 5)},
		{token.Semicolon, ";", makePos(15, 6)},

		{token.Integer, "6", makePos(16, 1)},
		{token.LessThan, "<", makePos(16, 3)},
		{token.Integer, "10", makePos(16, 5)},
		{token.GreaterThan, ">", makePos(16, 8)},
		{token.Integer, "7", makePos(16, 10)},
		{token.Semicolon, ";", makePos(16, 11)},

		{token.If, "if", makePos(18, 1)},
		{token.LParen, "(", makePos(18, 4)},
		{token.Integer, "5", makePos(18, 5)},
		{token.LessThan, "<", makePos(18, 7)},
		{token.Integer, "10", makePos(18, 9)},
		{token.RParen, ")", makePos(18, 11)},
		{token.LBrace, "{", makePos(18, 13)},
		{token.Return, "return", makePos(19, 2)},
		{token.True, "true", makePos(19, 9)},
		{token.Semicolon, ";", makePos(19, 13)},
		{token.RBrace, "}", makePos(20, 1)},
		{token.Else, "else", makePos(20, 3)},
		{token.LBrace, "{", makePos(20, 8)},
		{token.Return, "return", makePos(21, 2)},
		{token.False, "false", makePos(21, 9)},
		{token.Semicolon, ";", makePos(21, 14)},
		{token.RBrace, "}", makePos(22, 1)},
		{token.Semicolon, ";", makePos(22, 2)},

		{token.Integer, "10", makePos(24, 1)},
		{token.Equal, "==", makePos(24, 4)},
		{token.Integer, "10", makePos(24, 7)},
		{token.Semicolon, ";", makePos(24, 9)},

		{token.Integer, "10", makePos(25, 1)},
		{token.NotEqual, "!=", makePos(25, 4)},
		{token.Integer, "9", makePos(25, 7)},
		{token.Semicolon, ";", makePos(25, 8)},

		{token.String, "foobar", makePos(26, 1)},
		{token.Semicolon, ";", makePos(26, 9)},

		{token.String, "foo bar", makePos(27, 1)},
		{token.Semicolon, ";", makePos(27, 10)},

		{token.LSquare, "[", makePos(28, 1)},
		{token.Integer, "1", makePos(28, 2)},
		{token.Comma, ",", makePos(28, 3)},
		{token.Integer, "2", makePos(28, 5)},
		{token.RSquare, "]", makePos(28, 6)},
		{token.Semicolon, ";", makePos(28, 7)},

		{token.LBrace, "{", makePos(29, 1)},
		{token.String, "foo", makePos(29, 2)},
		{token.Colon, ":", makePos(29, 7)},
		{token.String, "bar", makePos(29, 9)},
		{token.RBrace, "}", makePos(29, 14)},
		{token.Semicolon, ";", makePos(29, 15)},

		{token.Nil, "nil", makePos(30, 1)},
		{token.Semicolon, ";", makePos(30, 4)},

		{token.Comment, "Single line hash comment", makePos(31, 1)},
		{token.Comment, "Single line slash comment", makePos(32, 1)},
		{token.Comment, " Multi-line\ncomment ", makePos(33, 1)},

		{token.Float, "12.5", makePos(36, 1)},
		{token.Semicolon, ";", makePos(36, 5)},

		{token.String, "\n\r\t\v\f\\\"\b\\' Hello", makePos(38, 1)},
		{token.Semicolon, ";", makePos(38, 27)},

		{token.String, "\\n\\r\\t\\v\\f\\\\\\\"\\b' Hello", makePos(39, 1)},
		{token.Semicolon, ";", makePos(39, 27)},

		{token.Integer, "0x2F54a", makePos(41, 1)},
		{token.Semicolon, ";", makePos(41, 9)},

		{token.LOr, "or", makePos(42, 1)},

		{token.LAnd, "and", makePos(44, 1)},

		{token.For, "for", makePos(46, 1)},

		{token.Arrow, "->", makePos(48, 1)},

		{token.EOF, "", makePos(49, 0)},
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
