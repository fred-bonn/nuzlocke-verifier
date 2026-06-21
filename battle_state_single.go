package main

import (
	"log"
)

type singleBattleState struct {
	activePlayerSlot   *slot
	activeOpponentSlot *slot
	player             *trainer
	opponent           *trainer
	actions            actionQueue
	weather            weatherState
}

func (sbs *singleBattleState) execute() {
	log.Println("Starting battle...")

	for k := 0; !sbs.player.lost && !sbs.opponent.lost; k++ {
		log.Println("=====")
		log.Printf("Turn %d:\n", k+1)
		log.Printf("%s %d/%d - %s %d/%d", sbs.activePlayerSlot.mon.base.Name, sbs.activePlayerSlot.mon.hp, sbs.activePlayerSlot.mon.maxHP(), sbs.activeOpponentSlot.mon.base.Name, sbs.activeOpponentSlot.mon.hp, sbs.activeOpponentSlot.mon.maxHP())

		sbs.gatherActions()
		sbs.actions.sort(sbs)
		for len(sbs.actions.queue) > 0 {
			action, _ := sbs.actions.queue.pop()
			action.invoke(sbs)
		}
		resolveEndOfTurn(sbs)
		// if the end of turn causes mons to faint, empty the queue for replace actions
		for len(sbs.actions.queue) > 0 {
			action, _ := sbs.actions.queue.pop()
			action.invoke(sbs)
		}
	}
	log.Println("=====")
	log.Println("Ending battle...")
}

func (sbs *singleBattleState) gatherActions() {
	sbs.actions.queue.push(sbs.player.nextAction(sbs, sbs.activePlayerSlot))
	sbs.actions.queue.push(sbs.opponent.nextAction(sbs, sbs.activeOpponentSlot))
}

func (sbs *singleBattleState) getAllSlots() []*slot {
	return []*slot{
		sbs.activePlayerSlot,
		sbs.activeOpponentSlot,
	}
}

func (sbs *singleBattleState) getOtherSlots(s *slot) []*slot {
	if s == sbs.activePlayerSlot {
		return []*slot{sbs.activeOpponentSlot}
	}
	return []*slot{sbs.activePlayerSlot}
}

func (sbs *singleBattleState) getOpponentSlot(s *slot) *slot {
	if s == sbs.activePlayerSlot {
		return sbs.activeOpponentSlot
	}
	return sbs.activePlayerSlot
}

func (sbs *singleBattleState) getActions() *actionQueue {
	return &sbs.actions
}

func (sbs *singleBattleState) getWeather() weatherState {
	return sbs.weather
}

func (sbs *singleBattleState) setWeather(w weatherState, turns int) {
	sbs.weather = w
}

func initSingleBattleState(player, opponent trainer, playerParty, opponentParty []*pokemon) *singleBattleState {
	player.pokemonParty = playerParty
	opponent.pokemonParty = opponentParty

	res := singleBattleState{
		activePlayerSlot: &slot{
			mon:       playerParty[0],
			trainer:   &player,
			firstTurn: true,
		},
		activeOpponentSlot: &slot{
			mon:       opponentParty[0],
			trainer:   &opponent,
			firstTurn: true,
		},
		player:   &player,
		opponent: &opponent,
		actions: actionQueue{
			queue: make(priorityQueue[action], 0, 3),
		},
	}

	resolveOnEntry(&res)

	return &res
}
