package parser

import (
	"fmt"

	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/lexer"
	"github.com/shksa/monkey/token"
)

// Parser is a type for representing the parser object.
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	Errors    []string
}

// New returns a pointer to a newly created Parser object.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.readNextToken()
	p.readNextToken() // To set the curToken and the nextToken
	return p
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
	default:
		return nil
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
