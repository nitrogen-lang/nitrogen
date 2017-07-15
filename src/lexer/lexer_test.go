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
		};

		let result = add(five, ten);
		!-/5*;
		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} else {
			return false;
		}

		10 == 10;
		10 != 9;
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
    `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.DEF, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.CONST, "always"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		// Test UTF-8 processing
		{token.IDENT, "世界"},

		// Test numbers in ident name
		{token.IDENT, "ident2"},

		// Modulo
		{token.INT, "6"},
		{token.MODULUS, "%"},
		{token.INT, "3"},

		{token.DEF, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.DEF, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.INT, "5"},
		{token.ASTERISK, "*"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},

		{token.INT, "10"},
		{token.EQ, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.LSQUARE, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RSQUARE, "]"},
		{token.SEMICOLON, ";"},

		{token.LBRACE, "{"},
		{token.STRING, "foo"},
		{token.COLON, ":"},
		{token.STRING, "bar"},
		{token.RBRACE, "}"},

		{token.NULL, "nil"},

		{token.COMMENT, "Single line hash comment"},
		{token.COMMENT, "Single line slash comment"},
		{token.COMMENT, " Multi-line\ncomment "},

		{token.FLOAT, "12.5"},
		// Ensure bad floats are lexed, parser handles errors
		{token.FLOAT, "12.5.7"},

		{token.STRING, "\n\r\t\v\f\\\"\b\\' Hello"},
		{token.STRING, `\n\r\t\v\f\\\"\b' Hello`},

		{token.EOF, ""},
	}

	l := NewString(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. Expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. Expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
