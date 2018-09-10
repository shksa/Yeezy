package parser

import (
	"fmt"
	"testing"

	"github.com/shksa/monkey/ast"
	"github.com/shksa/monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedExpression int
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier, tt.expectedExpression) {
			return
		}
	}
}

func testLetStatement(t *testing.T, stmt ast.StatementNode, name string, value int) bool {
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

	if !testIntegerLiteral(t, letStmt.Value, int64(value)) {
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

func TestReturnStatements(t *testing.T) {
	input := `
	return 5;
	return 5 + 10;
	return add(5, 10);
	`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatementNode)
		if !ok {
			t.Errorf("stmt is not a *ast.ReturnStatementNode. got=%T", stmt)
		}

		if retStmt.TokenLiteral() != "return" {
			t.Errorf("retStmt.TokenLiteral() is not 'return'. got=%s", retStmt.TokenLiteral())
		}
	}
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
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
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
			"a * b;",
			"(a * b)",
		},
		{
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true;",
			"true",
		},
		{
			"false;",
			"false",
		},
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true;",
			"((3 < 5) == true)",
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
