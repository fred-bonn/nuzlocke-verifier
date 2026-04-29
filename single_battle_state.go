package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type singleBattleState struct {
	activePlayerSlot   *slot
	activeOpponentSlot *slot
	player             *trainer
	opponent           *trainer
	actions            []action
}

func (sbs *singleBattleState) setMon(old, new *pokemon.Pokemon) {
	if old == sbs.activePlayerSlot.mon {
		sbs.activePlayerSlot.mon = new
	} else {
		sbs.activeOpponentSlot.mon = new
	}
}

func (sbs *singleBattleState) getMon(slot *slot) *pokemon.Pokemon {
	if slot == sbs.activePlayerSlot {
		return sbs.activePlayerSlot.mon
	}
	return sbs.activeOpponentSlot.mon
}

func (sbs *singleBattleState) execute() {
	log.Println("Starting battle...")

	for k := 0; k < 3; k++ {
		log.Println("=====")
		log.Printf("Turn %d:\n", k+1)
		sbs.gatherActions()
		sortActions(sbs.actions)
		for _, action := range sbs.actions {
			action.invoke(sbs)
		}
	}
	log.Println("=====")
	log.Println("Ending battle...")
}

func (sbs *singleBattleState) gatherActions() {
	sbs.actions = make([]action, 0, 2)
	sbs.actions = append(sbs.actions, sbs.player.nextAction(sbs, sbs.activePlayerSlot))
	sbs.actions = append(sbs.actions, sbs.opponent.nextAction(sbs, sbs.activeOpponentSlot))
}

func (sbs *singleBattleState) getOtherSlots(s *slot) []*slot {
	if s == sbs.activePlayerSlot {
		return []*slot{sbs.activeOpponentSlot}
	}
	return []*slot{sbs.activePlayerSlot}
}

func initSingleBattleState(player, opponent trainer) (*singleBattleState, error) {
	if len(player.pokemonParty) == 0 || len(opponent.pokemonParty) == 0 {
		return nil, fmt.Errorf("player or opponent has no pokemon in their party")
	}

	return &singleBattleState{
		activePlayerSlot:   &slot{mon: player.pokemonParty[0]},
		activeOpponentSlot: &slot{mon: opponent.pokemonParty[0]},
		player:             &player,
		opponent:           &opponent,
	}, nil
}
