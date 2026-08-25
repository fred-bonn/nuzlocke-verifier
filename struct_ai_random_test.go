package main

import (
	"strings"
	"testing"
)

func TestLearningAiReturnsAndRecordsFirstActionForState(t *testing.T) {
	player := testSwitchPokemon("player", 100, 100, 100, 100, nil)
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, nil)
	bs := testSwitchBattleState(player, opponent, opponent)
	bs.player.player = true
	preferred := &Move{Name: "preferred", Power: 1, PP: 1, Class: physicalClass}
	other := &Move{Name: "other", Power: 100, PP: 1, Class: physicalClass}
	player.moves = []*Move{preferred, other}
	la := newLearningAi()
	got, _ := la.evaluateActions(bs, []*moveAction{
		{userSlot: bs.activePlayerSlot, targetSlot: bs.activeOpponentSlot, move: preferred},
		{userSlot: bs.activePlayerSlot, targetSlot: bs.activeOpponentSlot, move: other},
	})
	options := la.policy[discretizeBattleState(bs)]
	if got.move != preferred || len(options) != 2 || options[0] != "move:preferred" || options[1] != "move:other" {
		t.Fatalf("learning AI did not retain all move options: got %v", options)
	}
}

func TestLearningAiReturnsAndRecordsFirstSwitchForState(t *testing.T) {
	player := testSwitchPokemon("player", 100, 100, 100, 100, nil)
	replacement := testSwitchPokemon("replacement", 100, 100, 100, 100, nil)
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, nil)
	bs := testSwitchBattleState(player, replacement, opponent)
	bs.player.player = true
	la := newLearningAi()
	if got := la.evaluteSwitchIns(bs, []*pokemon{replacement}, bs.activeOpponentSlot); got != replacement {
		t.Fatal("learning AI did not select first replacement")
	}
	options := la.policy[discretizeBattleState(bs)]
	if len(options) != 1 || options[0] != "switch:replacement" {
		t.Fatalf("learning AI did not retain switch option: got %v", options)
	}
}

func TestDiscretizeBattleStateIsCanonical(t *testing.T) {
	player := testSwitchPokemon("player", 100, 100, 100, 100, &Move{Name: "tackle", Power: 100, PP: 1, Class: physicalClass})
	opponent := testSwitchPokemon("opponent", 100, 50, 100, 100, &Move{Name: "sonic boom", Power: 1, PP: 1, Class: specialClass})
	player.ailments[burnAilment] = &ailment{state: burnAilment}
	player.ailments[confusionAilment] = &ailment{state: confusionAilment}
	bs := testSwitchBattleState(player, opponent, opponent)
	bs.player.player = true
	player.hp = 20
	opponent.hp = 20

	first := discretizeBattleState(bs)
	second := discretizeBattleState(bs)
	if first != second {
		t.Fatalf("state encoding is not stable: first=%s second=%s", first, second)
	}
	if !contains(first, `"player_mon_is_faster":true`) {
		t.Fatalf("state did not record player speed: %s", first)
	}
	if !contains(first, `"player_ailments":["burn","confusion"]`) {
		t.Fatalf("state did not canonicalize player ailments: %s", first)
	}
	if !contains(first, `"player_mon_has_move_that_kills":true`) {
		t.Fatalf("state did not record player kill move: %s", first)
	}
	if !contains(first, `"opponent_moves_that_kill":["sonic boom"]`) {
		t.Fatalf("state did not record opponent kill move: %s", first)
	}
}

func contains(value, substring string) bool {
	return strings.Contains(value, substring)
}

func TestSingleBattleStateResetRestoresInitialBattleState(t *testing.T) {
	playerAI := &learningAi{policy: map[string][]string{}}
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
	if got := len(sbs.player.ai.(*learningAi).policy); got != 0 {
		t.Fatalf("reset unexpectedly changed learning policy: got %d entries", got)
	}
}
