package main

import "testing"

func TestIsImmuneToPowderMoves(t *testing.T) {
	tests := map[string]struct {
		t    pokemonType
		a    ability
		want bool
	}{
		"grass":              {grassType, intimidateAbility, true},
		"grass overcoat:":    {grassType, overcoatAbility, true},
		"not grass 1":        {flyingType, intimidateAbility, false},
		"not grass 2":        {waterType, intimidateAbility, false},
		"not grass overcoat": {fireType, overcoatAbility, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				base: BasePokemon{
					Types: []pokemonType{
						tc.t,
					},
				},
				ability: tc.a,
			}

			if got := mon.isImmuneToPowderMoves(); got != tc.want {
				t.Errorf("mon.isImmuneToPowderMoves() = %t, want %t", got, tc.want)
			}
		})
	}
}
