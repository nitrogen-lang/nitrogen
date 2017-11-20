package parser

import (
	"fmt"

	"github.com/nitrogen-lang/nitrogen/src/ast"
	"github.com/nitrogen-lang/nitrogen/src/lexer"
	"github.com/nitrogen-lang/nitrogen/src/token"
)

const (
	priLowest      int = iota
	priCompare         // and, or
	priEquals          // ==
	priLessGreater     // > or <
	priSum             // +, -
	priProduct         // *, /
	priPrefix          // -x or !x
	priCall            // myFunction(x)
	priIndex           // array[index]
	priAssign
)

var precedences = map[token.TokenType]int{
	token.LAnd:          priCompare,
	token.LOr:           priCompare,
	token.Equal:         priEquals,
	token.NotEqual:      priEquals,
	token.LessThanEq:    priEquals,
	token.GreaterThanEq: priEquals,
	token.LessThan:      priLessGreater,
	token.GreaterThan:   priLessGreater,
	token.Plus:          priSum,
	token.Dash:          priSum,
	token.Slash:         priProduct,
	token.Asterisk:      priProduct,
	token.Modulo:        priProduct,
	token.LParen:        priCall,
	token.LSquare:       priIndex,
	token.Arrow:         priIndex,
	token.Assign:        priAssign,
	token.PlusAssign:    priAssign,
	token.MinusAssign:   priAssign,
	token.TimesAssign:   priAssign,
	token.SlashAssign:   priAssign,
	token.ModAssign:     priAssign,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	lastToken token.Token
	curToken  token.Token
	peekToken token.Token

	insertedTokens []token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		insertedTokens: make([]token.Token, 0, 5),
	}

	// Prefix parsing functions
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Identifier, p.parseIdentifier)
	p.registerPrefix(token.Integer, p.parseIntegerLiteral)
	p.registerPrefix(token.Float, p.parseFloatLiteral)
	p.registerPrefix(token.Nil, p.parseNullLiteral)
	p.registerPrefix(token.String, p.parseStringLiteral)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LSquare, p.parseArrayLiteral)
	p.registerPrefix(token.LBrace, p.parseHashLiteral)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Dash, p.parsePrefixExpression)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Try, p.parseTryCatch)

	// Infix parsing functions
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Dash, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Modulo, p.parseInfixExpression)
	p.registerInfix(token.Equal, p.parseInfixExpression)
	p.registerInfix(token.NotEqual, p.parseInfixExpression)
	p.registerInfix(token.LessThanEq, p.parseInfixExpression)
	p.registerInfix(token.GreaterThanEq, p.parseInfixExpression)
	p.registerInfix(token.LessThan, p.parseInfixExpression)
	p.registerInfix(token.GreaterThan, p.parseInfixExpression)
	p.registerInfix(token.LAnd, p.parseCompareExpression)
	p.registerInfix(token.LOr, p.parseCompareExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)
	p.registerInfix(token.LSquare, p.parseIndexExpression)
	p.registerInfix(token.Arrow, p.parseIndexExpression)
	p.registerInfix(token.Assign, p.parseAssignmentStatement)
	p.registerInfix(token.PlusAssign, p.parseCompoundAssign)
	p.registerInfix(token.MinusAssign, p.parseCompoundAssign)
	p.registerInfix(token.TimesAssign, p.parseCompoundAssign)
	p.registerInfix(token.SlashAssign, p.parseCompoundAssign)
	p.registerInfix(token.ModAssign, p.parseCompoundAssign)

	// Read the first two tokens to populate curToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefix(tt token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tt] = fn
}

func (p *Parser) registerInfix(tt token.TokenType, fn infixParseFn) {
	p.infixParseFns[tt] = fn
}

func (p *Parser) addError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	p.errors = append(p.errors, msg)
}

func (p *Parser) addErrorWithPos(format string, args ...interface{}) {
	if len(args) > 0 {
		args = append([]interface{}{p.curToken.Pos.Line, p.curToken.Pos.Col}, args...)
	} else {
		args = []interface{}{p.curToken.Pos.Line, p.curToken.Pos.Col}
	}
	msg := fmt.Sprintf("at line %d, col %d "+format, args...)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	p.addError(
		"at line %d, col %d Incorrect next token. Expected %q, got %q",
		p.peekToken.Pos.Line,
		p.peekToken.Pos.Col,
		t.String(),
		p.peekToken.Type.String(),
	)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	p.addErrorWithPos("Invalid prefix: %s", t)
}

func (p *Parser) nextToken() {
	p.advanceToken()
	for p.peekTokenIs(token.Comment) {
		p.advancePeekToken()
	}

	if p.curTokenIs(token.Illegal) {
		p.addErrorWithPos("Got illegal token: %s", p.curToken.Literal)
		p.nextToken()
	}
}

func (p *Parser) advanceToken() {
	p.lastToken = p.curToken
	p.curToken = p.peekToken
	p.advancePeekToken()
}

func (p *Parser) advancePeekToken() {
	if len(p.insertedTokens) > 0 {
		p.peekToken = p.insertedTokens[len(p.insertedTokens)-1]
		p.insertedTokens = p.insertedTokens[:len(p.insertedTokens)-1]
	} else {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) insertToken(t token.Token) {
	p.insertedTokens = append(p.insertedTokens, p.peekToken)
	p.peekToken = t
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Filename: p.curToken.Filename}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		if len(p.errors) > 0 {
			return program
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(priLowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatements() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}
	p.nextToken()

	for !p.curTokenIs(token.RBrace) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(priLowest))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(priLowest))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(priPrefix)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	return p.parseGroupedExpressionC(token.RParen)
}

func (p *Parser) parseGroupedExpressionC(end token.TokenType) ast.Expression {
	p.nextToken()
	exp := p.parseExpression(priLowest)

	if !p.expectPeek(end) {
		return nil
	}
	return exp
}
