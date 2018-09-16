package evaluator

import (
	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/object"
)

// Eval takes in the AST and evaluates it, returning Monkey objects
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evaluateStatements(node.Statements)

	case *ast.ExpressionStatementNode:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteralNode:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evaluateStatements(stmtNodes []ast.StatementNode) object.Object {
	var result object.Object

	for _, stmtNode := range stmtNodes {
		result = Eval(stmtNode)
	}

	return result
}
