package lexer

import (
	"fmt"
	"testing"

	"github.com/shksa/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+(){},;`

	// tests is a list of output expectations.
	tests := []struct {
		Type, Literal string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLAN, ";"},
	}

	lexer := New(input)

	for i, expected := range tests {
		tok := lexer.NextToken()
		fmt.Printf("%q \n", tok.Literal)
		if tok.Type != expected.Type {
			t.Fatalf("tests[%d] - token.Type is wrong. expected %q, got %q", i, expected.Type, tok.Type)
		}

		if tok.Literal != expected.Literal {
			t.Fatalf("tests[%d] - token.Literal is wrong. expected %q, got %q", i, expected.Literal, tok.Literal)
		}
	}
}
