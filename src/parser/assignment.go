package parser

import (
	"fmt"

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
	case token.Interface:
		return p.parseInterfaceDefStatement()
	case token.For:
		return p.parseForLoop()
	case token.While:
		return p.parseWhileLoop()
	case token.Loop:
		return p.parseInfiniteLoop()
	case token.Import:
		return p.parseImport()
	case token.Delete:
		return p.parseDelete()
	case token.Use:
		return p.parseUseStatement()
	case token.Continue:
		stat := &ast.ContinueStatement{
			Token: p.curToken,
		}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	case token.Break:
		stat := &ast.BreakStatement{
			Token: p.curToken,
		}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	case token.Pass:
		stat := &ast.PassStatement{
			Token: p.curToken,
		}
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}
		return stat
	case token.EOF:
		panic("Something messed up big time")
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseUseStatement() ast.Statement {
	stmt := &ast.DefStatement{
		Token: p.curToken,
		Const: true,
	}

	p.nextToken()

	node := p.parseExpression(priLowest)
	exp, ok := node.(*ast.AttributeExpression)
	if !ok {
		p.addErrorWithCurPos("use expected an attribute expression")
		return nil
	}
	stmt.Value = exp

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()

		stmt.Name = &ast.Identifier{Value: ImportName(exp.Index.String())}
		if stmt.Name.Value == "" {
			p.addErrorWithCurPos("use statement does not create a valid identifier")
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

func (p *Parser) parseDelete() ast.Statement {
	stmt := &ast.DeleteStatement{Token: p.curToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt.Name = p.curToken.Literal

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseImport() ast.Statement {
	stmt := &ast.ImportStatement{Token: p.curToken}

	if !p.expectPeek(token.String) {
		return nil
	}

	if p.curToken.Literal == "" {
		p.addErrorWithCurPos("import path cannot be empty")
		return nil
	}

	stmt.Path = &ast.StringLiteral{Token: p.curToken, Value: []rune(p.curToken.Literal)}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()

		stmt.Name = &ast.Identifier{Value: ImportName(stmt.Path.String())}
		if stmt.Name.Value == "" {
			p.addErrorWithCurPos("import path does not create a valid identifier")
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
			p.addErrorWithCurPos("Constant defined with no value")
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
		if fun.Name != "(anonymous)" || fun.FQName != "(anonymous)" {
			p.addErrorWithCurPos("Function definition with let cannot have two names")
		}
		fun.Name = stmt.Name.String()
		fun.FQName = stmt.Name.String()
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

	if p.peekTokenIs(token.Semicolon, token.RBrace) {
		if p.peekTokenIs(token.Semicolon) {
			p.nextToken()
		}

		stmt.Value = &ast.NullLiteral{Token: createKeywordToken("null")}
		return stmt
	}
	p.nextToken()

	exp := p.parseExpression(priLowest)
	if exp == nil {
		stmt.Value = &ast.NullLiteral{Token: createKeywordToken("null")}
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

	startToken := p.curToken
	stmt := &ast.DefStatement{Token: createKeywordToken("let")}

	var ok bool
	stmt.Value, ok = p.parseExpression(priLowest).(ast.Expression)
	if !ok {
		return nil
	}

	fun, ok := stmt.Value.(*ast.FunctionLiteral)
	if !ok {
		p.addErrorWithPos(startToken.Pos, "Expected something else")
		return nil
	}

	if fun.Name == "(anonymous)" {
		p.addErrorWithPos(startToken.Pos, "Anonymous function with no definition or name")
		return nil
	}

	stmt.Name = &ast.Identifier{Token: startToken, Value: fun.Name}

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

	exp, ok := p.parseExpression(priLowest).(ast.Expression)
	if !ok {
		return nil
	}
	stmt.Value = exp

	if class, ok := stmt.Value.(*ast.ClassLiteral); ok {
		class.Name = stmt.Name.String()
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseInterfaceDefStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseInterfaceDefStatement")
	}
	if !p.peekTokenIs(token.Identifier) {
		return p.parseExpressionStatement()
	}

	InterfaceToken := p.curToken
	if !p.expectPeek(token.Identifier) {
		return nil
	}

	stmt := &ast.DefStatement{Token: createKeywordToken("let")}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	p.insertToken(InterfaceToken)
	p.nextToken()

	exp, ok := p.parseExpression(priLowest).(ast.Expression)
	if !ok {
		return nil
	}
	stmt.Value = exp

	if iface, ok := stmt.Value.(*ast.InterfaceLiteral); ok {
		iface.Name = stmt.Name.String()
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

	var ok bool
	stmt.Value, ok = p.parseExpression(priLowest).(ast.Expression)
	if !ok {
		p.addErrorWithPos(p.curToken.Pos, "Invalid assignment")
		return nil
	}

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
