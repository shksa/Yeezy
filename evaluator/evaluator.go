package evaluator

import (
	"fmt"

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
	object.Boolean value or a object.NULL value
*/

/* Evaluation of statements inside Blocks
- Inside blocks with return statements, the evaluation of a return statement will return a object.ReturnValue.
- The evaluation of block statements is such that when a return statement is evaluated, the result which is an object.ReturnValue,
	is returned and therefore the evaluation of rest of the statements in the block are skipped.
- So the evaluation of such block statements will return a object.ReturnValue and will reach the top-level/program evaluation because
	the if-expressions which consist of block statements, will return whatever the block statement returns.
- In the top-level, the object.ReturnValue it will be unwrapped to get the actual value and will be returned to the user.
*/

// Eval takes in the AST and evaluates it, returning Monkey objects
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	// Statements
	case *ast.Program:
		return evaluateProgram(node.Statements)

	case *ast.ExpressionStatementNode:
		return Eval(node.Expression)

	case *ast.BlockStatementNode:
		return evaluateBlockStatement(node) // Can return a *object.ReturnValue

	case *ast.ReturnStatementNode:
		val := Eval(node.ReturnValue)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val} // Need to keep track of return value so that we can decide later whether to stop evaluation or not

	// Expressions
	case *ast.IntegerLiteralNode:
		return &object.Integer{Value: node.Value}

	case *ast.BooleanNode:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpressionNode:
		operand := Eval(node.Right) // operand may be object.Integer, object.Boolean, or object.Null, object.Error
		if isError(operand) {
			return operand
		}
		return evaluatePrefixExpression(node.Operator, operand)

	case *ast.InfixExpressionNode:
		leftOperand := Eval(node.Left) // leftOperand may be object.Integer, object.Boolean, or object.Null, object.Error
		if isError(leftOperand) {
			return leftOperand
		}

		rightOperand := Eval(node.Right) // rightOperand may be object.Integer, object.Boolean, or object.Null, object.Error
		if isError(rightOperand) {
			return rightOperand
		}

		return evaluateInfixExpression(node.Operator, leftOperand, rightOperand)

	case *ast.IfExpressionNode:
		return evaluateIfExpression(node) // If-expression will return whatever its block statement will return.
	}

	return nil
}

func evaluateProgram(stmtNodes []ast.StatementNode) object.Object {
	var result object.Object

	for _, stmtNode := range stmtNodes {
		result = Eval(stmtNode)

		switch result := result.(type) {
		case *object.ReturnValue: // Evaluation of further statements is ended because a return statement is encountered.
			return result.Value // The actual Object value is unwrapped from the object.ReturnValue.

		case *object.Error: // Evaluation of further statements is ended because an error is encountered.
			return result
		}
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
		return newError("unknown operator: %s%s", operator, operand.Type())
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

	default: // For integers
		return FALSE
	}
}

func evaluateMinusPrefixOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER {
		return newError(`invalid prefix operator "-" for operand type %s`, operand.Type())
	}
	value := operand.(*object.Integer).Value

	return &object.Integer{Value: -value}
	// This is where Go is performing the negation operation.
	// ex:- for operand = 5, -5 is returned, for operand = -5, +5 is returned.
	// Go knows how to do integer arithmetic, so we make Go do it.
}

func evaluateInfixExpression(operator string, leftOperand, rightOperand object.Object) object.Object {
	switch {
	case leftOperand.Type() != rightOperand.Type():
		return newError("operand type mismatch for operator %q : %s %s %s", operator, leftOperand.Type(), operator, rightOperand.Type())

	case leftOperand.Type() == object.INTEGER && rightOperand.Type() == object.INTEGER:
		return evaluateIntegerInfixExpression(operator, leftOperand, rightOperand)

	// For the next cases, the leftOperand and rightOperand are *object.Boolean, either TRUE or FALSE values
	case operator == "==":
		return nativeBoolToBooleanObject(leftOperand == rightOperand) // Pointer comparision to check for equality b/w 2 boolean object pointers.

	case operator == "!=":
		return nativeBoolToBooleanObject(leftOperand != rightOperand)

	default:
		return newError("invalid operator %q between %s values: %s %s %s", operator, leftOperand.Type(), leftOperand.Inspect(), operator, rightOperand.Inspect())
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
		return newError("invalid operator %q between %s values: %s %s %s", operator, leftOperand.Type(), leftValue, operator, rightValue)
	}
}

func evaluateIfExpression(node *ast.IfExpressionNode) object.Object {
	conditionValue := Eval(node.Condition)

	if isError(conditionValue) {
		return conditionValue
	}

	if isTruthy(conditionValue) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case TRUE:
		return true
	case FALSE:
		return false
	case NULL:
		return false
	default:
		return true
	}
}

func evaluateBlockStatement(block *ast.BlockStatementNode) object.Object { // can return object.ReturnValue if the block has return statements.
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			resultType := result.Type()
			if resultType == object.ERROROBJ || resultType == object.RETURNOBJ {
				return result // The return value is not explicitly unwrapped and returned as is. (as *object.ReturnValue or *object.Error)
			}
		}
	}

	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROROBJ
	}
	return false
}
