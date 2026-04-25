package main

import "testing"

func TestCleanPokemonNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases": {"pIkAcHu", "pikachu"},
		"dot":        {"Mr. Mime", "mr-mime"},
		"Farfetch’d": {"Farfetch’d", "farfetchd"},
		"empty":      {"", ""},
		"numerals":   {"Porygon2", "porygon2"},
		"dash":       {"Ho-Oh", "ho-oh"},
		"regional":   {"Arcanine-Hisui", "arcanine-hisui"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanPokemonName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestCleanMoveNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases": {"ThUndEr SHoCk", "thunder-shock"},
		"empty":      {"", ""},
		"dash":       {"Tri-Attack", "tri-attack"},
		"numerals":   {"conversion 2", "conversion-2"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanMoveName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}
