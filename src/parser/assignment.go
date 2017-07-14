package parser

import (
	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/token"
) 

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.DEF:
		fallthrough
	case token.CONST:
		return p.parseDefStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.FUNCTION:
		return p.parseFuncDefStatement()
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseDefStatement() ast.Statement {
	stmt := &ast.DefStatement{Token: p.curToken}

	stmt.Const = p.curTokenIs(token.CONST)

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fun, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fun.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	if p.curTokenIs(token.SEMICOLON) {
		stmt.Value = &ast.NullLiteral{Token: createKeywordToken("null")}
		return stmt
	}

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseFuncDefStatement() ast.Statement {
	if !p.peekTokenIs(token.IDENT) {
		return p.parseExpressionStatement()
	}

	fToken := p.curToken
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt := &ast.DefStatement{Token: createKeywordToken("let")}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.insertToken(fToken)
	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if fun, ok := stmt.Value.(*ast.FunctionLiteral); ok {
		fun.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.SEMICOLON) {
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

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
