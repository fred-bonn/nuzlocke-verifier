package main

import "testing"

func TestCleanPokemonNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases": {input: "pIkAcHu", want: "pikachu"},
		"dot":        {input: "Mr. Mime", want: "mr. mime"},
		"Farfetch’d": {input: "Farfetch’d", want: "farfetch’d"},
		"empty":      {input: "", want: ""},
		"numerals":   {input: "Porygon2", want: "porygon2"},
		"hyphen":     {input: "Ho-Oh", want: "ho-oh"},
		"regional":   {input: "Arcanine-Hisui", want: "arcanine-hisui"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanPokemonName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestHasHyphen(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"empty": {input: "", want: false},
		"mon":   {input: "ho-oh", want: true},
		"mon2":  {input: "pikachu", want: false},
		"mon3":  {input: "arcanine-hisui", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := hasHyphen(tc.input); got != tc.want {
				t.Errorf("%s: hasHyphen(%q) = %t, want %t", name, tc.input, got, tc.want)
			}
		})
	}
}

func IsRegionalPokemon(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"empty": {input: "", want: false},
		"mon":   {input: "ho-oh", want: false},
		"mon2":  {input: "pikachu", want: false},
		"mon3":  {input: "arcanine-hisui", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := hasHyphen(name); got != tc.want {
				t.Errorf("%s: isRegionalPokemon(%q) = %t, want %t", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestCleanMoveNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases": {input: "ThUndEr SHoCk", want: "thunder shock"},
		"empty":      {input: "", want: ""},
		"dash":       {input: "Tri-Attack", want: "tri attack"},
		"numerals":   {input: "conversion 2", want: "conversion 2"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanMoveName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}
