package parser

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  Token
	}{
		"skip whitespace": {
			{
				input: "\t\t   \r\n\t",
				want: Token{
					Type:    NEWLINE,
					Literal: "\\n",
				},
			},
		},
		"item": {
			{
				input: "@ item-name",
				want: Token{
					Type:    ITEM,
					Literal: "item-name",
				},
			},
		},
		"end": {
			{
				input: "",
				want: Token{
					Type:    EOF,
					Literal: "",
				},
			},
		},
		"level": {
			{
				input: ": 69",
				want: Token{
					Type:    IDENT,
					Literal: "69",
				},
			},
		},
		"nature": {
			{
				input: "Jolly Nature\n",
				want: Token{
					Type:    IDENT,
					Literal: "Jolly Nature",
				},
			},
		},
		"level token": {
			{
				input: "Level: 30",
				want: Token{
					Type:    LEVEL,
					Literal: "Level",
				},
			},
		},
	}

	for name, tcs := range tests {
		for _, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				l := New(string(tc.input))

				if got := l.NextToken(); got != tc.want {
					t.Errorf("%s: l.NextToken() = %q, want %q", name, got, tc.want)
				}
			})
		}
	}
}

func TestReadIdentifier(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  Token
	}{
		"mr mime": {
			{
				input: "Mr. Mime   \n",
				want: Token{
					Type:    IDENT,
					Literal: "Mr Mime",
				},
			},
		},
		"nidoran-f": {
			{
				input: "Nidoran-F\n",
				want: Token{
					Type:    IDENT,
					Literal: "Nidoran-F",
				},
			},
		},
		"farfetch": {
			{
				input: "   Farfetc'h\n",
				want: Token{
					Type:    IDENT,
					Literal: "Farfetch",
				},
			},
		},
		"porygon2": {
			{
				input: "Porygon2\n",
				want: Token{
					Type:    IDENT,
					Literal: "Porygon2",
				},
			},
		},
	}

	for name, tcs := range tests {
		for _, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				l := New(string(tc.input))

				if got := l.readIdent(); got != tc.want {
					t.Errorf("%s: parsePokemonLine(%q) = %q, want %q", name, tc.input, got, tc.want)
				}
				if l.ch != '\n' {
					t.Errorf("%s: l.ch = '%c', want '\n'", name, l.ch)
				}
			})
		}
	}
}
