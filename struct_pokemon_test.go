package main

import "testing"

func TestIsImmuneToPowderMoves(t *testing.T) {
	tests := map[string]struct {
		t    pokemonType
		a    ability
		want bool
	}{
		"grass":              {t: grassType, a: intimidateAbility, want: true},
		"grass overcoat:":    {t: grassType, a: overcoatAbility, want: true},
		"not grass 1":        {t: flyingType, a: intimidateAbility, want: false},
		"not grass 2":        {t: waterType, a: intimidateAbility, want: false},
		"not grass overcoat": {t: fireType, a: overcoatAbility, want: true},
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
		"positive stage":                 {stage: 1, base: 100, item: noneItem, ability: noneAbility, weather: noneWeather, paralyzed: false, unburden: false, want: 150},
		"negative stage":                 {stage: -1, base: 100, item: noneItem, ability: noneAbility, weather: noneWeather, paralyzed: false, unburden: false, want: 66},
		"paralyzed":                      {stage: 0, base: 100, item: noneItem, ability: noneAbility, weather: noneWeather, paralyzed: true, unburden: false, want: 25},
		"positive stage paralyzed":       {stage: 1, base: 100, item: noneItem, ability: noneAbility, weather: noneWeather, paralyzed: true, unburden: false, want: 37},
		"iron ball":                      {stage: 0, base: 100, item: ironBall, ability: noneAbility, weather: noneWeather, paralyzed: false, unburden: false, want: 50},
		"negative stage iron ball":       {stage: -1, base: 100, item: ironBall, ability: noneAbility, weather: noneWeather, paralyzed: false, unburden: false, want: 33},
		"swift swim no rain":             {stage: 0, base: 100, item: noneItem, ability: swiftSwimAbility, weather: noneWeather, paralyzed: false, unburden: false, want: 100},
		"swift swim rain":                {stage: 0, base: 100, item: noneItem, ability: swiftSwimAbility, weather: rainWeather, paralyzed: false, unburden: false, want: 200},
		"sand rush rain":                 {stage: 0, base: 100, item: noneItem, ability: sandRushAbility, weather: rainWeather, paralyzed: false, unburden: false, want: 100},
		"sand rush sandstorm":            {stage: 0, base: 100, item: noneItem, ability: sandRushAbility, weather: sandstormWeather, paralyzed: false, unburden: false, want: 200},
		"slush rush rain":                {stage: 0, base: 100, item: noneItem, ability: slushRushAbility, weather: rainWeather, paralyzed: false, unburden: false, want: 100},
		"slush rush sandstorm iron ball": {stage: 0, base: 100, item: ironBall, ability: slushRushAbility, weather: hailWeather, paralyzed: false, unburden: false, want: 100},
		"chloro rain":                    {stage: 0, base: 100, item: noneItem, ability: chlorophyllAbility, weather: rainWeather, paralyzed: false, unburden: false, want: 100},
		"chloro sun positive stage":      {stage: 1, base: 100, item: noneItem, ability: chlorophyllAbility, weather: sunWeather, paralyzed: false, unburden: false, want: 300},
		"unburden no ability":            {stage: 0, base: 100, item: noneItem, ability: noneAbility, weather: noneWeather, paralyzed: false, unburden: true, want: 100},
		"unburden":                       {stage: 0, base: 100, item: noneItem, ability: unburdenAbility, weather: noneWeather, paralyzed: false, unburden: true, want: 200},
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
		"positive stage 1": {stage: 1, keenEye: false, want: struct {
			num int
			dem int
		}{num: 3, dem: 4}},
		"positive stage 5": {stage: 5, keenEye: false, want: struct {
			num int
			dem int
		}{num: 3, dem: 8}},
		"negative stage -1": {stage: -1, keenEye: false, want: struct {
			num int
			dem int
		}{num: 4, dem: 3}},
		"keen eye positve": {stage: 5, keenEye: true, want: struct {
			num int
			dem int
		}{num: 1, dem: 1}},
		"keen eye netative": {stage: -3, keenEye: true, want: struct {
			num int
			dem int
		}{num: 1, dem: 1}},
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
		"positive stage 1": {stage: 1, want: struct {
			num int
			dem int
		}{num: 4, dem: 3}},
		"negative stage -1": {stage: -1, want: struct {
			num int
			dem int
		}{num: 3, dem: 4}},
		"negative stage -5": {stage: -5, want: struct {
			num int
			dem int
		}{num: 3, dem: 8}},
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

func TestChangeStatStageBy(t *testing.T) {
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
		"para/para":     {has: paralysisAilment, check: paralysisAilment, want: true},
		"para/freeze":   {has: paralysisAilment, check: freezeAilment, want: false},
		"poison/poison": {has: poisonAilment, check: poisonAilment, want: true},
		"poison/toxic":  {has: poisonAilment, check: toxicAilment, want: false},
		"toxic/poison":  {has: toxicAilment, check: poisonAilment, want: false},
		"yawn/yawn":     {has: yawnAilment, check: yawnAilment, want: true},
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
		"para":        {has: paralysisAilment, want: true},
		"sleep":       {has: sleepAilment, want: true},
		"toxic":       {has: toxicAilment, want: true},
		"yawn":        {has: yawnAilment, want: false},
		"infatuation": {has: infatuationAilment, want: false},
		"confusion":   {has: confusionAilment, want: false},
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
		want        bool
	}{
		"burn":                  {ailment: burnAilment, pokemonType: normalType, ability: intimidateAbility, want: true},
		"burn fire":             {ailment: burnAilment, pokemonType: fireType, ability: intimidateAbility, want: false},
		"para":                  {ailment: paralysisAilment, pokemonType: normalType, ability: intimidateAbility, want: true},
		"para electric":         {ailment: paralysisAilment, pokemonType: electricType, ability: intimidateAbility, want: false},
		"freeze":                {ailment: freezeAilment, pokemonType: normalType, ability: noneAbility, want: true},
		"freeze ice":            {ailment: freezeAilment, pokemonType: iceType, ability: intimidateAbility, want: false},
		"freeze magma armor":    {ailment: freezeAilment, pokemonType: normalType, ability: magmaArmorAbility, want: false},
		"poison":                {ailment: poisonAilment, pokemonType: normalType, ability: intimidateAbility, want: true},
		"poison steel":          {ailment: poisonAilment, pokemonType: steelType, ability: intimidateAbility, want: false},
		"poison immunity":       {ailment: poisonAilment, pokemonType: normalType, ability: immunityAbility, want: false},
		"sleep":                 {ailment: sleepAilment, pokemonType: normalType, ability: noneAbility, want: true},
		"sleep vital spirit":    {ailment: sleepAilment, pokemonType: normalType, ability: vitalSpiritAbility, want: false},
		"yawn":                  {ailment: yawnAilment, pokemonType: normalType, ability: noneAbility, want: true},
		"yawn vital spirit":     {ailment: yawnAilment, pokemonType: normalType, ability: vitalSpiritAbility, want: false},
		"confusion":             {ailment: confusionAilment, pokemonType: normalType, ability: noneAbility, want: true},
		"infatuation":           {ailment: infatuationAilment, pokemonType: normalType, ability: noneAbility, want: true},
		"infatuation oblivious": {ailment: infatuationAilment, pokemonType: normalType, ability: obliviousAbility, want: false},
		"burn water veil":       {ailment: burnAilment, pokemonType: normalType, ability: waterVeilAbility, want: false},
		"para limber":           {ailment: paralysisAilment, pokemonType: normalType, ability: limberAbility, want: false},
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
			item, _ := registerItem(noneItem, &mon)
			mon.item = item

			got := mon.applyAilment(tc.ailment, nil, nil)
			if got != tc.want {
				t.Fatalf("mon.applyAilment(%s, nil, nil) = %t, want %t", tc.ailment.String(), got, tc.want)
			}
			if tc.want && mon.hasAilment(tc.ailment) == nil {
				t.Fatalf("mon.applyAilment(%s, nil, nil) did not leave %s applied", tc.ailment.String(), tc.ailment.String())
			}
			if !tc.want && mon.hasAilment(tc.ailment) != nil {
				t.Fatalf("mon.applyAilment(%s, nil, nil) applied %s despite a blocking condition", tc.ailment.String(), tc.ailment.String())
			}
		})
	}
}

