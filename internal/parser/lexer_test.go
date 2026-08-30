package parser

import (
	"os"
	"path/filepath"
	"strings"
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

func TestReadShowdownFileRejectsMalformedNatureLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad_nature.txt")
	content := "Pikachu\nLevel: 25\nJolly\nAbility: Static\n- Thunderbolt\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp showdown file: %v", err)
	}

	_, err := ReadShowdownFile(path)
	if err == nil {
		t.Fatal("ReadShowdownFile() = nil error, want malformed nature error")
	}
	if !strings.Contains(err.Error(), "nature line not enough fields") {
		t.Fatalf("ReadShowdownFile() error = %v, want nature format error", err)
	}
}

func TestReadShowdownFileRejectsUnknownStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad_status.txt")
	content := "Pikachu\nLevel: 25\nJolly Nature\nAbility: Static\nStatus: Confused\n- Thunderbolt\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp showdown file: %v", err)
	}

	_, err := ReadShowdownFile(path)
	if err == nil {
		t.Fatal("ReadShowdownFile() = nil error, want invalid status error")
	}
	if !strings.Contains(err.Error(), "unrecognized status") {
		t.Fatalf("ReadShowdownFile() error = %v, want status validation error", err)
	}
}

func TestReadShowdownFileAllowsMinimalPokemonWithoutStatusOrHP(t *testing.T) {
	path := filepath.Join(t.TempDir(), "minimal.txt")
	content := "Pikachu\nLevel: 25\nJolly Nature\nAbility: Static\n- Thunderbolt\n- Quick Attack\n- Tail Whip\n- Slam\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp showdown file: %v", err)
	}

	pokemon, err := ReadShowdownFile(path)
	if err != nil {
		t.Fatalf("ReadShowdownFile() returned error for minimal valid team entry: %v", err)
	}
	if len(pokemon) != 1 {
		t.Fatalf("ReadShowdownFile() = %d entries, want 1", len(pokemon))
	}
	if pokemon[0].HP != -1 {
		t.Fatalf("pokemon[0].HP = %d, want -1 when no HP line is present", pokemon[0].HP)
	}
	if pokemon[0].Status != "" {
		t.Fatalf("pokemon[0].Status = %q, want empty string when no status line is present", pokemon[0].Status)
	}
	if len(pokemon[0].Moves) != 4 {
		t.Fatalf("len(pokemon[0].Moves) = %d, want 4", len(pokemon[0].Moves))
	}
}
