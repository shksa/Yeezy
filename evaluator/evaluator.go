package evaluator

import (
	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/object"
)

// TRUE and False are refernences to the two boolean objects in Monkey
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

/* IMPORTANT
- The Host language Golang knows how to perform integer arithmetic operations, so in evaluating Monkey code, we use Golang's operators
	to perform the the operations.
- Eval() on a node as a root of a branch evaluates the whole branch to a single value like an object.Integer value or a
	object.Boolean value

*/

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

	case *ast.BooleanNode:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpressionNode:
		operandExpr := Eval(node.Right) // operandExpr may be object.Integer, object.Boolean, or object.Null
		return evaluatePrefixExpression(node.Operator, operandExpr)

	case *ast.InfixExpressionNode:
		leftOperand := Eval(node.Left)   // leftOperand may be object.Integer, object.Boolean, or object.Null
		rightOperand := Eval(node.Right) // rightOperand may be object.Integer, object.Boolean, or object.Null
		return evaluateInfixExpression(node.Operator, leftOperand, rightOperand)
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

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evaluatePrefixExpression(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		return evaluateBangPrefixOperatorExpression(operand)

	case "-":
		return evaluateMinusPrefixOperatorExpression(operand)

	default:
		return NULL
	}
}

func evaluateBangPrefixOperatorExpression(operand object.Object) *object.Boolean {
	switch operand {
	case TRUE:
		return FALSE

	case FALSE:
		return TRUE

	case NULL:
		return TRUE

	default:
		return FALSE
	}
}

func evaluateMinusPrefixOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER {
		return NULL
	}
	value := operand.(*object.Integer).Value

	return &object.Integer{Value: -value}
	// This is where Go is performing the negation operation.
	// ex:- for operand = 5, -5 is returned, for operand = -5, +5 is returned.
	// Go knows how to do integer arithmetic, so we make Go do it.
}

func evaluateInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	switch {
	case leftOperand.Type() == object.INTEGER && rightOperand.Type() == object.INTEGER:
		return evaluateIntegerInfixExpression(operator, leftOperand, rightOperand)

	// For the next cases, the leftOperand and rightOperand are *object.Boolean, either TRUE or FALSE
	case operator == "==":
		return nativeBoolToBooleanObject(leftOperand == rightOperand) // Pointer comparision to check for equality b/w 2 boolean object pointers.

	case operator == "!=":
		return nativeBoolToBooleanObject(leftOperand != rightOperand)

	default:
		return NULL
	}
}

func evaluateIntegerInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	leftValue := leftOperand.(*object.Integer).Value
	rightValue := rightOperand.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue} // Go is performing the addition operation.
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return NULL
	}
}
