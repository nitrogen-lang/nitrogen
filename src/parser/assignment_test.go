package parser

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
)

func TestDefStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewString(tt.input)
		p := New(l, nil)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testDefStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.DefStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"always x = 5;", "x", 5},
		{"always y = true;", "y", true},
		{"always foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewString(tt.input)
		p := New(l, nil)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testConstStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.DefStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return true;", true},
		{"return foobar;", "foobar"},
		{"return;", nil},
	}

	for _, tt := range tests {
		l := lexer.NewString(tt.input)
		p := New(l, nil)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.returnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}
		if !testLiteralExpression(t, returnStmt.Value, tt.expectedValue) {
			return
		}
	}
}

func TestFunctionSugar(t *testing.T) {
	input := `func hello(place) {
        return "Hello, " + place;
    }`

	l := lexer.NewString(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
			len(program.Statements))
	}

	stmt := program.Statements[0]
	if !testDefStatement(t, stmt, "hello") {
		return
	}

	val := stmt.(*ast.DefStatement).Value
	if _, ok := val.(*ast.FunctionLiteral); !ok {
		t.Fatalf("func sugar invalid, no function literal. got=%T", val)
	}
}

func TestLetFuncSugarStatement(t *testing.T) {
	input := `let hello = func hello_(place) {
        return "Hello, " + place;
    }`

	l := lexer.NewString(input)
	p := New(l, nil)
	p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Fatalf("let with func sugar expected to fail, but didn't")
	}
	if p.Errors()[0] != "at line 1, col 18 Incorrect next token. Expected \"(\", got \"IDENT\"" {
		t.Fatalf("Incorrect error. got \"%s\"", p.Errors()[0])
	}
}

func TestNullReturn(t *testing.T) {
	input := `func hello(place) {
        return
    }`

	l := lexer.NewString(input)
	p := New(l, nil)
	program := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("null return returned error: %s", p.Errors()[0])
	}

	astString := program.String()
	if astString != "func hello(place) {return null;}" {
		t.Fatalf("Incorrect null return parsing. Got %q", astString)
	}
}

func TestGeneralAssignments(t *testing.T) {
	input := `variable = "value";`

	l := lexer.NewString(input)
	p := New(l, nil)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.AssignStatement)
	if !ok {
		t.Fatalf("exp is not ast.AssignStatement. got=%T",
			stmt.Expression)
	}

	ident, ok := exp.Left.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp is not ast.Identifier. got=%T",
			stmt.Expression)
	}

	if ident.Value != "variable" {
		t.Fatalf("ident is not correct. expected=\"variable\", got=%s",
			ident.Value)
	}
}
