package main

import (
	"container/heap"
	"log"
)

type singleBattleState struct {
	activePlayerSlot   *slot
	activeOpponentSlot *slot
	player             *trainer
	opponent           *trainer
	actions            *ActionQueue
}

func (sbs *singleBattleState) execute() {
	log.Println("Starting battle...")
	heap.Init(sbs.actions)
	for k := 0; !sbs.player.lost && !sbs.opponent.lost; k++ {
		log.Println("=====")
		log.Printf("Turn %d:\n", k+1)
		log.Printf("%s %d/%d - %s %d/%d", sbs.activePlayerSlot.mon.Base.Name, sbs.activePlayerSlot.mon.Hp, sbs.activePlayerSlot.mon.Stats["hp"], sbs.activeOpponentSlot.mon.Base.Name, sbs.activeOpponentSlot.mon.Hp, sbs.activeOpponentSlot.mon.Stats["hp"])

		sbs.gatherActions()
		for sbs.actions.Len() > 0 {
			action := heap.Pop(sbs.actions).(action)
			action.invoke(sbs)
		}
		resolveEndOfTurn(sbs)
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

func (sbs *singleBattleState) injectReplaceAction(slot *slot, trainer *trainer, midTurn bool) {
	heap.Push(sbs.actions, &replaceAction{
		oldSlot: slot,
		trainer: trainer,
		midTurn: midTurn,
	})
}

func (sbs *singleBattleState) getTrainer(slot *slot) *trainer {
	if slot == sbs.activePlayerSlot {
		return sbs.player
	}
	return sbs.opponent
}

func (sbs *singleBattleState) getActions() *ActionQueue {
	return sbs.actions
}

func (sbs *singleBattleState) getAllSlots() []*slot {
	return []*slot{
		sbs.activePlayerSlot,
		sbs.activeOpponentSlot,
	}
}

func initSingleBattleState(player, opponent trainer) *singleBattleState {
	return &singleBattleState{
		activePlayerSlot: &slot{
			firstTurn: true,
		},
		activeOpponentSlot: &slot{
			firstTurn: true,
		},
		player:   &player,
		opponent: &opponent,
		actions:  &ActionQueue{},
	}
}
