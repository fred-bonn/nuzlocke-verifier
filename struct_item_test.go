package main

import (
	"testing"
)

func TestBerryItemsCureAilments(t *testing.T) {
	tests := map[string]struct {
		ailment  ailmentState
		item     itemState
		unnerved bool
		want     bool
	}{
		"cheri cures paralysis":       {ailment: paralysisAilment, item: cheriBerry, want: true},
		"cheri blocked by unnerve":    {ailment: paralysisAilment, item: cheriBerry, unnerved: true, want: false},
		"chesto cures sleep":          {ailment: sleepAilment, item: chestoBerry, want: true},
		"chesto does not cure poison": {ailment: poisonAilment, item: chestoBerry, want: false},
		"pecha cures poison":          {ailment: poisonAilment, item: pechaBerry, want: true},
		"rawst cures burn":            {ailment: burnAilment, item: rawstBerry, want: true},
		"aspear cures freeze":         {ailment: freezeAilment, item: aspearBerry, want: true},
		"persim cures confusion":      {ailment: confusionAilment, item: persimBerry, want: true},
		"persim does not cure freeze": {ailment: freezeAilment, item: persimBerry, want: false},
		"lum cures paralysis":         {ailment: paralysisAilment, item: lumBerry, want: true},
		"lum cures sleep":             {ailment: sleepAilment, item: lumBerry, want: true},
		"lum blocked by unnerve":      {ailment: paralysisAilment, item: lumBerry, unnerved: true, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				base: BasePokemon{
					Types: []pokemonType{normalType},
				},
				unnerved: tc.unnerved,
				ailments: make(map[ailmentState]*ailment),
			}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			if got := mon.applyAilment(tc.ailment, nil, nil); got != true {
				t.Fatalf("mon.applyAilment(%s, nil, nil) = %t, want true", tc.ailment.String(), got)
			}
			if got := mon.hasAilment(tc.ailment) == nil; got != tc.want {
				t.Fatalf("mon.hasAilment(%s) = %t, want %t", tc.ailment.String(), got, tc.want)
			}
		})
	}
}

func TestBerryItemsHeal(t *testing.T) {
	tests := map[string]struct {
		initialHP int
		maxHP     int
		item      itemState
		unnerved  bool
		want      int
	}{
		"oran berry restores hp":         {initialHP: 25, maxHP: 100, item: oranBerry, want: 35},
		"oran berry not below threshold": {initialHP: 60, maxHP: 100, item: oranBerry, want: 60},
		"sitrus berry restores hp":       {initialHP: 25, maxHP: 100, item: sitrusBerry, want: 50},
		"oran berry blocked by unnerve":  {initialHP: 25, maxHP: 100, item: oranBerry, unnerved: true, want: 25},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				hp:       tc.initialHP,
				unnerved: tc.unnerved,
			}
			mon.stats[hitPoints] = tc.maxHP
			item, _ := registerItem(tc.item, &mon)
			mon.item = item
			mon.checkItemTrigger(true, nil)

			if got := mon.hp; got != tc.want {
				t.Errorf("mon.hp = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestIsBerry(t *testing.T) {
	tests := map[string]struct {
		item itemState
		want bool
	}{
		"oran":         {item: oranBerry, want: true},
		"none":         {item: noneItem, want: false},
		"yache:":       {item: yacheBerry, want: true},
		"iron ball":    {item: ironBall, want: false},
		"liechi":       {item: liechiBerry, want: true},
		"mystic water": {item: mysticWater, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.item.isBerry(); got != tc.want {
				t.Errorf("%s.isBerry() = %t, want %t", tc.item.String(), got, tc.want)
			}
		})
	}
}

func TestIsChoice(t *testing.T) {
	tests := map[string]struct {
		item itemState
		want bool
	}{
		"choice scarf":  {item: choiceScarf, want: true},
		"choice specs":  {item: choiceSpecs, want: true},
		"focus sash:":   {item: focusSash, want: false},
		"none":          {item: noneItem, want: false},
		"oran":          {item: oranBerry, want: false},
		"silver powder": {item: silverPowder, want: false},
		"choice band":   {item: choiceBand, want: true},
		"assault vest":  {item: assaultVest, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.item.isChoice(); got != tc.want {
				t.Errorf("%s.isChoice() = %t, want %t", tc.item.String(), got, tc.want)
			}
		})
	}
}

