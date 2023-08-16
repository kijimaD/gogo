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
	SEMICOLON = ";"
)

const (
	CTYPE_VOID = "void"
	CTYPE_INT  = "int"
	CTYPE_CHAR = "char"
	CTYPE_STR  = "string"
)
