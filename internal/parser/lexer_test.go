package parser

import (
	"testing"
)

func TestNextTokenReadsTheExpectedTokenType(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  token
	}{
		"skips leading whitespace": {
			{
				input: "\t\t   \r\n\t",
				want: token{
					Type:    t_NEWLINE,
					Literal: "\\n",
				},
			},
		},
		"parses an item token": {
			{
				input: "@ item-name",
				want: token{
					Type:    t_ITEM,
					Literal: "item-name",
				},
			},
		},
		"returns EOF at the end of input": {
			{
				input: "",
				want: token{
					Type:    t_EOF,
					Literal: "",
				},
			},
		},
		"parses a level token": {
			{
				input: ": 69",
				want: token{
					Type:    t_IDENT,
					Literal: "69",
				},
			},
		},
		"parses a nature token": {
			{
				input: "Jolly Nature\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Jolly Nature",
				},
			},
		},
		"recognizes a level token": {
			{
				input: "Level: 30",
				want: token{
					Type:    t_LEVEL,
					Literal: "",
				},
			},
		},
		"recognizes a status token": {
			{
				input: "   Status>",
				want: token{
					Type:    t_STATUS,
					Literal: "",
				},
			},
		},
		"parses a move token": {
			{
				input: "\t - Tackle",
				want: token{
					Type:    t_MOVE,
					Literal: "-",
				},
			},
		},
	}

	for name, tcs := range tests {
		for _, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				l := newLexer(string(tc.input))

				if got := l.nextToken(); got != tc.want {
					t.Errorf("%s: l.NextToken() = %q, want %q", name, got, tc.want)
				}
			})
		}
	}
}

func TestReadIdentifierReadsPokemonNamesCorrectly(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  token
	}{
		"parses a name with a period": {
			{
				input: "Mr. Mime   \n",
				want: token{
					Type:    t_IDENT,
					Literal: "Mr Mime",
				},
			},
		},
		"parses a hyphenated name": {
			{
				input: "Nidoran-F\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Nidoran-F",
				},
			},
		},
		"parses a name with an apostrophe": {
			{
				input: "   Farfetc'h\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Farfetch",
				},
			},
		},
		"parses a name with a numeral": {
			{
				input: "Porygon2\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Porygon2",
				},
			},
		},
	}

	for name, tcs := range tests {
		for _, tc := range tcs {
			t.Run(name, func(t *testing.T) {
				l := newLexer(string(tc.input))

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
