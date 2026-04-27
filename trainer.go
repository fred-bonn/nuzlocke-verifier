package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type trainer struct {
	pokemonParty []*pokemon.Pokemon
	player       bool
	ai           ai
}

func (t trainer) nextAction(sbs battleState, current *pokemon.Pokemon, slot int) action {
	possibleActions := make([]action, 0)
	opponentSlot := 1 - slot
	for _, mon := range t.pokemonParty {
		if mon == current || mon.Fainted {
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
			slot: opponentSlot,
			move: move,
		})
	}

	return t.ai.evaluateActions(sbs, possibleActions)
}
