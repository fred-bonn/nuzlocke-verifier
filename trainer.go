package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type trainer struct {
	pokemonParty []*pokemon.Pokemon
	player       bool
	ai           ai
}

func (t trainer) nextAction(sbs *singleBattleState, current *pokemon.Pokemon) action {
	possibleActions := make([]action, 0)
	for _, mon := range t.pokemonParty {

		if mon == current {
			continue
		}
		possibleActions = append(possibleActions, &swapAction{
			old: current,
			new: mon,
		})
	}
	for _, move := range current.Moves {
		possibleActions = append(possibleActions, &moveAction{
			mon:  current,
			slot: 0,
			move: move,
		})
	}

	return t.ai.evaluateActions(sbs, possibleActions)
}
