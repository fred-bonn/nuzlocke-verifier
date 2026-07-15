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
		"oran  heals below threshold hp":       {initialHP: 50, maxHP: 100, item: oranBerry, want: 60},
		"oran  does not heal below threshold":  {initialHP: 51, maxHP: 100, item: oranBerry, want: 51},
		"sitrus heals below threshold hp":      {initialHP: 50, maxHP: 100, item: sitrusBerry, want: 75},
		"sitrus does not heal below threshold": {initialHP: 51, maxHP: 100, item: oranBerry, unnerved: true, want: 51},
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
		"oran berry":   {item: oranBerry, want: true},
		"no item":      {item: noneItem, want: false},
		"yache berry":  {item: yacheBerry, want: true},
		"iron ball":    {item: ironBall, want: false},
		"liechi berry": {item: liechiBerry, want: true},
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
		"no item":       {item: noneItem, want: false},
		"oran berry":    {item: oranBerry, want: false},
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

func TestPinchHealBerries(t *testing.T) {
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
		"figy with lonely nature":                          {item: figyBerry, nature: "lonely", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"figy with bold nature confuses":                   {item: figyBerry, nature: "bold", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"figy heals at threshold":                          {item: figyBerry, nature: "lonely", initialHp: 50, maxHp: 100, gluttony: true, wantHp: 100},
		"figy does not heal above threshold":               {item: figyBerry, nature: "bold", initialHp: 51, maxHp: 100, gluttony: true, wantHp: 51},
		"iapapa with adamant nature":                       {item: iapapaBerry, nature: "adamant", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"iapapa with hasty nasture confuses":               {item: iapapaBerry, nature: "hasty", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"wiki with rash nature":                            {item: wikiBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"wiki with impish nature confuses":                 {item: wikiBerry, nature: "impish", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"aguav with relaxed nature":                        {item: aguavBerry, nature: "relaxed", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"aguav with rash nature confuses":                  {item: aguavBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"mago wtih modest nature":                          {item: magoBerry, nature: "modest", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90},
		"mago with sassy nature confuses":                  {item: magoBerry, nature: "sassy", initialHp: 40, maxHp: 100, gluttony: true, wantHp: 90, wantConfused: true},
		"iapapa does not heal dead mons":                   {item: iapapaBerry, nature: "adamant", initialHp: 0, maxHp: 100, gluttony: true, wantHp: 0},
		"aguav blocked by unnerve":                         {item: aguavBerry, nature: "rash", initialHp: 40, maxHp: 100, gluttony: true, unnerved: true, wantHp: 40},
		"figy without gluttony does not heal":              {item: figyBerry, nature: "lonely", initialHp: 40, maxHp: 100, wantHp: 40},
		"figy without gluttony heals at threshold":         {item: figyBerry, nature: "lonely", initialHp: 25, maxHp: 100, wantHp: 75},
		"figy without gluttony does not heal at threshold": {item: figyBerry, nature: "lonely", initialHp: 26, maxHp: 100, wantHp: 26},
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

			if tc.initialHp != tc.wantHp && !item.consumed {
				t.Errorf("%s was not consumed", tc.item.String())
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

			if tc.wantStage == 1 && !item.consumed {
				t.Errorf("%s was not consumed", tc.item.String())
			}
		})
	}
}

func TestResistBerries(t *testing.T) {
	tests := map[string]struct {
		item          itemState
		pokemon       pokemonType
		unnerved      bool
		event         bool
		initialDamage int
		want          int
	}{
		"chople reduces fighting damage":          {item: chopleBerry, pokemon: fightingType, event: true, initialDamage: 100, want: 50},
		"chople ignores non-matching type":        {item: chopleBerry, pokemon: normalType, event: true, initialDamage: 100, want: 100},
		"chople blocked by unnerve":               {item: chopleBerry, pokemon: fightingType, event: true, initialDamage: 100, unnerved: true, want: 100},
		"berry wont activate if event is missing": {item: chopleBerry, pokemon: fightingType, initialDamage: 100, want: 100},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{unnerved: tc.unnerved}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			damage := tc.initialDamage
			if tc.event {
				mon.checkItemTrigger(true, resistBerryEvent{
					pokemonType: tc.pokemon,
					damage:      &damage,
				})
			} else {
				mon.checkItemTrigger(true, nil)
			}

			if got := damage; got != tc.want {
				t.Errorf("damage = %d, want %d", got, tc.want)
			}

			if tc.initialDamage != tc.want && !item.consumed {
				t.Errorf("%s was not consumed", tc.item.String())
			}
		})
	}
}

func TestTypeGems(t *testing.T) {
	tests := map[string]struct {
		item         itemState
		move         pokemonType
		initialPower int
		event        bool
		want         int
	}{
		"normal gem increases normal type":              {item: normalGem, move: normalType, initialPower: 100, event: true, want: 150},
		"normal gem does not increase normal fire type": {item: normalGem, move: fireType, initialPower: 100, event: true, want: 100},
		"gem wont activate if event is missing":         {item: normalGem, move: normalType, initialPower: 100, want: 100},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			mon := pokemon{}
			item, _ := registerItem(tc.item, &mon)
			mon.item = item

			power := tc.initialPower
			if tc.event {
				mon.checkItemTrigger(true, gemEvent{
					pokemonType: tc.move,
					power:       &power,
				})
			} else {
				mon.checkItemTrigger(true, nil)
			}

			if got := power; got != tc.want {
				t.Errorf("power = %d, want %d", got, tc.want)
			}

			if tc.initialPower != tc.want && !item.consumed {
				t.Errorf("%s was not consumed", tc.item.String())
			}
		})
	}
}