func TestPinchHealingBerries(t *testing.T) {
	tests := map[string]struct {
		item         itemState
		nature       string
		initialHp    int
		maxHp        int
		gluttony     bool
		unnerved     bool
		wantHp       int
		wantConfused bool
	}{
		"figy lonely":               {item: figyBerry, nature: "lonely", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"figy bold":                 {item: figyBerry, nature: "bold", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"figy lonely 50":            {item: figyBerry, nature: "lonely", initialHp: 50, maxHp: 100, gluttony: true, wantHp: 100},
		"figy bold 51":              {item: figyBerry, nature: "bold", initialHp: 51, maxHp: 100, gluttony: true, wantHp: 51},
		"iapapa adamant":            {item: iapapaBerry, nature: "adamant", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"iapapa hasty":              {item: iapapaBerry, nature: "hasty", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"wiki rash":                 {item: wikiBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"wiki impish":               {item: wikiBerry, nature: "impish", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"aguav relaxed":             {item: aguavBerry, nature: "relaxed", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"aguav rash":                {item: aguavBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"mago modest":               {item: magoBerry, nature: "modest", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"mago sassy":                {item: magoBerry, nature: "sassy", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"iapapa adamant dead":       {item: iapapaBerry, nature: "adamant", initialHp: 0, maxHp: 100, gluttony: true, wantHp: 0},
		"aguav rash unnerved":       {item: aguavBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, unnerved: true, wantHp: 40},
		"figy lonely no ability":    {item: figyBerry, nature: "lonely", initialHp: 40, maxHp: 100, wantHp: 40},
		"figy lonely no ability 25": {item: figyBerry, nature: "lonely", initialHp: 25, maxHp: 100, wantHp: 75},
		"figy lonely no ability 26": {item: figyBerry, nature: "lonely", initialHp: 26, maxHp: 100, wantHp: 26},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				ailments: make(map[ailmentState]*ailment),
				hp:       tc.initialHp,
				unnerved: tc.unnerved,
			}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item
			nat, _ := getNature(tc.nature)
			mon.nat = nat
			mon.stats[hitPoints] = tc.maxHp
			if tc.gluttony {
				mon.ability = gluttonyAbility
			}

			mon.checkItemTrigger(true, nil)
			if got := mon.hp; got != tc.wantHp {
				t.Errorf("mon.hp = %d, want %d", got, tc.wantHp)
			}

			ailment := mon.hasAilment(confusionAilment)
			if got := ailment != nil; got != tc.wantConfused {
				if got {
					t.Errorf("mon was confused by %s with a %s nature", tc.item.String(), tc.nature)
				} else {
					t.Errorf("mon was confused by %s with a %s nature", tc.item.String(), tc.nature)
				}
			}
		})
	}
}

func TestStatBoostBerries(t *testing.T) {
	tests := map[string]struct {
		item      itemState
		initialHp int
		maxHp     int
		gluttony  bool
		unnerved  bool
		wantStage int
	}{
		"liechi boosts attack at threshold":                   {item: liechiBerry, initialHp: 25, maxHp: 100, wantStage: 1},
		"liechi does not boost above threshold":               {item: liechiBerry, initialHp: 26, maxHp: 100, wantStage: 0},
		"liechi blocked by unnerve":                           {item: liechiBerry, initialHp: 25, maxHp: 100, unnerved: true, wantStage: 0},
		"liechi uses gluttony threshold":                      {item: liechiBerry, initialHp: 50, maxHp: 100, gluttony: true, wantStage: 1},
		"liechi does not boost with gluttony above threshold": {item: liechiBerry, initialHp: 51, maxHp: 100, gluttony: true, wantStage: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{
				stats:    []int{0, 0, 0, 0, 0, 0, 0, 0, 0},
				stages:   []int{0, 0, 0, 0, 0, 0, 0, 0},
				hp:       tc.initialHp,
				unnerved: tc.unnerved,
			}
			mon.stats[hitPoints] = tc.maxHp
			if tc.gluttony {
				mon.ability = gluttonyAbility
			}

			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			mon.checkItemTrigger(true, nil)
			if got := mon.stages[attack]; got != tc.wantStage {
				t.Errorf("mon.stages[attack] = %d, want %d", got, tc.wantStage)
			}
		})
	}
}

func TestResistBerries(t *testing.T) {
	tests := map[string]struct {
		item          itemState
		pokemon       pokemonType
		unnerved      bool
		initialDamage int
		want          int
	}{
		"chople reduces fighting damage":   {item: chopleBerry, pokemon: fightingType, initialDamage: 100, want: 50},
		"chople ignores non-matching type": {item: chopleBerry, pokemon: normalType, initialDamage: 100, want: 100},
		"chople blocked by unnerve":        {item: chopleBerry, pokemon: fightingType, initialDamage: 100, unnerved: true, want: 100},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{unnerved: tc.unnerved}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			damage := tc.initialDamage
			mon.checkItemTrigger(true, resistBerryEvent{
				pokemonType: tc.pokemon,
				damage:      &damage,
			})

			if got := damage; got != tc.want {
				t.Errorf("damage = %d, want %d", got, tc.want)
			}
		})
	}
}
