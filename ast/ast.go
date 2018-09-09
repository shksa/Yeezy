package ast

import (
	"bytes"

	"github.com/shksa/monkey/token"
)

// Node is an interface type for representing the interface for all nodes in the AST
type Node interface {
	TokenLiteral() string
	String() string
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
// Program is a type for representing the whole program tree.
type Program struct {
	Statements []StatementNode // slice of AST node pointers that implement the StatementNode interface
}

// TokenLiteral returns the token literal of the first statement the program holds.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// String will output the whole program's source code back as it is.
// This makes testing the structure of the AST very simple and easy.
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

/* NOTE
1. A statement is identified with the token it starts with.
	Example:- A let statement, starts with the LET token, A return statement, starts with the RETURN token.
2. So a statement's node will contain the token that identifies that statement.
*/

// LetStatementNode is a type for representing all "let" statement nodes. ex:= `let x = 5 * 6`
type LetStatementNode struct {
	Token token.Token // token.LET
	Name  *IdentifierNode
	Value ExpressionNode
}

// *LetStatementNode implements StatementNode interface.
func (ls *LetStatementNode) statementNode() {}

// TokenLiteral returns the LetStatementNode's token literal.
func (ls *LetStatementNode) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// IdentifierNode is a type for representing all identifier nodes.
type IdentifierNode struct {
	Token token.Token // token.IDENTIFIER
	Value string      // Value is the Token.Literal
}

// *LetStatementNode implements ExpressionNode interface.
func (i *IdentifierNode) expressionNode() {}

// TokenLiteral returns the IdentifierNode's token literal.
func (i *IdentifierNode) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IdentifierNode) String() string {
	return i.Value
}

// ReturnStatementNode is a type for representing all "return" statement nodes. ex:- return 777
type ReturnStatementNode struct {
	Token       token.Token // token.RETURN
	ReturnValue ExpressionNode
}

// *ReturnStatementNode implements StatementNode interface.
func (rs *ReturnStatementNode) statementNode() {}

// TokenLiteral returns the ReturnStatementNode's token literal.
func (rs *ReturnStatementNode) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// ExpressionStatementNode is a type for representing all "expression" statement nodes. ex:- (in top-level) 5 * 5 + 10
type ExpressionStatementNode struct {
	Token      token.Token // The first token of the expression
	Expression ExpressionNode
}

// *ExpressionStatementNode implements StatementNode interface.
func (es *ExpressionStatementNode) statementNode() {}

// TokenLiteral returns the ExpressionStatementNode's token literal.
func (es *ExpressionStatementNode) TokenLiteral() string {
	return es.Token.Literal
}
