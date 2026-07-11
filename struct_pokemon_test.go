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

func TestChangeStatStageBy2(t *testing.T) {
	tests := map[string]struct {
		initial   int
		stat      stat
		change    int
		offensive bool
		ability   ability
		want      int
	}{
		"increases stage":               {initial: 0, stat: attack, change: 1, want: 1},
		"decreases stage":               {initial: 0, stat: defense, change: -1, want: -1},
		"caps positive stage":           {initial: 5, stat: speed, change: 2, want: 6},
		"caps negative stage":           {initial: -5, stat: speed, change: -2, want: -6},
		"clear body blocks offense":     {initial: 0, stat: attack, change: -1, offensive: true, ability: clearBodyAbility, want: 0},
		"clear body blocks not offense": {initial: 0, stat: attack, change: -1, offensive: false, ability: clearBodyAbility, want: -1},
		"keen eye blocks accuracy":      {initial: 0, stat: accuracy, change: -1, ability: keenEyeAbility, want: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stages:  []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				ability: tc.ability,
			}
			mon.stages[tc.stat] = tc.initial

			mon.changeStatStageBy(tc.stat, tc.change, tc.offensive)

			if got := mon.stages[tc.stat]; got != tc.want {
				t.Errorf("mon.changeStatStageBy(%s, %d, %t) = %d, want %d", tc.stat, tc.change, tc.offensive, got, tc.want)
			}
		})
	}
}

func TestHasAilment(t *testing.T) {
	tests := map[string]struct {
		has   ailmentState
		check ailmentState
		want  bool
	}{
		"para/para":     {paralysisAilment, paralysisAilment, true},
		"para/freeze":   {paralysisAilment, freezeAilment, false},
		"poison/poison": {poisonAilment, poisonAilment, true},
		"poison/toxic":  {poisonAilment, toxicAilment, false},
		"toxic/poison":  {toxicAilment, poisonAilment, false},
		"yawn/yawn":     {yawnAilment, yawnAilment, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				ailments: map[ailmentState]*ailment{
					tc.has: {
						state: tc.has,
					},
				},
			}

			if got := mon.hasAilment(tc.check) != nil; tc.want != got {
				t.Errorf("mon.hasAilment(%s) = %t, want %t", tc.check.String(), got, tc.want)
			}
		})
	}
}

func TestHasNonVolatileAilment(t *testing.T) {
	tests := map[string]struct {
		has  ailmentState
		want bool
	}{
		"para":        {paralysisAilment, true},
		"sleep":       {sleepAilment, true},
		"toxic":       {toxicAilment, true},
		"yawn":        {yawnAilment, false},
		"infatuation": {infatuationAilment, false},
		"confusion":   {confusionAilment, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				ailments: map[ailmentState]*ailment{
					tc.has: {
						state: tc.has,
					},
				},
			}

			if got := mon.hasNonVolatileAilment(); tc.want != got {
				t.Errorf("mon.hasAilment() = %t, want %t", got, tc.want)
			}
		})
	}
}

func TestApplyAilment(t *testing.T) {
	tests := map[string]struct {
		ailment     ailmentState
		pokemonType pokemonType
		ability     ability
		item        itemState
		want        bool
		removed     bool
	}{
		"burn":                  {burnAilment, normalType, intimidateAbility, noneItem, true, false},
		"para":                  {paralysisAilment, normalType, intimidateAbility, noneItem, true, false},
		"para cheri":            {paralysisAilment, normalType, intimidateAbility, cheriBerry, true, true},
		"para chesto":           {paralysisAilment, normalType, intimidateAbility, chestoBerry, true, false},
		"burn fire":             {burnAilment, fireType, intimidateAbility, noneItem, false, false},
		"para electric":         {paralysisAilment, electricType, intimidateAbility, noneItem, false, false},
		"confusion lum":         {confusionAilment, normalType, intimidateAbility, lumBerry, true, true},
		"confusion persim":      {confusionAilment, normalType, intimidateAbility, persimBerry, true, true},
		"freeze ice":            {freezeAilment, iceType, intimidateAbility, noneItem, false, false},
		"freeze magma armor":    {freezeAilment, normalType, magmaArmorAbility, noneItem, false, false},
		"freeze aspear":         {freezeAilment, normalType, noneAbility, aspearBerry, true, true},
		"burn water veil":       {burnAilment, normalType, waterVeilAbility, noneItem, false, false},
		"para limber":           {paralysisAilment, normalType, limberAbility, noneItem, false, false},
		"poison":                {poisonAilment, normalType, intimidateAbility, noneItem, true, false},
		"poison steel":          {poisonAilment, steelType, intimidateAbility, noneItem, false, false},
		"poison pecha":          {poisonAilment, normalType, intimidateAbility, pechaBerry, true, true},
		"poison immunity":       {poisonAilment, normalType, immunityAbility, noneItem, false, false},
		"sleep vital spirit":    {sleepAilment, normalType, vitalSpiritAbility, noneItem, false, false},
		"yawn vital spirit":     {yawnAilment, normalType, vitalSpiritAbility, noneItem, false, false},
		"infatuation oblivious": {infatuationAilment, normalType, obliviousAbility, noneItem, false, false},
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
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			got := mon.applyAilment(tc.ailment, nil, nil)
			if got != tc.want {
				t.Fatalf("mon.applyAilment(%s, nil, nil) = %t, want %t", tc.ailment.String(), got, tc.want)
			}
			if got {
				if tc.removed != (mon.hasAilment(tc.ailment) == nil) {
					t.Fatalf("mon.applyAilment(%s, nil, nil) applied %s, but the state after is wrong, removed: %t, has: %t", tc.ailment.String(), tc.ailment.String(), tc.removed, mon.hasAilment(tc.ailment) != nil)
				}

				gotAgain := mon.applyAilment(tc.ailment, nil, nil)
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

func TestToxicConversion(t *testing.T) {
	tests := map[string]struct {
		move string
		want bool
	}{
		"poison fang": {"poison fang", true},
		"toxic":       {"toxic", true},
		"tackle":      {"tackle", false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				ailments: make(map[ailmentState]*ailment),
			}
			move := Move{
				Name: tc.move,
			}
			item, _ := registerItem(noneItem, &mon)
			mon.item = item

			mon.applyAilment(poisonAilment, &move, nil)
			if tc.want && mon.hasAilment(toxicAilment) == nil {
				t.Fatalf("mon.applyAilment(poison, %s, nil) did not convert the poison to toxic", tc.move)
			}
			if !tc.want && mon.hasAilment(toxicAilment) != nil {
				t.Fatalf("mon.applyAilment(poison, %s, nil) did converted the poison to toxic", tc.move)
			}
		})
	}
}

func TestIsGrounded(t *testing.T) {
	tests := map[string]struct {
		pokemonType pokemonType
		ability     ability
		item        itemState
		want        bool
	}{
		"flying":             {flyingType, noneAbility, noneItem, false},
		"flying iron ball":   {flyingType, noneAbility, ironBall, true},
		"intim":              {normalType, intimidateAbility, noneItem, true},
		"intim iron ball":    {normalType, intimidateAbility, ironBall, true},
		"levitate":           {normalType, levitateAbility, noneItem, false},
		"levitate iron ball": {normalType, levitateAbility, ironBall, true},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				base: BasePokemon{
					Types: []pokemonType{tc.pokemonType},
				},
				ability: tc.ability,
			}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			got := mon.isGrounded()
			if got != tc.want {
				t.Fatalf("mon.isGrounded() = %t, want %t", got, tc.want)
			}
		})
	}
}
