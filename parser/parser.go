package parser

import (
	"fmt"

	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/lexer"
	"github.com/shksa/monkey/token"
)

// Parser is a type for representing the parser object.
type Parser struct {
	l                     *lexer.Lexer
	curToken              token.Token
	nextToken             token.Token
	Errors                []string
	ParseFnForPrefixToken map[string]prefixTokenParseFn
	ParseFnForInfixToken  map[string]infixTokenParseFn
}

// New returns a pointer to a newly created Parser object.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.readNextToken()
	p.readNextToken()                                             // To set the init value of curToken and the nextToken.
	p.ParseFnForPrefixToken = make(map[string]prefixTokenParseFn) // Need to assign a non-nil map, otherwise cannot assign to a nil nap.
	p.registerParseFuncForPrefixToken(token.IDENTIFIER, p.parseIdentifier)
	return p
}

// These are the 2 types of parsing functions associated with a token.
// prefixTokenParseFn is a function type that should be called when the associated token is in a prefix position.
// infixTokenParseFn is a function type that should be called when the associated token is in a infix position.
type (
	prefixTokenParseFn func() ast.ExpressionNode
	infixTokenParseFn  func(ast.ExpressionNode) ast.ExpressionNode
)

// registerParseFuncForPrefixToken adds entries to the parser's ParseFnForPrefixToken map.
func (p *Parser) registerParseFuncForPrefixToken(tok token.Token, fn prefixTokenParseFn) {
	p.ParseFnForPrefixToken[tok.Type] = fn
}

// registerParseFuncForInfixToken adds entries to the parser's ParseFnForInfixToken map.
func (p *Parser) registerParseFuncForInfixToken(tok token.Token, fn infixTokenParseFn) {
	p.ParseFnForInfixToken[tok.Type] = fn
}

func (p *Parser) readNextToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

// ParseProgram parses the whole program in a top-down recursive way to
// construct the full AST and returns it
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // root node of AST

	for !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.readNextToken()
	}

	return program
}

// parseStatement parses statements based on the current token info.
// Because the type of a statement is determined by it's FIRST token.
func (p *Parser) parseStatement() ast.StatementNode {
	switch p.curToken {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseLetStatement returns a pointer to a let statement node i.e *ast.LetStatementNode
func (p *Parser) parseLetStatement() *ast.LetStatementNode {
	stmt := &ast.LetStatementNode{Token: p.curToken}
	if !p.nextTokenIs(token.IDENTIFIER) {
		p.recordUnexpectedTokenError(token.IDENTIFIER)
		return nil
	}
	p.readNextToken()

	stmt.Name = &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}

	if !p.nextTokenIs(token.ASSIGN) {
		p.recordUnexpectedTokenError(token.ASSIGN)
		return nil
	}
	p.readNextToken()

	// TODO: SKIPPING THE EXPRESSIONS UNTILL A SEMICOLON IS ENCOUNTERED
	for !p.curTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(tok token.Token) bool {
	return p.curToken.Type == tok.Type
}

func (p *Parser) nextTokenIs(tok token.Token) bool {
	return p.nextToken.Type == tok.Type
}

func (p *Parser) recordUnexpectedTokenError(expectedTok token.Token) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", expectedTok.Type, p.nextToken.Type)
	p.Errors = append(p.Errors, msg)
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatementNode {
	stmt := &ast.ReturnStatementNode{Token: p.curToken}

	p.readNextToken()

	// TODO: SKIPPING THE EXPRESSIONS UNTILL A SEMICOLON IS ENCOUNTERED
	for p.curToken != token.SEMICOLON {
		p.readNextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatementNode {
	stmt := &ast.ExpressionStatementNode{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextTokenIs(token.SEMICOLON) { // semicolon is optional
		p.readNextToken()
	}

	return stmt
}

// Operator precedences. With these constants we can answer questions like :-
// i) Does the * operator have a higher precedence than the == operator?
// ii) Does a prefix operator have a higher preference than a call expression?
const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < or >
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

// parseExpression does the following:-
// 1. Checks whether the current token in prefix position has a parsing function associated with it.
// 2. If it does, it calls the parsing function and that returns an expression node.
func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	parseFn := p.ParseFnForPrefixToken[p.curToken.Type]
	if parseFn == nil {
		return nil
	}
	leftExpr := parseFn()

	return leftExpr
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}
}
