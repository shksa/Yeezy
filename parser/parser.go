package parser

import (
	"fmt"
	"strconv"

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
	p.registerParseFuncForPrefixToken(token.INT, p.parseIntegerLiteral)
	p.registerParseFuncForPrefixToken(token.BANG, p.parsePrefixExpression)
	p.registerParseFuncForPrefixToken(token.MINUS, p.parsePrefixExpression)
	p.ParseFnForInfixToken = make(map[string]infixTokenParseFn)
	p.registerParseFuncForInfixToken(token.PLUS, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.MINUS, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.ASTERISK, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.LT, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.GT, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.EQ, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.NOTEQ, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.LPAREN, p.parseInfixExpression)
	p.registerParseFuncForInfixToken(token.SLASH, p.parseInfixExpression)
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

func (p *Parser) noPrefixParseFnError(tok token.Token) {
	msg := fmt.Sprintf("No parse function found for %q token at prefix position", tok.Literal)
	p.Errors = append(p.Errors, msg)
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
		// At this point, p.curToken is a token.SEMICOLON
		p.readNextToken()
		// After p.readNextToken(), p.curToken is the first token of a new statement.
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
	letStmtNode := &ast.LetStatementNode{Token: p.curToken}
	if !p.nextTokenIs(token.IDENTIFIER) {
		p.unexpectedTokenError(token.IDENTIFIER)
		return nil
	}
	p.readNextToken()

	letStmtNode.Name = &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}

	if !p.nextTokenIs(token.ASSIGN) {
		p.unexpectedTokenError(token.ASSIGN)
		return nil
	}
	p.readNextToken()

	if !(p.nextTokenIs(token.INT) || p.nextTokenIs(token.IDENTIFIER) || p.nextTokenIs(token.BANG) || p.nextTokenIs(token.MINUS)) {
		p.unexpectedTokenError(token.Token{Type: "EXPRESSION"})
		return nil
	}
	p.readNextToken()

	letStmtNode.Value = p.parseExpression(LOWEST)
	// TODO: SKIPPING THE EXPRESSIONS UNTILL A SEMICOLON IS ENCOUNTERED
	// for !p.curTokenIs(token.SEMICOLON) {
	// 	p.readNextToken()
	// }
	if !p.nextTokenIs(token.SEMICOLON) {
		p.unexpectedTokenError(token.SEMICOLON)
		return nil
	}
	p.readNextToken() //p.curToken should token.SEMICOLON at the end of parsing a statement

	return letStmtNode
}

func (p *Parser) curTokenIs(tok token.Token) bool {
	return p.curToken.Type == tok.Type
}

func (p *Parser) nextTokenIs(tok token.Token) bool {
	return p.nextToken.Type == tok.Type
}

func (p *Parser) unexpectedTokenError(expectedTok token.Token) {
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

	if !p.nextTokenIs(token.SEMICOLON) {
		p.unexpectedTokenError(token.SEMICOLON)
		return nil
	}
	p.readNextToken() //p.curToken should token.SEMICOLON at the end of parsing a statement

	return stmt
}

// Operator precedences. With these constants we can answer questions like :-
// i) Does the * operator have a higher precedence than the == operator?
// ii) Does a prefix operator have a higher preference than a call expression?
const (
	_           int = iota
	LOWEST          // #1
	EQUALS          // ==, #2
	LESSGREATER     // < or >, #3
	SUM             // +, #4
	PRODUCT         // *, #5
	PREFIX          // -X or !X, #6
	CALL            // myFunction(X), #7
)

// parseExpression does the following:-
// 1. Checks whether the current token in prefix position has a parsing function associated with it.
// 2. If it does, it calls the parsing function and that returns an expression node.
// 3. Checks if the next token has higher precedence than it's precedence parameter.
// 4. If it does, then the expression node built before becomes the left arm of the new infix node that
//		will be built by calling the next token's infix parse function.
// 5. point 3 continues untill a semicolon is encountered or the next token's precedence is lower than the
//		func's precedence parameter.
// 6. Finally returns the parsed expression node.
func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefixParseFn := p.ParseFnForPrefixToken[p.curToken.Type]
	if prefixParseFn == nil {
		p.noPrefixParseFnError(p.curToken)
		return nil
	}
	leftExprNode := prefixParseFn()

	// Say there are 3 tokens. 1 + 2 * 3
	// And p.curToken reps 2, which means p.nextToken reps *
	// leftExprNode contains the integerLiteralNode for the 2.
	// The loop condition checks the following,
	// 1. If the left-binding power of the * token is greater than the right binding power of +,
	//		then the node for 2 becomes the left arm of the infix expression with * as the infix operator.
	// That means the parsed expression would be nested this way -> 1 + (2 * 3)
	for !p.nextTokenIs(token.SEMICOLON) && precedence < p.nextTokenPrecedence() {
		infixParseFn := p.ParseFnForInfixToken[p.nextToken.Type]
		if infixParseFn == nil {
			return leftExprNode
		}

		p.readNextToken()

		leftExprNode = infixParseFn(leftExprNode)
	}

	return leftExprNode
}

func (p *Parser) parseIdentifier() ast.ExpressionNode {
	return &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.ExpressionNode {
	intLitNode := &ast.IntegerLiteralNode{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		errMsg := fmt.Sprintf("cannot parse %q as an int64", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)
		return nil
	}
	intLitNode.Value = value
	return intLitNode
}

func (p *Parser) parsePrefixExpression() ast.ExpressionNode {
	prefixExprNode := &ast.PrefixExpressionNode{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.readNextToken()
	prefixExprNode.Right = p.parseExpression(PREFIX)
	return prefixExprNode
}

var precedences = map[token.Token]int{
	token.EQ:       EQUALS,      // 2
	token.NOTEQ:    EQUALS,      // 2
	token.LT:       LESSGREATER, // 3
	token.GT:       LESSGREATER, // 3
	token.PLUS:     SUM,         // 4
	token.MINUS:    SUM,         // 4
	token.SLASH:    PRODUCT,     // 5
	token.ASTERISK: PRODUCT,     // 5
	token.LPAREN:   CALL,        // 7
}

func (p *Parser) curTokenPrecedence() int {
	if p, ok := precedences[p.curToken]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) nextTokenPrecedence() int {
	if p, ok := precedences[p.nextToken]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.ExpressionNode) ast.ExpressionNode {
	infixExprNode := &ast.InfixExpressionNode{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}

	infixOpPrecedence := p.curTokenPrecedence()
	p.readNextToken()
	infixExprNode.Right = p.parseExpression(infixOpPrecedence)

	return infixExprNode
}
