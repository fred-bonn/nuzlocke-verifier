package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fred-bonn/nuz/internal/engine"
	"github.com/fred-bonn/nuz/internal/parser"
	"github.com/fred-bonn/nuz/internal/pokeapi"
)

func TestConfigValidateInputLoadsShowdownPartyFromLocalData(t *testing.T) {
	cfg := &config{client: pokeapi.NewClient()}

	playerData, err := os.ReadFile("showdown_demo_files/player.txt")
	if err != nil {
		t.Fatalf("failed reading player party file: %v", err)
	}
	party, err := cfg.validateInput(string(playerData))
	if err != nil {
		t.Fatalf("validateInput returned unexpected error: %v", err)
	}
	if len(party) != 5 {
		t.Fatalf("unexpected party size: got %d want 5", len(party))
	}
	if party[0].Base.Name != "horsea" {
		t.Fatalf("unexpected first pokemon name: %q", party[0].Base.Name)
	}
	if len(party[0].Moves) != 2 {
		t.Fatalf("unexpected move count for horsea: got %d want 2", len(party[0].Moves))
	}
	if party[0].Moves[0].Name != "bubble beam" {
		t.Fatalf("unexpected first move name: %q", party[0].Moves[0].Name)
	}
}

func TestConfigNameHelpersAndHiddenPower(t *testing.T) {
	if got := apiName("Mr. Mime"); got != "mr-mime" {
		t.Fatalf("apiName returned %q, want %q", got, "mr-mime")
	}
	if got := cleanName("Mr. Mime"); got != "mr. mime" {
		t.Fatalf("cleanName returned %q, want %q", got, "mr. mime")
	}

	move, err := generateHiddenPower("hidden-power-fire")
	if err != nil {
		t.Fatalf("generateHiddenPower returned unexpected error: %v", err)
	}
	if move.Name != "hidden power" || move.Type != engine.FireType || move.Power != 60 || move.Accuracy != 100 {
		t.Fatalf("unexpected generated hidden power: %+v", move)
	}

	if _, err := generateHiddenPower("hidden-power-unknown"); err == nil {
		t.Fatal("generateHiddenPower should reject invalid hidden power types")
	}
}

func TestParserReadShowdownFileRejectsMalformedInput(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.txt")
	content := "Pikachu\nJolly Nature\nAbility: Static\n- Tackle\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	if _, err := parser.ReadShowdownFile(path); err == nil {
		t.Fatal("ReadShowdownFile should reject malformed party data")
	}
}

func TestCleanNameFormatsPokemonNamesConsistently(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases mixed-case names":    {input: "pIkAcHu", want: "pikachu"},
		"preserves punctuation in names": {input: "Mr. Mime", want: "mr. mime"},
		"preserves apostrophes in names": {input: "Farfetch’d", want: "farfetch’d"},
		"handles empty input":            {input: "", want: ""},
		"keeps numerals in names":        {input: "Porygon2", want: "porygon2"},
		"normalizes hyphenated names":    {input: "Ho-Oh", want: "ho-oh"},
		"normalizes regional forms":      {input: "Arcanine-Hisui", want: "arcanine-hisui"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanPokemonName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestHasHyphenDetectsHyphenatedNames(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"detects no hyphen in empty input":      {input: "", want: false},
		"detects a hyphen in a hyphenated name": {input: "ho-oh", want: true},
		"detects no hyphen in a simple name":    {input: "pikachu", want: false},
		"detects no hyphen in a regional form":  {input: "arcanine-hisui", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := hasHyphen(tc.input); got != tc.want {
				t.Errorf("%s: hasHyphen(%q) = %t, want %t", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestIsRegionalPokemonDetectsRegionalForms(t *testing.T) {
	tests := map[string]struct {
		input string
		want  bool
	}{
		"detects no regional form in empty input":       {input: "", want: false},
		"detects no regional form in a hyphenated name": {input: "ho-oh", want: false},
		"detects no regional form in a simple name":     {input: "pikachu", want: false},
		"detects a regional form in a regional name":    {input: "arcanine-hisui", want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := isRegionalPokemon(tc.input); got != tc.want {
				t.Errorf("%s: isRegionalPokemon(%q) = %t, want %t", name, tc.input, got, tc.want)
			}
		})
	}
}

func TestCleanNameFormatsMoveNamesConsistently(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"lowercases mixed-case move names": {input: "ThUndEr SHoCk", want: "thunder shock"},
		"handles empty move input":         {input: "", want: ""},
		"normalizes hyphenated move names": {input: "Tri-Attack", want: "tri attack"},
		"preserves numerals in move names": {input: "conversion 2", want: "conversion 2"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := cleanName(tc.input); got != tc.want {
				t.Errorf("%s: cleanMoveName(%q) = %q, want %q", name, tc.input, got, tc.want)
			}
		})
	}
}
