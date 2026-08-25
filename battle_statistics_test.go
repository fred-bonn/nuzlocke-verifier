package main

import "testing"

func TestBattleStatisticsRecord(t *testing.T) {
	survivor := &pokemon{base: BasePokemon{Name: "survivor"}}
	fainted := &pokemon{base: BasePokemon{Name: "fainted"}, fainted: true}
	player := &trainer{pokemonParty: []*pokemon{survivor, fainted}}
	statistics := newBattleStatistics(player.pokemonParty)

	statistics.record(player)
	player.lost = true
	statistics.record(player)

	if statistics.battleCount != 2 || statistics.winCount != 1 {
		t.Fatalf("unexpected battle totals: battles=%d wins=%d", statistics.battleCount, statistics.winCount)
	}
	if got := statistics.pokemonSurvivors[0]; got != 2 {
		t.Fatalf("survivor count = %d, want 2", got)
	}
	if got := statistics.pokemonSurvivors[1]; got != 0 {
		t.Fatalf("fainted Pokemon survivor count = %d, want 0", got)
	}
}
