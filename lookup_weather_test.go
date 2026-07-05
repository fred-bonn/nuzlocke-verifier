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
		ability     ability
		goggles     bool
		pokemonType pokemonType
		want        bool
	}{
		"overcoat sand":    {sandstormWeather, overcoatAbility, false, normalType, false},
		"overcoat hail":    {hailWeather, overcoatAbility, false, normalType, false},
		"magic guard sand": {sandstormWeather, magicGuardAbility, false, normalType, false},
		"magic guard hail": {hailWeather, magicGuardAbility, false, normalType, false},
		"intim":            {sandstormWeather, intimidateAbility, false, normalType, true},
		"intim goggles":    {sandstormWeather, intimidateAbility, true, normalType, false},
		"sand rush 1":      {sandstormWeather, sandRushAbility, false, normalType, false},
		"sand rush 2":      {hailWeather, sandRushAbility, false, normalType, true},
		"sand rush 3":      {hailWeather, sandRushAbility, true, normalType, false},
		"ice body 1":       {hailWeather, iceBodyAbility, false, normalType, false},
		"ice body 2":       {sandstormWeather, iceBodyAbility, false, normalType, true},
		"ice body 3":       {sandstormWeather, iceBodyAbility, true, normalType, false},
		"ground sand":      {sandstormWeather, intimidateAbility, false, groundType, false},
		"steel sand":       {sandstormWeather, intimidateAbility, false, steelType, false},
		"rock sand":        {sandstormWeather, intimidateAbility, false, rockType, false},
		"ice sand":         {sandstormWeather, intimidateAbility, false, iceType, true},
		"ice hail":         {hailWeather, intimidateAbility, false, iceType, false},
		"ground hail":      {hailWeather, intimidateAbility, false, groundType, true},
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
		"rain":      {rainWeather, "to rain"},
		"sun":       {sunWeather, "turned harsh"},
		"sandstorm": {sandstormWeather, "brewed"},
		"hail":      {hailWeather, "to hail"},
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
