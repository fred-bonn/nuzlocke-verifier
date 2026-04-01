package parser

import "fmt"

type LexError struct {
	Path   string
	Line   int
	Column int
	Char   byte
}

type ParseError struct {
	Previous Token
	Current  Token
}

func (e LexError) Error() string {
	return fmt.Sprintf("error: lexer failed at %d:%d in '%s' (illegal character: %q)", e.Line, e.Column, e.Path, e.Char)
}

func (e ParseError) Error() string {
	return fmt.Sprintf("error: parser failed: did not expect token %s after %s", e.Current, e.Previous)
}
