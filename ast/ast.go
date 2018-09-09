package ast

import "github.com/shksa/monkey/token"

// Node interface represents the interface for every node in the AST
type Node interface {
	TokenLiteral() string
}

// StatementNode is an interface type for representing the interface of all statement nodes in the AST.
type StatementNode interface {
	Node
	statementNode()
}

// ExpressionNode is an interface type for representing the interface of all expression nodes in the AST.
type ExpressionNode interface {
	Node
	expressionNode()
}

// A Program in Monkey is a series of statements.
// This Program node will be root of the AST.
// Program is a type for representing the whole program tree
type Program struct {
	Statements []StatementNode // slice of AST nodes that implement the StatementNode interface
}

// TokenLiteral returns the token literal of the first statement the program holds.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatementNode is a type for representing all "let" statement nodes. ex:= `let x = 5 * 6`
type LetStatementNode struct {
	Token      token.Token // token.LET
	Name       *IdentifierNode
	Expression ExpressionNode
}

func (ls *LetStatementNode) statementNode() {}

// TokenLiteral returns the LetStatementNode's token literal.
func (ls *LetStatementNode) TokenLiteral() string {
	return ls.Token.Literal
}

// IdentifierNode is a type for representing all identifier nodes.
type IdentifierNode struct {
	Token token.Token // token.IDENTIFIER
	Value string      // Value is the Token.Literal
}

func (i *IdentifierNode) expressionNode() {}

// TokenLiteral returns the IdentifierNode's token literal.
func (i *IdentifierNode) TokenLiteral() string {
	return i.Token.Literal
}
