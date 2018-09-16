package parser

import (
	"fmt"
	"testing"

	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedExpression interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10", "y", 10},
		{"let foobar = 838383;", "foobar", 838383},
		{"let youUgly = true", "youUgly", true},
		{"let youFat = false;", "youFat", false},
		{"let foo = bar;", "foo", "bar"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedExpression) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.StatementNode, name string, value interface{}) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral not 'let'. got=%q", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatementNode)
	if !ok {
		t.Errorf("stmt not *ast.LetStatementNode. got=%T", stmt)
		return false
	}

	if !testIdentifier(t, letStmt.Name, name) {
		return false
	}

	if !testLiteralExpression(t, letStmt.Value, value) {
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedExpression string
	}{
		{"return 5;", "5"},
		{"return 5 + 10;", "10"},
		{"return true", "true"},
		{"return foo", "foo"},
		{"return add(1 + 2, 3, 4)", "add((1 + 2), 3, 4)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testReturnStatement(t, stmt, tt.expectedExpression) {
			return
		}
	}
}

func testReturnStatement(t *testing.T, stmt ast.StatementNode, expectedExprStr string) bool {

	if stmt.TokenLiteral() != "return" {
		t.Errorf("stmt.TokenLiteral not 'return'. got=%q", stmt.TokenLiteral())
		return false
	}

	retStmt, ok := stmt.(*ast.ReturnStatementNode)
	if !ok {
		t.Errorf("stmt not *ast.ReturnStatementNode. got=%T", stmt)
		return false
	}

	if retStmt.ReturnValue.String() != expectedExprStr {
		return false
	}

	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}

	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.ExpressionNode, value int64) bool {
	integ, ok := exp.(*ast.IntegerLiteralNode)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteralNode. got=%T", exp)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.ExpressionNode, value string) bool {
	ident, ok := exp.(*ast.IdentifierNode)
	if !ok {
		t.Errorf("exp not *ast.IdentifierNode. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.ExpressionNode, value bool) bool {
	boolNode, ok := exp.(*ast.BooleanNode)
	if !ok {
		t.Errorf("exp not *ast.BooleanNode. got=%T", exp)
		return false
	}

	if boolNode.Value != value {
		t.Errorf("boolNode.Value not %t. got=%t", value, boolNode.Value)
		return false
	}

	if boolNode.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("boolNode.TokenLiteral not %t. got=%s", value, boolNode.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.ExpressionNode, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func TestIdentifierLiteralExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	testIdentifier(t, expStmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `555;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	empStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatementNode. got=%T", program.Statements[0])
	}

	testIntegerLiteral(t, empStmt.Expression, 555)

}

func TestBooleanLiteralExpression(t *testing.T) {
	input := `true;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(program.Statements))
	}

	empStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatementNode. got=%T", program.Statements[0])
	}

	testBooleanLiteral(t, empStmt.Expression, true)

}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatementNode. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpressionNode)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testInfixExpression(t *testing.T, exp ast.ExpressionNode, left interface{}, operator string, right interface{}) bool {

	infixExp, ok := exp.(*ast.InfixExpressionNode)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpressionNode. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, infixExp.Left, left) {
		return false
	}

	if infixExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infixExp.Operator)

		return false
	}

	if !testLiteralExpression(t, infixExp.Right, right) {
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatementNode. got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, expStmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a * b",
			"(a * b);",
		},
		{
			"-a * b",
			"((-a) * b);",
		},
		{
			"!-a",
			"(!(-a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4 ; -5 * 5",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"true",
			"true;",
		},
		{
			"false",
			"false;",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false);",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true);",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5));",
		},
		{
			"!(true == true)",
			"(!(true == true));",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d);",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)));",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g));",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpressions(t *testing.T) {
	input := `if (x < y) {x}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] not a *ast.ExpressionStatementNode. got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpressionNode)
	if !ok {
		t.Fatalf("exprStmt.Expression not a *ast.IfExpressionNode. got=%T", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(ifExpr.Consequence.Statements))
	}

	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Errorf("Consequence.Statements[0] is not an *ast.ExpressionStatementNode. got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequenceExprStmt.Expression, "x") {
		return
	}

	if ifExpr.Alternative != nil {
		t.Errorf("ifExpr.Alternative is not nil. got=%+v", ifExpr.Alternative)
	}
}

func TestIfElseExpressions(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] not a *ast.ExpressionStatementNode. got=%T", program.Statements[0])
	}

	ifExpr, ok := exprStmt.Expression.(*ast.IfExpressionNode)
	if !ok {
		t.Fatalf("exprStmt.Expression not a *ast.IfExpressionNode. got=%T", exprStmt.Expression)
	}

	if !testInfixExpression(t, ifExpr.Condition, "x", "<", "y") {
		return
	}

	if len(ifExpr.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statement. got=%d", len(ifExpr.Consequence.Statements))
	}

	consequenceExprStmt, ok := ifExpr.Consequence.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Errorf("Consequence.Statements[0] is not an *ast.ExpressionStatementNode. got=%T", ifExpr.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequenceExprStmt.Expression, "x") {
		return
	}

	if len(ifExpr.Alternative.Statements) != 1 {
		t.Errorf("alternative is not 1 statement. got=%d", len(ifExpr.Alternative.Statements))
	}

	alternativeExprStmt, ok := ifExpr.Alternative.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Errorf("Alternative.Statements[0] is not an *ast.ExpressionStatementNode. got=%T", ifExpr.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternativeExprStmt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `func(x, y) {
		x + y
	}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("program.Statements[0] not a *ast.ExpressionStatementNode. got=%T", program.Statements[0])
	}

	funcExpr, ok := exprStmt.Expression.(*ast.FunctionLiteralNode)
	if !ok {
		t.Fatalf("exprStmt.Expression is not an *ast.FunctionLiteralNode. got=%T", exprStmt.Expression)
	}

	if len(funcExpr.Parameters) != 2 {
		t.Fatalf("len(funcExpr.Parameters) is not 2. got=%d", len(funcExpr.Parameters))
	}

	testIdentifier(t, funcExpr.Parameters[0], "x")
	testIdentifier(t, funcExpr.Parameters[1], "y")

	if len(funcExpr.Body.Statements) != 1 {
		t.Fatalf("funcExpr.Body.Statements has not 1 statements. got=%d\n", len(funcExpr.Body.Statements))
	}

	bodyStmt, ok := funcExpr.Body.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("funcExpr body stmt is not ast.ExpressionStatement. got=%T", funcExpr.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")

}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "func() {};", expectedParams: []string{}},
		{input: "func(x) {};", expectedParams: []string{"x"}},
		{input: "func(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatementNode)
		funcExpr := exprStmt.Expression.(*ast.FunctionLiteralNode)

		if len(funcExpr.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(funcExpr.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testIdentifier(t, funcExpr.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
	}

	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatementNode)
	if !ok {
		t.Fatalf("exprStmt is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	callExp, ok := exprStmt.Expression.(*ast.CallExpressionNode)
	if !ok {
		t.Fatalf("exprStmt.Expression is not ast.CallExpression. got=%T", exprStmt.Expression)
	}

	if !testIdentifier(t, callExp.Function, "add") {
		return
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(callExp.Arguments))
	}

	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionArgumentsParsing(t *testing.T) {
	tests := []struct {
		input        string
		expectedArgs []string
	}{
		{input: "add();", expectedArgs: []string{}},
		{input: "add(x * y);", expectedArgs: []string{"(x * y)"}},
		{input: "add(x * y, y, z);", expectedArgs: []string{"(x * y)", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		exprStmt := program.Statements[0].(*ast.ExpressionStatementNode)
		callExp := exprStmt.Expression.(*ast.CallExpressionNode)

		if len(callExp.Arguments) != len(tt.expectedArgs) {
			t.Errorf("length of arguments wrong. want %d, got=%d\n", len(tt.expectedArgs), len(callExp.Arguments))
		}

		for i, argExp := range tt.expectedArgs {
			if callExp.Arguments[i].String() != argExp {
				t.Errorf("callExp.Arguments[%d] is not %q. got=%q", i, argExp, callExp.Arguments[i].String())
			}
		}
	}
}
