package parser

import "fmt"

type lexError struct {
	Line   int
	Column int
	Char   byte
}

type parseError struct {
	Previous token
	Current  token
}

func (e lexError) Error() string {
	return fmt.Sprintf("error: lexer failed at %d:%d (illegal character: %q)", e.Line, e.Column, e.Char)
}

func (e parseError) Error() string {
	return fmt.Sprintf("error: parser failed: did not expect token %s after %s", e.Current, e.Previous)
}
