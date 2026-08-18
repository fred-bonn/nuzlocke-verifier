package main

import "testing"

func TestLearningAiTracksPolicyFromBattleState(t *testing.T) {
	la := &learningAi{policy: map[string]int{"damage": 0, "status": 0, "healing": 0, "switch": 0}}

	before := la.policy["damage"]
	la.observeOutcome("damage", 3)
	if got := la.policy["damage"]; got <= before {
		t.Fatalf("learning policy did not increase after a strong outcome: before=%d after=%d", before, got)
	}

	if len(la.policy) == 0 {
		t.Fatal("expected a non-empty learning policy for battle-state decisions")
	}
}

func TestSingleBattleStateResetRestoresInitialBattleState(t *testing.T) {
	playerAI := &learningAi{policy: map[string]int{"damage": 7, "status": 3, "healing": 2, "switch": 5}}
	playerParty := []*pokemon{{
		base:  BasePokemon{Name: "player mon"},
		stats: []int{100, 10, 10, 10, 10, 10},
		hp:    100,
		moves: []*Move{{Name: "tackle"}},
	}}
	opponentParty := []*pokemon{{
		base:  BasePokemon{Name: "opp mon"},
		stats: []int{100, 10, 10, 10, 10, 10},
		hp:    100,
		moves: []*Move{{Name: "tackle"}},
	}}

	sbs := initSingleBattleState(
		trainer{ai: playerAI, player: true, fieldEffects: make(map[fieldEffect]int)},
		trainer{ai: rnbAi{}, fieldEffects: make(map[fieldEffect]int)},
		playerParty,
		opponentParty,
		noneWeather,
	)

	sbs.activePlayerSlot.mon.hp = 1
	sbs.activePlayerSlot.firstTurn = false
	sbs.player.lost = true
	sbs.activeOpponentSlot.mon.hp = 2
	sbs.opponent.lost = true
	if err := sbs.reset(); err != nil {
		t.Fatalf("reset returned error: %v", err)
	}

	if sbs.activePlayerSlot.mon.hp != 100 || sbs.activeOpponentSlot.mon.hp != 100 {
		t.Fatalf("reset did not restore hp values: player=%d opponent=%d", sbs.activePlayerSlot.mon.hp, sbs.activeOpponentSlot.mon.hp)
	}
	if sbs.player.lost || sbs.opponent.lost {
		t.Fatal("reset did not clear trainer loss state")
	}
	if !sbs.activePlayerSlot.firstTurn || !sbs.activeOpponentSlot.firstTurn {
		t.Fatal("reset did not restore first-turn state")
	}
	if got := sbs.player.ai.(*learningAi).policy["damage"]; got != 7 {
		t.Fatalf("reset unexpectedly changed learning policy: got %d want 7", got)
	}
}
