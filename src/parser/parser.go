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
	token.BitwiseOr:     priSum,
	token.Carrot:        priSum,
	token.Slash:         priProduct,
	token.Asterisk:      priProduct,
	token.Modulo:        priProduct,
	token.ShiftLeft:     priProduct,
	token.ShiftRight:    priProduct,
	token.BitwiseAnd:    priProduct,
	token.BitwiseAndNot: priProduct,
	token.LParen:        priCall,
	token.Implements:    priCall,
	token.LSquare:       priIndex,
	token.Dot:           priIndex,
	token.Assign:        priAssign,
	token.PlusAssign:    priAssign,
	token.MinusAssign:   priAssign,
	token.TimesAssign:   priAssign,
	token.SlashAssign:   priAssign,
	token.ModAssign:     priAssign,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Node
	Settings      struct {
		Debug bool
	}
)

type Parser struct {
	l        *lexer.Lexer
	errors   []string
	settings *Settings

	lastToken token.Token
	curToken  token.Token
	peekToken token.Token

	insertedTokens []token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer, settings *Settings) *Parser {
	if settings == nil {
		settings = &Settings{}
	}

	p := &Parser{
		l:              l,
		settings:       settings,
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
	p.registerPrefix(token.Class, p.parseClassLiteral)
	p.registerPrefix(token.Interface, p.parseInterfaceLiteral)
	p.registerPrefix(token.New, p.parseMakeExpression)
	p.registerPrefix(token.Do, p.parseDoExpression)
	p.registerPrefix(token.Recover, p.parseRecoverExpression)

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
	p.registerInfix(token.Dot, p.parseAttributeExpression)
	p.registerInfix(token.Assign, p.parseAssignmentStatement)
	p.registerInfix(token.PlusAssign, p.parseCompoundAssign)
	p.registerInfix(token.MinusAssign, p.parseCompoundAssign)
	p.registerInfix(token.TimesAssign, p.parseCompoundAssign)
	p.registerInfix(token.SlashAssign, p.parseCompoundAssign)
	p.registerInfix(token.ModAssign, p.parseCompoundAssign)
	p.registerInfix(token.ShiftLeft, p.parseInfixExpression)
	p.registerInfix(token.ShiftRight, p.parseInfixExpression)
	p.registerInfix(token.BitwiseAnd, p.parseInfixExpression)
	p.registerInfix(token.BitwiseAndNot, p.parseInfixExpression)
	p.registerInfix(token.BitwiseOr, p.parseInfixExpression)
	p.registerInfix(token.Carrot, p.parseInfixExpression)
	p.registerInfix(token.Implements, p.parseInfixExpression)

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

func (p *Parser) addErrorWithCurPos(format string, args ...interface{}) {
	p.addErrorWithPos(p.curToken.Pos, format, args...)
}

func (p *Parser) addErrorWithPos(pos token.Position, format string, args ...interface{}) {
	if len(args) > 0 {
		args = append([]interface{}{pos.Filename, pos.Line, pos.Col}, args...)
	} else {
		args = []interface{}{pos.Filename, pos.Line, pos.Col}
	}
	msg := fmt.Sprintf("%s:\n  line %d, col %d:\n    "+format, args...)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(t token.TokenType) {
	p.addErrorWithPos(p.peekToken.Pos, "Incorrect next token. Expected %q, got %q", t.String(),
		p.peekToken.Type.String())
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	if t == token.BitwiseOr {
		p.addErrorWithCurPos("Invalid prefix: %q. Did you mean \"or\"?", t.String())
	} else if t == token.BitwiseAnd {
		p.addErrorWithCurPos("Invalid prefix: %q. Did you mean \"and\"?", t.String())
	} else {
		p.addErrorWithCurPos("Invalid prefix: %q", t.String())
	}
}

func (p *Parser) nextToken() {
	p.advanceToken()
	for p.peekTokenIs(token.Comment) {
		p.advancePeekToken()
	}

	if p.curTokenIs(token.Illegal) {
		p.addErrorWithCurPos("Got illegal token: %s", p.curToken.Literal)
		p.nextToken()
	}

	if p.settings.Debug {
		fmt.Println(p.curToken)
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
	if p.settings.Debug {
		fmt.Println("ParseProgram")
	}
	program := &ast.Program{Filename: p.curToken.Pos.Filename}
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

func (p *Parser) parseExpression(precedence int) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseExpression")
	}

	var prefix prefixParseFn
	skips := 0

	// Attempt to parse a valid expression by skipping up to 2 tokens
	// if the current token isn't a valid prefix.
	for skips = 0; skips < 2; skips += 1 {
		prefix = p.prefixParseFns[p.curToken.Type]

		if prefix != nil {
			break
		}

		p.noPrefixParseFnError(p.curToken.Type)
		p.nextToken()
	}

	if skips == 2 {
		return nil
	}

	var leftExp ast.Node
	leftExp = prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		currPos := p.curToken.Pos

		p.nextToken()
		lexp, ok := leftExp.(ast.Expression)
		if !ok {
			p.addErrorWithPos(currPos, "Failed parsing expression")
			return nil
		}
		leftExp = infix(lexp)
	}

	return leftExp
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	if p.settings.Debug {
		fmt.Println("parseExpressionStatement")
	}
	var stmt ast.Statement
	curToken := p.curToken

	exp := p.parseExpression(priLowest)
	if exp == nil {
		return nil
	}

	if expStmt, ok := exp.(ast.Statement); ok {
		stmt = expStmt
	} else {
		stmt = &ast.ExpressionStatement{
			Token:      curToken,
			Expression: exp.(ast.Expression),
		}
	}

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseBlockStatements() *ast.BlockStatement {
	if p.settings.Debug {
		fmt.Println("parseBlockStatements")
	}
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}
	p.nextToken()

	for !p.curTokenIs(token.RBrace) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseSingleStmtBlock() *ast.BlockStatement {
	if p.settings.Debug {
		fmt.Println("parseSingleStmtBlock")
	}
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: make([]ast.Statement, 1),
	}
	p.nextToken()

	stmt := p.parseStatement()
	if stmt != nil {
		block.Statements[0] = stmt
	}

	return block
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseExpressionList")
	}
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	expr := p.parseExpression(priLowest)
	if expr == nil {
		return nil
	}
	list = append(list, expr.(ast.Expression))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		if p.curTokenIs(end) {
			return list
		}
		expr := p.parseExpression(priLowest)
		if expr == nil {
			return nil
		}
		list = append(list, expr.(ast.Expression))
	}

	if !p.peekTokenIs(end) {
		p.addErrorWithCurPos("I was expecting %q, but instead I got %q. Did you forget a comma?", end.String(), p.peekToken.Literal)
		return nil
	}

	p.nextToken()
	return list
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parsePrefixExpression")
	}
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(priPrefix).(ast.Expression)
	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Node {
	if p.settings.Debug {
		fmt.Println("parseInfixExpression")
	}
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	var ok bool
	expression.Right, ok = p.parseExpression(precedence).(ast.Expression)
	if !ok {
		return nil
	}
	return expression
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	return p.parseGroupedExpressionC(token.RParen)
}

func (p *Parser) parseGroupedExpressionC(end token.TokenType) ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseGroupedExpressionC")
	}
	p.nextToken()
	exp := p.parseExpression(priLowest).(ast.Expression)

	if !p.expectPeek(end) {
		return nil
	}
	return exp
}

func (p *Parser) parseGroupedExpressionE() ast.Expression {
	if p.settings.Debug {
		fmt.Println("parseGroupedExpressionE")
	}
	p.nextToken()
	exp := p.parseExpression(priLowest)
	if exp == nil {
		return nil
	}

	return exp.(ast.Expression)
}
