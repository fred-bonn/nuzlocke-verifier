package parser

import "fmt"

type tokenType string

const (
	t_START   tokenType = "START"
	t_ILLEGAL tokenType = "ILLEGAL"
	t_EOF     tokenType = "EOF"

	t_IDENT   tokenType = "IDENT"
	t_ITEM    tokenType = "ITEM"
	t_LEVEL   tokenType = "LEVEL"
	t_ABILITY tokenType = "ABILITY"
	t_IVS     tokenType = "IVS"
	t_STATUS  tokenType = "STATUS"
	t_HP      tokenType = "HP"
	t_MOVE    tokenType = "MOVE"
	t_NEWLINE tokenType = "NEWLINE"
)

type token struct {
	Type    tokenType
	Literal string
}

func (t token) String() string {
	return fmt.Sprintf("(%s, %s)", t.Type, t.Literal)
}
