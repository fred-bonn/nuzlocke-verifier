package main

import (
	"testing"
)

func TestParsePokemonLine(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"with item":    {"Pikachu @ Light Ball", "Pikachu"},
		"without item": {"Charizard", "Charizard"},
		"empty":        {"", ""},
		"extra spaces": {"Bulbasaur  @     Leftovers  ", "Bulbasaur"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parsePokemonLine(tc.input); got != tc.want {
				t.Errorf("%s: parsePokemonLine(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestParseLevelLine(t *testing.T) {
	tests := map[string]struct {
		input string
		want  int
	}{
		"valid":          {"Level: 25", 25},
		"with spaces":    {"  Level:  30  ", 30},
		"invalid format": {"Lvl: 20", 0},
		"non-numeric":    {"Level: Twenty", 0},
		"empty":          {"", 0},
		"missing space":  {"Level:50", 50},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parseLevelLine(tc.input); got != tc.want {
				t.Errorf("%s: parseLevelLine(%q) = %d, want %d", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestParseNatureLine(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"valid":          {"Modest Nature", "modest"},
		"wrong sep":      {"Adamant:Nature", ""},
		"invalid nature": {"hello Nature", ""},
		"extra spaces":   {"     Bashful   Nature    ", "bashful"},
		"empty":          {"", ""},
		"missing nature": {"bold", ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parseNatureLine(tc.input); got != tc.want {
				t.Errorf("%s: parseNatureLine(%q) = %s, want %s", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestParseMoveLine(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"valid":               {"- Tackle", "Tackle"},
		"wrong pre":           {"~ pound", ""},
		"extra space invalid": {"  -   ember ", ""},
		"extra space valid":   {"-   ember ", "ember"},
		"empty":               {"", ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := parseMoveLine(tc.input); got != tc.want {
				t.Errorf("%s: parseMoveLine(%q) = %s, want %s", name, tc.input, got, tc.want)
			}
		})
	}
}