func TestChangeHpBy(t *testing.T) {
	tests := map[string]struct {
		initialHP int
		maxHP     int
		change    int
		want      int
	}{
		"increase hp":       {initialHP: 50, maxHP: 100, change: 20, want: 70},
		"increase over max": {initialHP: 90, maxHP: 100, change: 20, want: 100},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats: []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				hp:    tc.initialHP,
			}
			mon.stats[hitPoints] = tc.maxHP

			mon.changeHpBy(tc.change)

			if got := mon.hp; got != tc.want {
				t.Fatalf("mon.changeHpBy(%d) hp = %d, want %d", tc.change, got, tc.want)
			}
		})
	}
}

func TestToxicConversion(t *testing.T) {
	tests := map[string]struct {
		move string
		want bool
	}{
		"poison fang": {move: "poison fang", want: true},
		"toxic":       {move: "toxic", want: true},
		"tackle":      {move: "tackle", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				ailments: make(map[ailmentState]*ailment),
			}
			move := Move{
				Name: tc.move,
			}

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
		"flying":             {pokemonType: flyingType, ability: noneAbility, item: noneItem, want: false},
		"flying iron ball":   {pokemonType: flyingType, ability: noneAbility, item: ironBall, want: true},
		"intim":              {pokemonType: normalType, ability: intimidateAbility, item: noneItem, want: true},
		"intim iron ball":    {pokemonType: normalType, ability: intimidateAbility, item: ironBall, want: true},
		"levitate":           {pokemonType: normalType, ability: levitateAbility, item: noneItem, want: false},
		"levitate iron ball": {pokemonType: normalType, ability: levitateAbility, item: ironBall, want: true},
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

			if got := mon.isGrounded(); got != tc.want {
				t.Fatalf("mon.isGrounded() = %t, want %t", got, tc.want)
			}
		})
	}
}
