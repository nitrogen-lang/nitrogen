package parser

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		fallthrough
	case token.Always:
		return p.parseDefStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Function:
		return p.parseFuncDefStatement()
	case token.Class:
		return p.parseClassDefStatement()
	case token.For:
		return p.parseForLoop()
	case token.Throw:
		p.nextToken()
		t := &ast.ThrowStatement{
			Expression: p.parseExpression(priLowest),
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
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseDefStatement() ast.Statement {
	stmt := &ast.DefStatement{Token: p.curToken}

	stmt.Const = p.curTokenIs(token.Always)

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

	stmt.Value = p.parseExpression(priLowest)

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
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	if p.curTokenIs(token.Semicolon) {
		stmt.Value = &ast.NullLiteral{Token: createKeywordToken("null")}
		return stmt
	}

	stmt.Value = p.parseExpression(priLowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFuncDefStatement() ast.Statement {
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

	stmt.Value = p.parseExpression(priLowest)

	if fun, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fun.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseClassDefStatement() ast.Statement {
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

	stmt.Value = p.parseExpression(priLowest)

	if class, ok := stmt.Value.(*ast.ClassLiteral); ok {
		class.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignmentStatement(left ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	stmt.Value = p.parseExpression(priLowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

// parseCompoundAssign will take a compound assignment like += and
// turn it into a normal assignment statement using the given left
// expression as the left side of a normal arithmatic operation
func (p *Parser) parseCompoundAssign(left ast.Expression) ast.Expression {
	stmt := &ast.AssignStatement{
		Token: p.curToken,
		Left:  left,
	}

	p.nextToken()

	right := p.parseExpression(priLowest)

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
