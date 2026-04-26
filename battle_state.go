package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"

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

func (sbs *singleBattleState) execute() {
	log.Println("Starting battle...")
	var actions []action

	for k := 0; k < 10; k++ {
		actions = sbs.gatherActions()
		sort.Slice(actions, func(i, j int) bool {
			if actions[i].prio() < actions[j].prio() {
				return false
			} else if actions[i].prio() > actions[j].prio() {
				return true
			} else if actions[i].speed() < actions[j].speed() {
				return false
			} else if actions[i].speed() > actions[j].speed() {
				return true
			}
			return (rand.Int() % 2) == 0
		})
		for _, action := range actions {
			action.invoke(sbs)
		}
	}
}

func (sbs *singleBattleState) gatherActions() []action {
	actions := make([]action, 0, 2)
	actions = append(actions, sbs.player.nextAction(sbs, sbs.activePlayerMon))
	actions = append(actions, sbs.opponent.nextAction(sbs, sbs.activeOpponentMon))
	return actions
}
