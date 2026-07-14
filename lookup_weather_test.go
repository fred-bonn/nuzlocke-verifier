package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestWeatherAffectsMon(t *testing.T) {
	tests := map[string]struct {
		weather     weatherState
		ability     abilityState
		goggles     bool
		pokemonType pokemonType
		want        bool
	}{
		"overcoat sand":    {weather: sandstormWeather, ability: overcoatAbility, goggles: false, pokemonType: normalType, want: false},
		"overcoat hail":    {weather: hailWeather, ability: overcoatAbility, goggles: false, pokemonType: normalType, want: false},
		"magic guard sand": {weather: sandstormWeather, ability: magicGuardAbility, goggles: false, pokemonType: normalType, want: false},
		"magic guard hail": {weather: hailWeather, ability: magicGuardAbility, goggles: false, pokemonType: normalType, want: false},
		"intim":            {weather: sandstormWeather, ability: intimidateAbility, goggles: false, pokemonType: normalType, want: true},
		"intim goggles":    {weather: sandstormWeather, ability: intimidateAbility, goggles: true, pokemonType: normalType, want: false},
		"sand rush 1":      {weather: sandstormWeather, ability: sandRushAbility, goggles: false, pokemonType: normalType, want: false},
		"sand rush 2":      {weather: hailWeather, ability: sandRushAbility, goggles: false, pokemonType: normalType, want: true},
		"sand rush 3":      {weather: hailWeather, ability: sandRushAbility, goggles: true, pokemonType: normalType, want: false},
		"ice body 1":       {weather: hailWeather, ability: iceBodyAbility, goggles: false, pokemonType: normalType, want: false},
		"ice body 2":       {weather: sandstormWeather, ability: iceBodyAbility, goggles: false, pokemonType: normalType, want: true},
		"ice body 3":       {weather: sandstormWeather, ability: iceBodyAbility, goggles: true, pokemonType: normalType, want: false},
		"ground sand":      {weather: sandstormWeather, ability: intimidateAbility, goggles: false, pokemonType: groundType, want: false},
		"steel sand":       {weather: sandstormWeather, ability: intimidateAbility, goggles: false, pokemonType: steelType, want: false},
		"rock sand":        {weather: sandstormWeather, ability: intimidateAbility, goggles: false, pokemonType: rockType, want: false},
		"ice sand":         {weather: sandstormWeather, ability: intimidateAbility, goggles: false, pokemonType: iceType, want: true},
		"ice hail":         {weather: hailWeather, ability: intimidateAbility, goggles: false, pokemonType: iceType, want: false},
		"ground hail":      {weather: hailWeather, ability: intimidateAbility, goggles: false, pokemonType: groundType, want: true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dummy := &pokemon{
				base: BasePokemon{
					Types: []pokemonType{
						tc.pokemonType,
					},
				},
				ability: tc.ability,
			}

			if tc.goggles {
				dummy.item = &item{
					state: safetyGoggles,
				}
			} else {
				dummy.item = &item{
					state:    noneItem,
					consumed: true,
				}
			}

			if got := tc.weather.affectsMon(dummy); got != tc.want {
				t.Errorf("%s: tc.weather.affectsMon(%q) = %t, want %t", name, tc.ability, got, tc.want)
			}
		})
	}
}

func TestWeatherOnset(t *testing.T) {
	tests := map[string]struct {
		weather  weatherState
		contains string
	}{
		"rain":      {weather: rainWeather, contains: "to rain"},
		"sun":       {weather: sunWeather, contains: "turned harsh"},
		"sandstorm": {weather: sandstormWeather, contains: "brewed"},
		"hail":      {weather: hailWeather, contains: "to hail"},
	}

	oldVerbose := *verbose
	*verbose = true
	t.Cleanup(func() {
		*verbose = oldVerbose
	})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var bf bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			t.Cleanup(func() {
				os.Stdout = oldStdout
			})

			tc.weather.onset()
			w.Close()
			io.Copy(&bf, r)

			if !strings.Contains(bf.String(), tc.contains) {
				t.Errorf("%s: tc.weather.onset() did not contain \"%s\" in standard output: %s", name, tc.contains, bf.String())
			}
		})
	}
}
