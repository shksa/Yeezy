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
	lexer.readNextChar() // To initialize lexer.ch, lexer.postion, lexer.nextPosition
	return lexer
}

// readNextChar reads the next char in the input string and stores it in the lexer's current char (ch) field.
func (l *Lexer) readNextChar() {
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
		tok = token.ASSIGN
	case '+':
		tok = token.PLUS
	case '-':
		tok = token.MINUS
	case '*':
		tok = token.ASTERISK
	case '/':
		tok = token.SLASH
	case '!':
		tok = token.BANG
	case '(':
		tok = token.LPAREN
	case ')':
		tok = token.RPAREN
	case '{':
		tok = token.LBRACE
	case '}':
		tok = token.RBRACE
	case ',':
		tok = token.COMMA
	case ';':
		tok = token.SEMICOLON
	case '<':
		tok = token.LT
	case '>':
		tok = token.GT
	case 0:
		tok = token.EOF
	default:
		if isLetter(l.ch) {
			letterStringLiteral := l.readLetterString()
			tok = token.GetTokenForLetterStringLiteral(letterStringLiteral)
			return tok
		} else if isDigit(l.ch) {
			tok = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = token.ILLEGAL
			tok.Literal = string(l.ch)
		}
	}
	l.readNextChar()
	return tok
}

// readLetterString returns an string of letters from input source code
func (l *Lexer) readLetterString() string {
	position := l.position
	for isLetter(l.ch) {
		l.readNextChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readNextChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readNextChar()
	}
}

// isLetter determines what characters can be used in identifiers and keywords
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
