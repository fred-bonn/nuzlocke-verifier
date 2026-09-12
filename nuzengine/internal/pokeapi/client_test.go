package pokeapi

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClientFetchPokemonHappyPath(t *testing.T) {
	client := Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://pokeapi.co/api/v2/pokemon/charizard" {
				t.Fatalf("unexpected URL: %s", req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"id":6,"name":"charizard","types":[{"type":{"name":"fire"}}],"stats":[{"base_stat":84,"stat":{"name":"hp"}}]}`)),
				Header:     make(http.Header),
			}, nil
		})},
		baseURL: "https://pokeapi.co/api/v2",
	}

	pokemon, err := client.FetchPokemon("charizard")
	if err != nil {
		t.Fatalf("FetchPokemon returned unexpected error: %v", err)
	}
	if pokemon.Name != "charizard" {
		t.Fatalf("unexpected pokemon name: %s", pokemon.Name)
	}
	if pokemon.Types[0].Type.Name != "fire" {
		t.Fatalf("unexpected primary type: %s", pokemon.Types[0].Type.Name)
	}
}

func TestClientFetchMoveHappyPath(t *testing.T) {
	client := Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != "https://pokeapi.co/api/v2/move/flamethrower" {
				t.Fatalf("unexpected URL: %s", req.URL.String())
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader(`{"name":"flamethrower","power":90,"pp":15,"accuracy":100,"type":{"name":"fire"},"damage_class":{"name":"special"},"meta":{"crit_rate":0,"category":{"name":"damage"}}}`)),
				Header:     make(http.Header),
			}, nil
		})},
		baseURL: "https://pokeapi.co/api/v2",
	}

	move, err := client.FetchMove("flamethrower")
	if err != nil {
		t.Fatalf("FetchMove returned unexpected error: %v", err)
	}
	if move.Name != "flamethrower" {
		t.Fatalf("unexpected move name: %s", move.Name)
	}
	if move.Power != 90 {
		t.Fatalf("unexpected move power: %d", move.Power)
	}
}

func TestClientFetchPokemonHandlesTransportErrors(t *testing.T) {
	client := Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return nil, fmt.Errorf("boom")
		})},
		baseURL: "https://pokeapi.co/api/v2",
	}

	if _, err := client.FetchPokemon("missingno"); err == nil {
		t.Fatal("expected transport error")
	}
}

func TestClientFetchMoveHandlesHTTPFailureStatus(t *testing.T) {
	client := Client{
		httpClient: &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Status:     "404 Not Found",
				Body:       io.NopCloser(strings.NewReader(`{"detail":"Not found"}`)),
				Header:     make(http.Header),
			}, nil
		})},
		baseURL: "https://pokeapi.co/api/v2",
	}

	if _, err := client.FetchMove("missing-move"); err == nil {
		t.Fatal("expected HTTP 404 error")
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
