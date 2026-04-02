package parser

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  token
	}{
		"skip whitespace": {
			{
				input: "\t\t   \r\n\t",
				want: token{
					Type:    t_NEWLINE,
					Literal: "\\n",
				},
			},
		},
		"item": {
			{
				input: "@ item-name",
				want: token{
					Type:    t_ITEM,
					Literal: "item-name",
				},
			},
		},
		"end": {
			{
				input: "",
				want: token{
					Type:    t_EOF,
					Literal: "",
				},
			},
		},
		"level": {
			{
				input: ": 69",
				want: token{
					Type:    t_IDENT,
					Literal: "69",
				},
			},
		},
		"nature": {
			{
				input: "Jolly Nature\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Jolly Nature",
				},
			},
		},
		"level token": {
			{
				input: "Level: 30",
				want: token{
					Type:    t_LEVEL,
					Literal: "",
				},
			},
		},
		"status token": {
			{
				input: "   Status>",
				want: token{
					Type:    t_STATUS,
					Literal: "",
				},
			},
		},
		"move": {
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

func TestReadIdentifier(t *testing.T) {
	tests := map[string][]struct {
		input string
		want  token
	}{
		"mr mime": {
			{
				input: "Mr. Mime   \n",
				want: token{
					Type:    t_IDENT,
					Literal: "Mr Mime",
				},
			},
		},
		"nidoran-f": {
			{
				input: "Nidoran-F\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Nidoran-F",
				},
			},
		},
		"farfetch": {
			{
				input: "   Farfetc'h\n",
				want: token{
					Type:    t_IDENT,
					Literal: "Farfetch",
				},
			},
		},
		"porygon2": {
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
