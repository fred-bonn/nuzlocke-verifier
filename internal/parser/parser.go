package parser

import (
	"fmt"
	"os"
)

func Parse() ([]token, error) {
	path := "./internal/parser/example.txt"

	input, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %s", "<path>")
	}

	var res []token
	l := newLexer(string(input))

	for tok := l.nextToken(); tok.Type != t_EOF; tok = l.nextToken() {
		if tok.Type == t_ILLEGAL {
			return nil, lexError{
				Line:   l.line,
				Column: l.column,
				Char:   l.ch,
			}
		}
		res = append(res, tok)
	}

	// verify the sequence of tokens is correct here

	return res, nil
}
