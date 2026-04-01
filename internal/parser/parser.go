package parser

import (
	"fmt"
	"os"
)

func Parse() {

	input, err := os.ReadFile("./internal/parser/example.txt")
	if err != nil {
		fmt.Println("failed")
		return
	}

	l := New(string(input))

	for tok := l.NextToken(); tok.Type != EOF; tok = l.NextToken() {
		fmt.Printf("%s -> %q\n", tok.Type, tok.Literal)
	}
}
