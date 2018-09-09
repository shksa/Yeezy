package ast

import (
	"testing"

	"github.com/shksa/monkey/token"
)

func TestString(t *testing.T) {
	// Test for 'let myVar = anotherVar;'
	program := &Program{
		Statements: []StatementNode{
			&LetStatementNode{
				Token: token.LET,
				Name: &IdentifierNode{
					Token: token.Token{Type: "IDENTIFIER", Literal: "myVar"},
					Value: "myVar",
				},
				Value: &IdentifierNode{
					Token: token.Token{Type: "IDENTIFIER", Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() is wrong. got=%q \n", program.String())
	}
}
