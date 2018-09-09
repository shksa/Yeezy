package token

// Token is data-structure that represents tokens of the language.
type Token struct {
	Type    string
	Literal string
}

// List of all tokens in the language.
var (
	// Operators
	ASSIGN   = Token{"ASSIGN", "="}
	PLUS     = Token{"PLUS", "+"}
	MINUS    = Token{"MINUS", "-"}
	BANG     = Token{"BANG", "!"}
	ASTERISK = Token{"ASTERISK", "*"}
	SLASH    = Token{"SLASH", "/"}

	LT = Token{"LT", "<"}
	GT = Token{"GT", ">"}

	EQ    = Token{"EQ", "=="}
	NOTEQ = Token{"NOTEQ", "!="}

	// Delimiters
	COMMA     = Token{"COMMA", ","}
	SEMICOLON = Token{"SEMICOLAN", ";"}

	// Brackets
	LPAREN = Token{"LPAREN", "("}
	RPAREN = Token{"RPAREN", ")"}
	LBRACE = Token{"LBRACE", "{"}
	RBRACE = Token{"RBRACE", "}"}

	// Keywords
	FUNCTION = Token{"FUNCTION", "func"}
	LET      = Token{"LET", "let"}
	IF       = Token{"IF", "if"}
	ELSE     = Token{"ELSE", "else"}
	RETURN   = Token{"RETURN", "return"}
	TRUE     = Token{"TRUE", "true"}
	FALSE    = Token{"FALSE", "false"}

	// Identifiers + Literals
	IDENTIFIER = Token{Type: "IDENTIFIER"} // add, foobar, x, y, ...
	INT        = Token{Type: "INT"}        // 23, 4343, 989898

	// Special tokens
	ILLEGAL = Token{Type: "ILLEGAL"}
	EOF     = Token{"EOF", ""}
)

// keywords table maps all the keyword token literals to their token values
var keywords = map[string]Token{
	"func":   FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

// GetTokenForLetterStringLiteral returns token for a letter-string literal.
func GetTokenForLetterStringLiteral(literal string) Token {
	if keywordToken, ok := keywords[literal]; ok {
		return keywordToken
	}
	identifierToken := IDENTIFIER
	identifierToken.Literal = literal
	return identifierToken
}
