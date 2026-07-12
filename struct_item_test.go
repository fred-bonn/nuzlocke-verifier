package main

import "testing"

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
				t.Fatalf("mon.hp hp = %d, want %d", got, tc.want)
			}
		})
	}
}
