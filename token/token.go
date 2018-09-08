package token

// Token is data-structure that represents tokens of the language.
type Token struct {
	Type    string
	Literal string
}

// List of all types of tokens in the language.
const (
	// Operators
	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"

	// Delimiters
	COMMA     = "COMMA"
	SEMICOLAN = "SEMICOLAN"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	// Keywords
	FUNCTION = "func"
	LET      = "let"

	// Identifiers + Literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 23, 4343, 989898

	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)
