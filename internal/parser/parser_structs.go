package parser

import "fmt"

type parser struct {
	tokens []token
	pos    int
}

type parseError struct {
	Previous token
	Current  token
}

func (e parseError) Error() string {
	return fmt.Sprintf("error: parser failed: did not expect token %s after %s", e.Current, e.Previous)
}

func newParser(tokens []token) *parser {
	p := &parser{
		tokens: tokens,
		pos:    0,
	}

	return p
}
