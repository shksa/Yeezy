package ast

import (
	"testing"

	"github.com/shksa/yeezy/token"
)

func TestString(t *testing.T) {
	// Test for 'let myVar = anotherVar;'
	program := &Program{
		Statements: []StatementNode{
			&LetStatementNode{
				Token: token.LET,
				Iden: &IdentifierNode{
					Token: token.Token{Type: "IDENTIFIER", Literal: "myVar"},
					Name:  "myVar",
				},
				Value: &IdentifierNode{
					Token: token.Token{Type: "IDENTIFIER", Literal: "anotherVar"},
					Name:  "anotherVar",
				},
			},
		},
	}
	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() is wrong. got=%q \n", program.String())
	}
}
