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
	lexer.advance() // To initialize lexer.ch, lexer.postion, lexer.nextPosition
	return lexer
}

// The purpose of advance is to advance our position in the input string.
func (l *Lexer) advance() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0 // 0 is the ASCII code for the "NUL" character and signifies either "we haven't read anything yet" or "end of file" for us.
	} else {
		l.ch = l.input[l.nextPosition]
		l.position = l.nextPosition
		l.nextPosition++ // l.nextPosition always points to the position where we're going to read from next.
	}
}

/*
NextToken returns the next token from the source code.
It looks at the current char and returns a token depending on what char it is.
It advances the lexer's position before returning the token
*/
func (l *Lexer) NextToken() token.Token {
	l.skipWhiteSpace()
	var tok token.Token
	switch l.ch {
	case '=':
		tok.Type = token.ASSIGN
		tok.Literal = string(l.ch)
	case '+':
		tok.Type = token.PLUS
		tok.Literal = string(l.ch)
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = string(l.ch)
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = string(l.ch)
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = string(l.ch)
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = string(l.ch)
	case ',':
		tok.Type = token.COMMA
		tok.Literal = string(l.ch)
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = string(l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readLetterString()
			tok.Type = token.TypeOfLetterString(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	}
	l.advance()
	return tok
}

// readLetterString returns an string of letters by advancing the lexer's position
// if the current char is a letter untill a non-letter char is encountered.
func (l *Lexer) readLetterString() string {
	position := l.position
	for isLetter(l.ch) {
		l.advance()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.advance()
	}
	return l.input[position:l.position]
}

// isLetter determines what characters can be used in identifiers and keywords
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.advance()
	}
}
