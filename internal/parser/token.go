package parser

import "fmt"

type TokenType string

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENT   TokenType = "IDENT"
	ITEM    TokenType = "ITEM"
	LEVEL   TokenType = "LEVEL"
	ABILITY TokenType = "ABILITY"
	IVS     TokenType = "IVS"
	STATUS  TokenType = "STATUS"
	HP      TokenType = "HP"
	MOVE    TokenType = "MOVE"
	NEWLINE TokenType = "NEWLINE"
)

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return fmt.Sprintf("(%s, %s)", t.Type, t.Literal)
}
