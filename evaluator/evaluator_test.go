package evaluator

import (
	"testing"

	"github.com/shksa/monkey/lexer"
	"github.com/shksa/monkey/object"
	"github.com/shksa/monkey/parser"
)

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	newEnv := object.NewEnvironment()
	return Eval(program, newEnv)
}

func TestIntegerExpressionEval(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput int64
	}{
		{"5", 5},
		{"25", 25},
		{"-5", -5},
		{"-25", -25},
		{"1 + 2 + 3", 6},
		{"2 * 2 * 3", 12},
		{"-50 + 0 + 50", 0},
		{"5 * 2 + 1", 11},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expectedOutput)
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestBooleanExpressionEval(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expectedOutput)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput bool
	}{
		{"!false", true},
		{"!true", false},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expectedOutput)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	return obj == NULL
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
				return 10;
				}
				return 1;
			}
			`,
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			`operand type mismatch for operator "+" : INTEGER + BOOLEAN`,
		},
		{
			"5 + true; 5;",
			`operand type mismatch for operator "+" : INTEGER + BOOLEAN`,
		},
		{
			"-true",
			`invalid prefix operator "-" for operand type BOOLEAN`,
		},
		{
			"true + false;",
			`invalid operator "+" between BOOLEAN values: true + false`,
		},
		{
			"5; true + false; 5",
			`invalid operator "+" between BOOLEAN values: true + false`,
		},
		{
			"if (10 > 1) { true + false; }",
			`invalid operator "+" between BOOLEAN values: true + false`,
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			`invalid operator "+" between BOOLEAN values: true + false`,
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)

		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x){ x + 2; };"

	evaluated := testEval(input)

	fnObj, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not a function. got=%T (%+v)", fnObj, fnObj)
	}

	if len(fnObj.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fnObj.Parameters)
	}

	if fnObj.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fnObj.Parameters[0])
	}

	expectedBody := "{(x + 2);}"

	if fnObj.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fnObj.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(a) { return a; }; identity(5);", 5},
		{"let double = func(b) { b * 2; }; double(5);", 10},
		{"let add = func(c, d) { c + d; }; add(5, 5);", 10},
		{"let add = func(e, f) { e + f; }; add(5 + 5, add(5, 5));", 20},
		{"func(g) { g; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
		let newAdder = func(x) {
			func(y) { x + y };
		};

		let addTwo = newAdder(2);
		
		addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}
