package parser

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENT   TokenType = "IDENT"
	ITEM    TokenType = "ITEM"
	MOVE    TokenType = "MOVE"
	NEWLINE TokenType = "NEWLINE"
)

type Token struct {
	Type    TokenType
	Literal string
}
