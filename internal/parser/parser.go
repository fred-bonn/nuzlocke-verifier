package parser

import (
	"fmt"
	"os"
)

func Parse() ([]Token, error) {
	path := "./internal/parser/example.txt"

	input, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %s", "<path>")
	}

	var res []Token
	l := New(string(input))

	for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		if tok.Type == ILLEGAL {
			return nil, LexError{
				Path:   path,
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
