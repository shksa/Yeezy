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
	SEMICOLON = "SEMICOLAN"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"

	// Identifiers + Literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 23, 4343, 989898

	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

// keywords table maps keyword literals to their token types
var keywords = map[string]string{
	"func": FUNCTION,
	"let":  LET,
}

// TypeOfLetterString returns token type of the letter-string literal argument
func TypeOfLetterString(literal string) string {
	if tokenTypeOfKeyword, ok := keywords[literal]; ok {
		return tokenTypeOfKeyword
	}
	return IDENT
}
