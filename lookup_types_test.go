package main

import "testing"

func TestGetEffectiveness(t *testing.T) {
	tests := map[string]struct {
		attacking pokemonType
		defending pokemonType
		want      effectiveness
	}{
		"normal/ghost":     {attacking: normalType, defending: ghostType, want: immuneEffectivensss},
		"ghost/normal":     {attacking: ghostType, defending: normalType, want: immuneEffectivensss},
		"fighting/normal":  {attacking: fightingType, defending: normalType, want: superEffectiveness},
		"psychic/fighting": {attacking: psychicType, defending: fightingType, want: superEffectiveness},
		"steel/psychic":    {attacking: steelType, defending: psychicType, want: normalEffectivenss},
		"bug/steel":        {attacking: bugType, defending: steelType, want: resistedEffectiveness},
		"ground/bug":       {attacking: groundType, defending: bugType, want: resistedEffectiveness},
		"water/ground":     {attacking: waterType, defending: groundType, want: superEffectiveness},
		"ice/water":        {attacking: iceType, defending: waterType, want: resistedEffectiveness},
		"fire/ice":         {attacking: fireType, defending: iceType, want: superEffectiveness},
		"fairy/fire":       {attacking: fairyType, defending: fireType, want: resistedEffectiveness},
		"dragon/fairy":     {attacking: dragonType, defending: fairyType, want: immuneEffectivensss},
		"electric/dragon":  {attacking: electricType, defending: dragonType, want: resistedEffectiveness},
		"grass/electric":   {attacking: grassType, defending: electricType, want: normalEffectivenss},
		"rock/grass":       {attacking: rockType, defending: grassType, want: normalEffectivenss},
		"flying/rock":      {attacking: flyingType, defending: rockType, want: resistedEffectiveness},
		"dark/flying":      {attacking: darkType, defending: flyingType, want: normalEffectivenss},
		"poison/dark":      {attacking: poisonType, defending: darkType, want: normalEffectivenss},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := getEffectiveness(tc.attacking, tc.defending); got != tc.want {
				t.Errorf("%s: getEffectiveness(%s, %s) = %.1f, want %.1f", name, tc.attacking.String(), tc.defending.String(), got, tc.want)
			}
		})
	}
}
