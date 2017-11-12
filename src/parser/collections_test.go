package parser

import (
	"testing"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
)

func TestParsingArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	array, ok := stmt.Expression.(*ast.Array)
	if !ok {
		t.Fatalf("exp not ast.Array. got=%T", stmt.Expression)
	}

	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}

	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"
	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestParsingHashLiteralsStringKeys(t *testing.T) {
	input := `{"one": 1, "two": 2, "three": 3}`

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestParsingEmptyHashLiteral(t *testing.T) {
	input := "{}"

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}
}

func TestParsingHashLiteralsWithExpressions(t *testing.T) {
	input := `{"one": 0 + 1, "two": 10 - 8, "three": 15 / 5}`

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	tests := map[string]func(ast.Expression){
		"one": func(e ast.Expression) {
			testInfixExpression(t, e, 0, "+", 1)
		},
		"two": func(e ast.Expression) {
			testInfixExpression(t, e, 10, "-", 8)
		},
		"three": func(e ast.Expression) {
			testInfixExpression(t, e, 15, "/", 5)
		},
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
			continue
		}
		testFunc, ok := tests[literal.String()]
		if !ok {
			t.Errorf("No test function for key %q found", literal.String())
			continue
		}
		testFunc(value)
	}
}

func TestHashAssignments(t *testing.T) {
	input := `m["key"] = "value";`

	l := lexer.NewString(input)
	p := New(l)
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

	left, ok := exp.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("left is not ast.IndexExpression. got=%T",
			exp.Left)
	}

	if left.Left.String() != "m" {
		t.Fatalf("left hash ident is not correct. expected=m, got=%s",
			left.Left.String())
	}

	if left.Index.String() != "key" {
		t.Fatalf("left index is not correct. expected=\"key\", got=%q",
			left.Index.String())
	}

	if exp.Value.String() != "value" {
		t.Fatalf("assigment value is not correct. expected=\"value\", got=%q",
			exp.Value.String())
	}
}

func TestArrayAssignments(t *testing.T) {
	input := `m[0] = "value";`

	l := lexer.NewString(input)
	p := New(l)
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

	left, ok := exp.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("left is not ast.IndexExpression. got=%T",
			exp.Left)
	}

	if left.Left.String() != "m" {
		t.Fatalf("left hash ident is not correct. expected=m, got=%s",
			left.Left.String())
	}

	if left.Index.String() != "0" {
		t.Fatalf("left index is not correct. expected=\"0\", got=%q",
			left.Index.String())
	}

	if exp.Value.String() != "value" {
		t.Fatalf("assigment value is not correct. expected=\"value\", got=%q",
			exp.Value.String())
	}
}

func TestHashArrowOperator(t *testing.T) {
	input := `m->value;`

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp is not ast.IndexExpression. got=%T",
			stmt.Expression)
	}

	if exp.Left.String() != "m" {
		t.Fatalf("exp hash ident is not correct. expected=m, got=%s",
			exp.Left.String())
	}

	if exp.Index.String() != "value" {
		t.Fatalf("exp index is not correct. expected=value, got=%q",
			exp.Index.String())
	}
}

func TestHashArrowOperatorString(t *testing.T) {
	input := `m->"value";`

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp is not ast.IndexExpression. got=%T",
			stmt.Expression)
	}

	if exp.Left.String() != "m" {
		t.Fatalf("exp hash ident is not correct. expected=m, got=%s",
			exp.Left.String())
	}

	if exp.Index.String() != "value" {
		t.Fatalf("exp index is not correct. expected=value, got=%q",
			exp.Index.String())
	}
}

func TestParsingHashLiteralsMultiLine(t *testing.T) {
	input := `{
                "one": 1,
                "two": 2,
                "three": 3,
              }`

	l := lexer.NewString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	hash, ok := stmt.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("exp is not ast.HashLiteral. got=%T", stmt.Expression)
	}

	if len(hash.Pairs) != 3 {
		t.Errorf("hash.Pairs has wrong length. got=%d", len(hash.Pairs))
	}

	expected := map[string]int64{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	for key, value := range hash.Pairs {
		literal, ok := key.(*ast.StringLiteral)
		if !ok {
			t.Errorf("key is not ast.StringLiteral. got=%T", key)
		}
		expectedValue := expected[literal.String()]
		testIntegerLiteral(t, value, expectedValue)
	}
}

func TestInvalidHashLiteralMultiLine(t *testing.T) {
	input := `{
                "one": 1,
                "two": 2,
                "three": 3
              }`

	l := lexer.NewString(input)
	p := New(l)
	p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Fatalf("No errors for invalid hash literal, missing comma")
	}

	if p.Errors()[0] != "<4,26> Hash pairs must end with a comma" {
		t.Fatalf("Incorrect error message. Got %q", p.Errors()[0])
	}
}
