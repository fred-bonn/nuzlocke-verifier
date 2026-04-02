package parser

import "fmt"

type lexError struct {
	Line   int
	Column int
	Char   byte
}

func (e lexError) Error() string {
	return fmt.Sprintf("error: lexer failed at %d:%d (illegal character: %q)", e.Line, e.Column, e.Char)
}
