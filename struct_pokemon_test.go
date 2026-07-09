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

func TestEvasionFraction(t *testing.T) {
	tests := map[string]struct {
		stage   int
		keenEye bool
		want    struct {
			num int
			dem int
		}
	}{
		"positive stage 1": {1, false, struct {
			num int
			dem int
		}{3, 4}},
		"positive stage 5": {5, false, struct {
			num int
			dem int
		}{3, 8}},
		"negative stage -1": {-1, false, struct {
			num int
			dem int
		}{4, 3}},
		"keen eye positve": {5, true, struct {
			num int
			dem int
		}{1, 1}},
		"keen eye netative": {-3, true, struct {
			num int
			dem int
		}{1, 1}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				stages: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}
			mon.stages[evasion] = tc.stage

			if num, dem := mon.evasionFraction(tc.keenEye); num != tc.want.num || dem != tc.want.dem {
				t.Errorf("mon.evasionFraction(%t) = (%d, %d), want (%d, %d)", tc.keenEye, num, dem, tc.want.num, tc.want.dem)
			}
		})
	}
}

func TestAccuracyFration(t *testing.T) {
	tests := map[string]struct {
		stage int
		want  struct {
			num int
			dem int
		}
	}{
		"positive stage 1": {1, struct {
			num int
			dem int
		}{4, 3}},
		"negative stage -1": {-1, struct {
			num int
			dem int
		}{3, 4}},
		"negative stage -5": {-5, struct {
			num int
			dem int
		}{3, 8}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				stages: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
			}
			mon.stages[accuracy] = tc.stage

			if num, dem := mon.accuracyFraction(); num != tc.want.num || dem != tc.want.dem {
				t.Errorf("mon.evasionFraction() = (%d, %d), want (%d, %d)", num, dem, tc.want.num, tc.want.dem)
			}
		})
	}
}

func TestApplyAilment(t *testing.T) {
	tests := map[string]struct {
		ailment     ailmentState
		move        string
		pokemonType pokemonType
		ability     ability
		item        itemState
		want        bool
		toxic       bool
		removed     bool
	}{
		"burn":                  {burnAilment, "none", normalType, intimidateAbility, noneItem, true, false, false},
		"para":                  {paralysisAilment, "none", normalType, intimidateAbility, noneItem, true, false, false},
		"para cheri":            {paralysisAilment, "none", normalType, intimidateAbility, cheriBerry, true, false, true},
		"para chesto":           {paralysisAilment, "none", normalType, intimidateAbility, chestoBerry, true, false, false},
		"burn fire":             {burnAilment, "none", fireType, intimidateAbility, noneItem, false, false, false},
		"para electric":         {paralysisAilment, "none", electricType, intimidateAbility, noneItem, false, false, false},
		"confusion lum":         {confusionAilment, "none", normalType, intimidateAbility, lumBerry, true, false, true},
		"confusion persim":      {confusionAilment, "none", normalType, intimidateAbility, persimBerry, true, false, true},
		"freeze ice":            {freezeAilment, "none", iceType, intimidateAbility, noneItem, false, false, false},
		"freeze magma armor":    {freezeAilment, "none", normalType, magmaArmorAbility, noneItem, false, false, false},
		"freeze aspear":         {freezeAilment, "none", normalType, noneAbility, aspearBerry, true, false, true},
		"burn water veil":       {burnAilment, "none", normalType, waterVeilAbility, noneItem, false, false, false},
		"para limber":           {paralysisAilment, "none", normalType, limberAbility, noneItem, false, false, false},
		"poison":                {poisonAilment, "none", normalType, intimidateAbility, noneItem, true, false, false},
		"poison steel":          {poisonAilment, "none", steelType, intimidateAbility, noneItem, false, false, false},
		"poison pecha":          {poisonAilment, "none", normalType, intimidateAbility, pechaBerry, true, false, true},
		"poison immunity":       {poisonAilment, "none", normalType, immunityAbility, noneItem, false, false, false},
		"poison toxic":          {poisonAilment, "toxic", normalType, noneAbility, noneItem, true, true, false},
		"poison poison fang":    {poisonAilment, "poison fang", normalType, noneAbility, noneItem, true, true, false},
		"sleep vital spirit":    {sleepAilment, "none", normalType, vitalSpiritAbility, noneItem, false, false, false},
		"yawn vital spirit":     {yawnAilment, "none", normalType, vitalSpiritAbility, noneItem, false, false, false},
		"infatuation oblivious": {infatuationAilment, "none", normalType, obliviousAbility, noneItem, false, false, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				base: BasePokemon{
					Types: []pokemonType{tc.pokemonType},
				},
				ability:  tc.ability,
				ailments: make(map[ailmentState]*ailment),
			}
			move := Move{
				Name: tc.move,
			}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			got := mon.applyAilment(tc.ailment, &move, nil)
			if got != tc.want {
				t.Fatalf("mon.applyAilment(%s, %s, nil) = %t, want %t", tc.ailment.String(), tc.move, got, tc.want)
			}
			if got {
				if tc.removed != (mon.hasAilment(tc.ailment) == nil) {
					if tc.ailment != poisonAilment || !tc.toxic {
						t.Fatalf("mon.applyAilment(%s, %s, nil) applied %s, but the state after is wrong, removed: %t, has: %t", tc.ailment.String(), tc.move, tc.ailment.String(), tc.removed, mon.hasAilment(tc.ailment) != nil)
					}
				}

				if !tc.removed && tc.ailment == poisonAilment && tc.toxic && (mon.hasAilment(poisonAilment) != nil || mon.hasAilment(toxicAilment) == nil) {
					t.Fatalf("poison was supposed to be converted to toxic due to %s", tc.move)
				}

				gotAgain := mon.applyAilment(tc.ailment, &move, nil)
				if got && tc.removed && !gotAgain {
					t.Fatalf("%s was removed but did not get re-applied the second time", tc.ailment.String())
				}
				if got && !tc.removed && gotAgain {
					t.Fatalf("%s was not removed but did get re-applied the second time", tc.ailment.String())
				}
			}
		})
	}
}
