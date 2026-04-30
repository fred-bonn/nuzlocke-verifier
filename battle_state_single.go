package main

import (
	"container/heap"
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type singleBattleState struct {
	activePlayerSlot   *slot
	activeOpponentSlot *slot
	player             *trainer
	opponent           *trainer
	actions            *ActionQueue
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
	heap.Init(sbs.actions)
	for k := 0; !sbs.player.lost && !sbs.opponent.lost; k++ {
		log.Println("=====")
		log.Printf("Turn %d:\n", k+1)

		sbs.gatherActions()
		for sbs.actions.Len() > 0 {
			action := heap.Pop(sbs.actions).(action)
			action.invoke(sbs)
		}
	}
	log.Println("=====")
	log.Println("Ending battle...")
}

func (sbs *singleBattleState) gatherActions() {
	heap.Push(sbs.actions, sbs.player.nextAction(sbs, sbs.activePlayerSlot))
	heap.Push(sbs.actions, sbs.opponent.nextAction(sbs, sbs.activeOpponentSlot))
}

func (sbs *singleBattleState) getOtherSlots(s *slot) []*slot {
	if s == sbs.activePlayerSlot {
		return []*slot{sbs.activeOpponentSlot}
	}
	return []*slot{sbs.activePlayerSlot}
}

func (sbs *singleBattleState) injectReplaceAction(slot *slot, trainer *trainer) {
	heap.Push(sbs.actions, &replaceAction{
		oldSlot: slot,
		trainer: trainer,
	})
}

func (sbs *singleBattleState) getTrainer(slot *slot) *trainer {
	if slot == sbs.activePlayerSlot {
		return sbs.player
	}
	return sbs.opponent
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
		actions:            &ActionQueue{},
	}, nil
}
