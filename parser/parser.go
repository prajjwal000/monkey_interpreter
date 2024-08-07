package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn

	errors []string
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

func (p *Parser) peekPrecedence() int {
	if pr, ok := precedences[p.peekToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if pr, ok := precedences[p.curToken.Type]; ok {
		return pr
	}
	return LOWEST
}

func (p *Parser) registerPrefix(to token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[to] = fn
}
func (p *Parser) registerInfix(to token.TokenType, fn infixParseFn) {
	p.infixParseFns[to] = fn
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{}}

	p.nextToken()
	p.nextToken()

	//Expression Parsing
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNC, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
    p.registerPrefix(token.LCURLY, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	//     defer untrace(trace("ParseStatement"))
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	//     defer untrace(trace("ParseLetStatement"))

	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Bind = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	//     defer untrace(trace("ParseReturnStatement"))
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	// 		defer untrace(trace("parseExpressionStatement"))

	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 		defer untrace(trace("parseExpression"))

	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// Expression parsers
func (p *Parser) parseIdentifier() ast.Expression {
	// 		defer untrace(trace("parseIdentifier"))
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}
func (p *Parser) parseBoolean() ast.Expression {
	// 		defer untrace(trace("parseBoolean"))
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}
func (p *Parser) parseIntegerLiteral() ast.Expression {
	// 		defer untrace(trace("parseIntegerLiteral"))
	lit := &ast.IntegerLiteral{Token: p.curToken}

	v, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as int", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = v

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	//     defer untrace(trace("ParseStringLiteral"))
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// 		defer untrace(trace("parsePrefixExpression"))

	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.On = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// 		defer untrace(trace("parseInfixExpression"))

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
	// 		defer untrace(trace("parseGroupedExpression"))
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	// 		defer untrace(trace("parseIfExpression"))
	exp := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LCURLY) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LCURLY) {
			return nil
		}
		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	// 		defer untrace(trace("parseBlockExpression"))
	p.nextToken()

	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	for !p.curTokenIs(token.RCURLY) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	//     defer untrace(trace("ParseFunctionLiteral"))
	f := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	f.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LCURLY) {
		return nil
	}
	f.Body = p.parseBlockStatement()

	return f
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	//     defer untrace(trace("ParseFunctionParameters"))
	parameters := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return parameters
	}

	p.nextToken()
	i := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	parameters = append(parameters, i)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		i := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		parameters = append(parameters, i)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return parameters
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	//     defer untrace(trace("ParseCallExpression"))
	exp := &ast.CallExpression{Token: p.curToken, Function: f}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	//     defer untrace(trace("ParseExpressionList"))
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return list
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	//     defer untrace(trace("ParseArrayLiteral"))
	arr := &ast.ArrayLiteral{Token: p.curToken}

	arr.Elements = p.parseExpressionList(token.RBRACKET)

	return arr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	//     defer untrace(trace("ParseIndexExpression"))
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()

	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

func (p *Parser) parseHashLiteral() ast.Expression {
    hash := &ast.HashLiteral{Token: p.curToken}
    hash.Pairs = make(map[ast.Expression]ast.Expression)

    for !p.peekTokenIs(token.RCURLY) {
        p.nextToken()

        key := p.parseExpression(LOWEST)

        if !p.expectPeek(token.COLON) {
            return nil
        }

        p.nextToken()

        value := p.parseExpression(LOWEST)
        hash.Pairs[key] = value

        if !p.peekTokenIs(token.RCURLY) && !p.expectPeek(token.COMMA) {
            return nil
        }
    }

    if !p.expectPeek(token.RCURLY) {
        return nil
    }

    return hash
}

// Helper funcs
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.addError(t)
		return false
	}
}

//error funcs

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) addError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
