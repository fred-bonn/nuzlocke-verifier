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

func TestEffectiveStat(t *testing.T) {
	tests := map[string]struct {
		stat  stat
		crit  bool
		stage int
		base  int
		want  int
	}{
		"positive stage":                      {stat: attack, stage: 1, base: 100, want: 150},
		"negative stage":                      {stat: attack, stage: -1, base: 100, want: 66},
		"defense crit ignores positive stage": {stat: defense, crit: true, stage: 1, base: 100, want: 100},
		"attack crit ignores negative stage":  {stat: attack, crit: true, stage: -1, base: 100, want: 100},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				stages: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}
			mon.stats[tc.stat] = tc.base
			mon.stages[tc.stat] = tc.stage

			if got := mon.effectiveStat(tc.stat, tc.crit); got != tc.want {
				t.Errorf("mon.effectiveStat(%s, %t) = %d, want %d", tc.stat, tc.crit, got, tc.want)
			}
		})
	}
}

func TestEffectiveSpeed(t *testing.T) {
	tests := map[string]struct {
		stage     int
		base      int
		item      itemState
		ability   ability
		weather   weatherState
		paralyzed bool
		unburden  bool
		want      int
	}{
		"positive stage":                 {1, 100, noneItem, noneAbility, noneWeather, false, false, 150},
		"negative stage":                 {-1, 100, noneItem, noneAbility, noneWeather, false, false, 66},
		"paralyzed":                      {0, 100, noneItem, noneAbility, noneWeather, true, false, 25},
		"positive stage paralyzed":       {1, 100, noneItem, noneAbility, noneWeather, true, false, 37},
		"iron ball":                      {0, 100, ironBall, noneAbility, noneWeather, false, false, 50},
		"negative stage iron ball":       {-1, 100, ironBall, noneAbility, noneWeather, false, false, 33},
		"swift swim no rain":             {0, 100, noneItem, swiftSwimAbility, noneWeather, false, false, 100},
		"swift swim rain":                {0, 100, noneItem, swiftSwimAbility, rainWeather, false, false, 200},
		"sand rush rain":                 {0, 100, noneItem, sandRushAbility, rainWeather, false, false, 100},
		"sand rush sandstorm":            {0, 100, noneItem, sandRushAbility, sandstormWeather, false, false, 200},
		"slush rush rain":                {0, 100, noneItem, slushRushAbility, rainWeather, false, false, 100},
		"slush rush sandstorm iron ball": {0, 100, ironBall, slushRushAbility, hailWeather, false, false, 100},
		"chloro rain":                    {0, 100, noneItem, chlorophyllAbility, rainWeather, false, false, 100},
		"chloro sun positive stage":      {1, 100, noneItem, chlorophyllAbility, sunWeather, false, false, 300},
		"unburden no ability":            {0, 100, noneItem, noneAbility, noneWeather, false, true, 100},
		"unburden":                       {0, 100, noneItem, unburdenAbility, noneWeather, false, true, 200},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				stages:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				ability:  tc.ability,
				ailments: make(map[ailmentState]*ailment),
				unburden: tc.unburden,
			}
			mon.stats[speed] = tc.base
			mon.stages[speed] = tc.stage
			item, _ := registerItem(tc.item, &mon)
			mon.item = item
			bbs := initBenchBattleState(tc.weather)
			if tc.paralyzed {
				mon.applyAilment(paralysisAilment, nil, nil)
			}

			if got := mon.effectiveSpeed(bbs); got != tc.want {
				t.Errorf("mon.effectiveSpeed(...) = %d, want %d", got, tc.want)
			}
		})
	}
}
