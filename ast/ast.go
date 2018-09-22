package ast

import (
	"bytes"
	"strings"

	"github.com/shksa/monkey/token"
)

// Node is an interface type for representing the interface for all nodes in the AST
type Node interface {
	TokenLiteral() string
	String() string
}

// StatementNode is an interface type for representing all statement nodes in the AST.
type StatementNode interface {
	Node
	statementNode()
}

// ExpressionNode is an interface type for representing all expression nodes in the AST.
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

	outString := out.String()
	return outString
}

/* NOTE
1. A statement is identified with the token it starts with.
	Example:- A let statement, starts with the LET token, A return statement, starts with the RETURN token.
2. So a statement's node will contain the token that identifies that statement.
*/

// LetStatementNode is a type for representing all "let" statements in AST. ex:= `let x = 5 * 6`
type LetStatementNode struct {
	Token token.Token // token.LET
	Iden  *IdentifierNode
	Value ExpressionNode
}

// *LetStatementNode implements StatementNode interface.
func (ls *LetStatementNode) statementNode() {}

// TokenLiteral returns the LetStatementNode's token literal.
func (ls *LetStatementNode) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Iden.String())
	out.WriteString(" = ")
	out.WriteString(ls.Value.String())

	out.WriteString(";")
	outString := out.String()
	return outString
}

// IdentifierNode is a type for representing all "identifier" expressions in AST.
type IdentifierNode struct {
	Token token.Token // token.IDENTIFIER
	Name  string      // Name is the Token.Literal
}

// *LetStatementNode implements ExpressionNode interface.
func (i *IdentifierNode) expressionNode() {}

// TokenLiteral returns the IdentifierNode's token literal.
func (i *IdentifierNode) TokenLiteral() string { return i.Token.Literal }

func (i *IdentifierNode) String() string { return i.Name }

// ReturnStatementNode is a type for representing all "return" statements in AST. ex:- return 777
type ReturnStatementNode struct {
	Token       token.Token // token.RETURN
	ReturnValue ExpressionNode
}

// *ReturnStatementNode implements StatementNode interface.
func (rs *ReturnStatementNode) statementNode() {}

// TokenLiteral returns the ReturnStatementNode's token literal.
func (rs *ReturnStatementNode) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatementNode) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")
	out.WriteString(rs.ReturnValue.String())
	out.WriteString(";")
	outString := out.String()
	return outString
}

// ExpressionStatementNode is a type for representing all "expression" statements in AST. ex:- (in top-level) 5 * 5 + 10
type ExpressionStatementNode struct {
	Token      token.Token // The first token of the expression
	Expression ExpressionNode
}

// *ExpressionStatementNode implements StatementNode interface.
func (es *ExpressionStatementNode) statementNode() {}

// TokenLiteral returns the ExpressionStatementNode's token literal.
func (es *ExpressionStatementNode) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatementNode) String() string {
	var out bytes.Buffer
	out.WriteString(es.Expression.String())
	out.WriteString(";")
	return out.String()
}

// IntegerLiteralNode is a type for representing all "integer" literal expressions in AST.
type IntegerLiteralNode struct {
	Token token.Token // token.INT
	Value int64
}

// TokenLiteral returns the IntegerLiteralNode's token literal.
func (il *IntegerLiteralNode) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteralNode) expressionNode()      {}
func (il *IntegerLiteralNode) String() string       { return il.Token.Literal }

// PrefixExpressionNode is a type for representing all "prefix" expressions in AST.
type PrefixExpressionNode struct {
	Token    token.Token    // The prefix token ex:- token.BANG or token.MINUS.
	Operator string         // "!" or "-".
	Right    ExpressionNode // expression to the right of the operator.
}

// TokenLiteral returns the PrefixExpressionNode's token literal.
func (pe *PrefixExpressionNode) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpressionNode) expressionNode()      {}
func (pe *PrefixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	outString := out.String()
	return outString
}

// InfixExpressionNode is a type for representing all "infix" expressions in AST.
type InfixExpressionNode struct {
	Token    token.Token // The operator token, for 5 + 10, it's token.PLUS
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

// TokenLiteral returns the InfixExpressionNode's token literal.
func (ie *InfixExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpressionNode) expressionNode()      {}
func (ie *InfixExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	outString := out.String()
	return outString
}

// BooleanNode is a type for representing all "boolean" literal expressions in AST.
type BooleanNode struct {
	Token token.Token
	Value bool // The go bool values, true or false
}

// TokenLiteral returns the BooleanNode's token literal.
func (b *BooleanNode) TokenLiteral() string { return b.Token.Literal }
func (b *BooleanNode) String() string       { return b.Token.Literal }
func (b *BooleanNode) expressionNode()      {}

// IfExpressionNode is a type for representing all "if" expressions in AST.
type IfExpressionNode struct {
	Token       token.Token // The if token
	Condition   ExpressionNode
	Consequence *BlockStatementNode
	Alternative *BlockStatementNode
}

// TokenLiteral returns the IfExpressionNode's token literal.
func (ie *IfExpressionNode) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpressionNode) expressionNode()      {}
func (ie *IfExpressionNode) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(" ( " + ie.Condition.String() + " ) ")
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// BlockStatementNode is a type for representing all "block" statements in AST.
type BlockStatementNode struct {
	Token      token.Token // The { token
	Statements []StatementNode
}

// TokenLiteral returns the BlockStatementNode's token literal.
func (bs *BlockStatementNode) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatementNode) statementNode()       {}
func (bs *BlockStatementNode) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, stmt := range bs.Statements {
		out.WriteString(stmt.String())
	}
	out.WriteString("}")

	return out.String()
}

// FunctionLiteralNode is a type for representing all "function literal" expressions in AST.
type FunctionLiteralNode struct {
	Token      token.Token // the "func" token
	Parameters []*IdentifierNode
	Body       *BlockStatementNode
}

// TokenLiteral returns the FunctionLiteralNode's token literal.
func (fl *FunctionLiteralNode) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteralNode) expressionNode()      {}
func (fl *FunctionLiteralNode) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fl.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpressionNode is a type for representing all "call" expressions in AST.
type CallExpressionNode struct {
	Token     token.Token    // The left paren "(" token
	Function  ExpressionNode // either IdentifierNode or FunctionLiteralNode
	Arguments []ExpressionNode
}

// TokenLiteral returns the CallExpressionNode's token literal.
func (ce *CallExpressionNode) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpressionNode) expressionNode()      {}
func (ce *CallExpressionNode) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, p := range ce.Arguments {
		args = append(args, p.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// StringLiteralNode is a type for representing all "string" literal expressions in ast.
type StringLiteralNode struct {
	Token token.Token // token.STRING
	Value string
}

// TokenLiteral returns the StringLiteralNode's token literal.
func (sn *StringLiteralNode) TokenLiteral() string { return sn.Token.Literal }
func (sn *StringLiteralNode) expressionNode()      {}
func (sn *StringLiteralNode) String() string       { return sn.Token.Literal }
