package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"

	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	ASSIGN   = "="

	INT       = "INT"
	STRING    = "STRING"
	CHAR      = "CHAR"
	SEMICOLON = ";"
	COMMA     = ","
	LPAREN    = "("
	RPAREN    = ")"
)

type Ctype int

const (
	CTYPE_VOID Ctype = iota
	CTYPE_INT
	CTYPE_CHAR
	CTYPE_STR
)
