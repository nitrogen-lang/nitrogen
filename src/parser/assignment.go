package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseStatment")
	}
	switch p.curToken.Type {
	case token.Let:
		fallthrough
	case token.Const:
		return p.parseDefStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Function:
		return p.parseFuncDefStatement()
	case token.Class:
		return p.parseClassDefStatement()
	case token.For:
		return p.parseForLoop()
	case token.Import:
		return p.parseImport()
	case token.Throw:
		p.nextToken()
		t := &ast.ThrowStatement{
			Expression: p.parseExpression(priLowest).(ast.Expression),
		}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return t
	case token.Continue:
		stat := &ast.ContinueStatement{}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	case token.Break:
		stat := &ast.BreakStatement{}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	case token.Pass:
		stat := &ast.PassStatement{}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseImport() ast.Statement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(token.String) {
		return nil
	}

	if p.curToken.Literal == "" {
		p.addErrorWithPos("import path cannot be empty")
		return nil
	}

	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()

		stmt.Name = &ast.Identifier{Value: importName(stmt.Path.Value)}
		if stmt.Name.Value == "" {
			p.addErrorWithPos("import path does not create a valid identifier")
			return nil
		}
		return stmt
	}

	if !p.expectPeek(token.As) {
		return nil
	}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return stmt
}

func importName(path string) string {
	path = filepath.Base(path)
	path = path[:strings.LastIndex(path, ".")]
	if isIdent(path) {
		return path
	}
	return ""
}

func (p *Parser) parseDefStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseDefStatement")
	}
	stmt := &ast.DefStatement{Token: p.curToken}

	stmt.Const = p.curTokenIs(token.Const)

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()

		if stmt.Const {
			p.addErrorWithPos("Constant defined with no value")
			return nil
		}
		return stmt
	}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	p.nextToken()

	thing := p.parseExpression(priLowest)
	if thing == nil {
		return nil
	}
	stmt.Value = thing.(ast.Expression)

	if fun, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fun.Name = stmt.Name.String()
	}
	if class, ok := stmt.Value.(*ast.ClassLiteral); ok {
		class.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseReturnStatement")
	}
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	if p.curTokenIs(token.Semicolon) {
		stmt.Value = &ast.NullLiteral{Token: createKeywordToken("null")}
		return stmt
	}

	exp := p.parseExpression(priLowest)
	if exp == nil {
		stmt.Value = nil
	} else {
		stmt.Value = exp.(ast.Expression)
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFuncDefStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseFuncDefStatement")
	}
	if !p.peekTokenIs(token.Identifier) {
		return p.parseExpressionStatement()
	}

	fToken := p.curToken
	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt := &ast.DefStatement{Token: createKeywordToken("let")}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.insertToken(fToken)
	p.nextToken()

	stmt.Value = p.parseExpression(priLowest).(ast.Expression)

	if fun, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fun.Name = stmt.Name.String()
		fun.FQName = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseClassDefStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseClassDefStatement")
	}
	if !p.peekTokenIs(token.Identifier) {
		return p.parseExpressionStatement()
	}

	classToken := p.curToken
	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt := &ast.DefStatement{Token: createKeywordToken("let")}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.insertToken(classToken)
	p.nextToken()

	stmt.Value = p.parseExpression(priLowest).(ast.Expression)

	if class, ok := stmt.Value.(*ast.ClassLiteral); ok {
		class.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseAssignmentStatement")
	}
	stmt := &ast.AssignStatement{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	stmt.Value = p.parseExpression(priLowest).(ast.Expression)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseCompoundAssign will take a compound assignment like += and
// turn it into a normal assignment statement using the given left
// expression as the left side of a normal arithmatic operation
func (p *Parser) parseCompoundAssign(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseCompoundAssign")
	}
	stmt := &ast.AssignStatement{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	right := p.parseExpression(priLowest).(ast.Expression)

	switch stmt.Token.Type {
	case token.PlusAssign:
		stmt.Value = makeInfix(token.Plus, left, right)
	case token.MinusAssign:
		stmt.Value = makeInfix(token.Dash, left, right)
	case token.TimesAssign:
		stmt.Value = makeInfix(token.Asterisk, left, right)
	case token.SlashAssign:
		stmt.Value = makeInfix(token.Slash, left, right)
	case token.ModAssign:
		stmt.Value = makeInfix(token.Modulo, left, right)
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func makeInfix(tokenType token.TokenType, left, right ast.Expression) *ast.InfixExpression {
	return &ast.InfixExpression{
		Token: token.Token{
			Type:    tokenType,
			Literal: tokenType.String(),
		},
		Left:     left,
		Operator: tokenType.String(),
		Right:    right,
	}
}
