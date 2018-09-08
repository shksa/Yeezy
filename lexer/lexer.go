package lexer

import "github.com/shksa/monkey/token"

// Lexer is the object which generates tokens from source code.
type Lexer struct {
	input        string
	position     int  // points to current char
	nextPosition int  // points to next char
	ch           byte // current char under examination
}

/* NOTES
The reason for these two "pointers" pointing into our input string is the fact that we will need
to "peek" further into the input and look after the current character to see what comes up next.
1. "nextPosition" always points to the "next" character in the input.
2. "position" points to the character in the input that corresponds to the ch byte.
*/

// New returns a pointer to a newly created Lexer object.
func New(input string) *Lexer {
	lexer := &Lexer{input: input}
	return lexer
}

// The purpose of advance is to advance our position in the input string.
func (l *Lexer) advance() {
	if l.nextPosition > len(l.input) {
		l.ch = 0 // 0 is the ASCII code for the "NUL" character and signifies either "we haven't read anything yet" or "end of file" for us.
	} else {
		l.ch = l.input[l.nextPosition]
		l.position = l.nextPosition
		l.nextPosition++ // l.nextPosition always points to the position where we're going to read from next.
	}
}

/*
NextToken returns the next token of the source code.
It advances lexer's position from the previous char to the next char.
Looks at the current char and returns a token depending on what char it is.
*/
func (l *Lexer) NextToken() token.Token {
	l.advance()
	var tok token.Token
	switch l.ch {
	case '=':
		tok = token.Token{Type: token.ASSIGN, Literal: string(l.ch)}
	case '+':
		tok = token.Token{Type: token.PLUS, Literal: string(l.ch)}
	case '(':
		tok = token.Token{Type: token.LPAREN, Literal: string(l.ch)}
	case ')':
		tok = token.Token{Type: token.RPAREN, Literal: string(l.ch)}
	case '{':
		tok = token.Token{Type: token.LBRACE, Literal: string(l.ch)}
	case '}':
		tok = token.Token{Type: token.RBRACE, Literal: string(l.ch)}
	case ',':
		tok = token.Token{Type: token.COMMA, Literal: string(l.ch)}
	case ';':
		tok = token.Token{Type: token.SEMICOLAN, Literal: string(l.ch)}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	}
	return tok
}
