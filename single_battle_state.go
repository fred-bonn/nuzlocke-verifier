package main

import (
	"fmt"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type singleBattleState struct {
	activePlayerMon   *pokemon.Pokemon
	activeOpponentMon *pokemon.Pokemon
	player            *trainer
	opponent          *trainer
}

func (sbs *singleBattleState) setMon(old, new *pokemon.Pokemon) {
	if old == sbs.activePlayerMon {
		sbs.activePlayerMon = new
	} else {
		sbs.activeOpponentMon = new
	}
}

func (sbs *singleBattleState) getMon(slot int) *pokemon.Pokemon {
	if slot == 0 {
		return sbs.activePlayerMon
	}
	return sbs.activeOpponentMon
}

func (sbs *singleBattleState) execute() {
	log.Println("Starting battle...")
	var actions []action

	for k := 0; k < 20; k++ {
		log.Println("=====")
		log.Printf("Turn %d:\n", k+1)
		actions = sbs.gatherActions()
		sortActions(actions)
		for _, action := range actions {
			action.invoke(sbs)
		}
	}
	log.Println("=====")
	log.Println("Ending battle...")
}

func (sbs *singleBattleState) gatherActions() []action {
	actions := make([]action, 0, 2)
	actions = append(actions, sbs.player.nextAction(sbs, sbs.activePlayerMon, 0))
	actions = append(actions, sbs.opponent.nextAction(sbs, sbs.activeOpponentMon, 1))
	return actions
}

func initSingleBattleState(player, opponent trainer) (*singleBattleState, error) {
	if len(player.pokemonParty) == 0 || len(opponent.pokemonParty) == 0 {
		return nil, fmt.Errorf("player or opponent has no pokemon in their party")
	}

	return &singleBattleState{
		activePlayerMon:   player.pokemonParty[0],
		activeOpponentMon: opponent.pokemonParty[0],
		player:            &player,
		opponent:          &opponent,
	}, nil
}
