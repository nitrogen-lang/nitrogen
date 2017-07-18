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
	case token.For:
		return p.parseForLoop()
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

	if !p.expectPeek(token.Assign) {
		return nil
	}

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
