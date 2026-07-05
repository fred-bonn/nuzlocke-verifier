package main

import "testing"

func TestGetEffectiveness(t *testing.T) {
	tests := map[string]struct {
		attacking pokemonType
		defending pokemonType
		want      effectiveness
	}{
		"normal/ghost":     {normalType, ghostType, immuneEffectivensss},
		"ghost/normal":     {ghostType, normalType, immuneEffectivensss},
		"fighting/normal":  {fightingType, normalType, superEffectiveness},
		"psychic/fighting": {psychicType, fightingType, superEffectiveness},
		"steel/psychic":    {steelType, psychicType, normalEffectivenss},
		"bug/steel":        {bugType, steelType, resistedEffectiveness},
		"ground/bug":       {groundType, bugType, resistedEffectiveness},
		"water/ground":     {waterType, groundType, superEffectiveness},
		"ice/water":        {iceType, waterType, resistedEffectiveness},
		"fire/ice":         {fireType, iceType, superEffectiveness},
		"fairy/fire":       {fairyType, fireType, resistedEffectiveness},
		"dragon/fairy":     {dragonType, fairyType, immuneEffectivensss},
		"electric/dragon":  {electricType, dragonType, resistedEffectiveness},
		"grass/electric":   {grassType, electricType, normalEffectivenss},
		"rock/grass":       {rockType, grassType, normalEffectivenss},
		"flying/rock":      {flyingType, rockType, resistedEffectiveness},
		"dark/flying":      {darkType, flyingType, normalEffectivenss},
		"poison/dark":      {poisonType, darkType, normalEffectivenss},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := getEffectiveness(tc.attacking, tc.defending); got != tc.want {
				t.Errorf("%s: getEffectiveness(%s, %s) = %.1f, want %.1f", name, tc.attacking.String(), tc.defending.String(), got, tc.want)
			}
		})
	}
}
