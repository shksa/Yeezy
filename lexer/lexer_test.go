package lexer

import (
	"fmt"
	"testing"

	"github.com/shksa/monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
	let ten = 10;
	let add = func(x, y) {
		x + y;
	};
	let result = add(five, ten);`

	// tests is a list of output expectations.
	tests := []struct {
		Type, Literal string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
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
