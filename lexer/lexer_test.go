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
	let result = add(five, ten);
	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}
	5 == 5;
	5 != 10;
	`

	// tests is a list of output expectations.
	tests := []token.Token{
		token.LET,
		{Type: "IDENTIFIER", Literal: "five"},
		token.ASSIGN,
		{Type: "INT", Literal: "5"},
		token.SEMICOLON,
		token.LET,
		{Type: "IDENTIFIER", Literal: "ten"},
		token.ASSIGN,
		{Type: "INT", Literal: "10"},
		token.SEMICOLON,
		token.LET,
		{Type: "IDENTIFIER", Literal: "add"},
		token.ASSIGN,
		token.FUNCTION,
		token.LPAREN,
		{Type: "IDENTIFIER", Literal: "x"},
		token.COMMA,
		{Type: "IDENTIFIER", Literal: "y"},
		token.RPAREN,
		token.LBRACE,
		{Type: "IDENTIFIER", Literal: "x"},
		token.PLUS,
		{Type: "IDENTIFIER", Literal: "y"},
		token.SEMICOLON,
		token.RBRACE,
		token.SEMICOLON,
		token.LET,
		{Type: "IDENTIFIER", Literal: "result"},
		token.ASSIGN,
		{Type: "IDENTIFIER", Literal: "add"},
		token.LPAREN,
		{Type: "IDENTIFIER", Literal: "five"},
		token.COMMA,
		{Type: "IDENTIFIER", Literal: "ten"},
		token.RPAREN,
		token.SEMICOLON,
		token.BANG,
		token.MINUS,
		token.SLASH,
		token.ASTERISK,
		{Type: "INT", Literal: "5"},
		token.SEMICOLON,
		{Type: "INT", Literal: "5"},
		token.LT,
		{Type: "INT", Literal: "10"},
		token.GT,
		{Type: "INT", Literal: "5"},
		token.SEMICOLON,
		token.IF,
		token.LPAREN,
		{Type: "INT", Literal: "5"},
		token.LT,
		{Type: "INT", Literal: "10"},
		token.RPAREN,
		token.LBRACE,
		token.RETURN,
		token.TRUE,
		token.SEMICOLON,
		token.RBRACE,
		token.ELSE,
		token.LBRACE,
		token.RETURN,
		token.FALSE,
		token.SEMICOLON,
		token.RBRACE,
		{Type: "INT", Literal: "5"},
		token.EQ,
		{Type: "INT", Literal: "5"},
		token.SEMICOLON,
		{Type: "INT", Literal: "5"},
		token.NOTEQ,
		{Type: "INT", Literal: "10"},
		token.SEMICOLON,
		token.EOF,
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

func TestNextTokenSingleLineInput(t *testing.T) {
	input := `
	let five = 5;
	`

	// tests is a list of output expectations.
	tests := []token.Token{
		token.LET,
		{Type: "IDENTIFIER", Literal: "five"},
		token.ASSIGN,
		{Type: "INT", Literal: "5"},
		token.SEMICOLON,
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
