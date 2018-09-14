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
	p.registerParseFuncForPrefixToken(token.TRUE, p.parseBooleanLiteral)
	p.registerParseFuncForPrefixToken(token.FALSE, p.parseBooleanLiteral)
	p.registerParseFuncForPrefixToken(token.LPAREN, p.parseGroupedExpression)
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
	msg := fmt.Sprintf("No prefix parse function found for %s token", tok.Literal)
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
		stmt := p.parseStatement() // At this point, p.curToken will be the start token of a new statement. <- IMP. INVARIANT
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		// For single-line inputs, at this point, p.nextToken is always token.EOF <- IMP. INVARIANT
		// For multi-line inputs with semicolons,
		//		p.curToken is always token.SEMICOLON <- IMP. INVARIANT
		//		p.nextToken is the token at the start of the next line <- IMP. INVARIANT
		p.readNextToken()
		// For single-line inputs, at this point, p.curToken is always token.EOF <- IMP. INVARIANT
		// For multi-line inputs with semicolons,
		// 		p.curToken is always the token at the start of the next statement <- IMP. INVARIANT
	}

	return program
}

/*
IMPORTANT: All parsing is done from LEFT to RIGHT.
*/

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

func (p *Parser) expectAndReadNextTokenToBe(tok token.Token) bool {
	if !p.nextTokenIs(tok) {
		p.unexpectedTokenError(tok)
		return false
	}
	p.readNextToken()
	return true
}

// parseLetStatement returns a pointer to a let statement node i.e *ast.LetStatementNode
func (p *Parser) parseLetStatement() *ast.LetStatementNode {
	letStmtNode := &ast.LetStatementNode{Token: p.curToken}

	if didRead := p.expectAndReadNextTokenToBe(token.IDENTIFIER); !didRead {
		return nil
	}

	letStmtNode.Name = &ast.IdentifierNode{Token: p.curToken, Value: p.curToken.Literal}

	if didRead := p.expectAndReadNextTokenToBe(token.ASSIGN); !didRead {
		return nil
	}

	// next token is the expression ot start of an expression
	p.readNextToken()

	letStmtNode.Value = p.parseExpression(LOWEST)

	if didRead := p.expectAndReadNextTokenToBe(token.SEMICOLON); !didRead {
		return nil
	}

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
	retStmtNode := &ast.ReturnStatementNode{Token: p.curToken}

	p.readNextToken()

	// TODO: SKIPPING THE EXPRESSIONS UNTILL A SEMICOLON IS ENCOUNTERED
	for p.curToken != token.SEMICOLON {
		p.readNextToken()
	}
	return retStmtNode
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatementNode {
	exprStmtNode := &ast.ExpressionStatementNode{Token: p.curToken}

	exprStmtNode.Expression = p.parseExpression(LOWEST)

	// semicolon is optional for single-line inputs.
	// 1.	For a single-line input without a semicolon at the end,
	//			p.curToken will be the last token.
	// 			p.nextToken will be token.EOF.
	// 2. For a single-line input with a semicolon at the end,
	// 			p.curToken will be the last token.
	// 			p.nextToken will be token.SEMICOLON.
	//			calling p.readNextToken() will make p.curToken = token.SEMICOLON
	//			and p.nextToken = token.EOF
	// So for a single-line input, after this if-block,
	// 			p.curToken = last token in the line. ex- ; or 55 or foobar
	//			p.nextToken = always token.EOF. <- IMP. INVARIANT FOR SINGLE-LINE INPUTS
	if p.nextTokenIs(token.SEMICOLON) {
		p.readNextToken()
	}

	return exprStmtNode
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
	for !p.nextTokenIs(token.SEMICOLON) && precedence < p.nextTokenPrecedence() { // The next token can be a semicolon or a eof or a RPAREN
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
	intLiteralNode := &ast.IntegerLiteralNode{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		errMsg := fmt.Sprintf("cannot parse %q as an int64", p.curToken.Literal)
		p.Errors = append(p.Errors, errMsg)
		return nil
	}
	intLiteralNode.Value = value
	return intLiteralNode
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
	return LOWEST // This is returned when the p.nextToken is a token.EOF or token.SEMICOLON.
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

func (p *Parser) parseBooleanLiteral() ast.ExpressionNode {
	return &ast.BooleanNode{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.ExpressionNode {
	p.readNextToken()

	exp := p.parseExpression(LOWEST) // Why LOWEST? Because we want the expression in the paren not be swayed by the precedence power
	// of the operator that is to the left of the paren. Therefore LOWEST signifies to the parseExpression function that the operator
	// in the left of the expression is of the lowest precedence, so that the expression will not be under the right-binding power of
	// that operator.

	if didRead := p.expectAndReadNextTokenToBe(token.RPAREN); !didRead {
		return nil
	}

	return exp
}
