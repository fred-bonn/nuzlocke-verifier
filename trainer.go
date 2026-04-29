package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type trainer struct {
	pokemonParty []*pokemon.Pokemon
	player       bool
	ai           ai
}

func (t trainer) nextAction(sbs battleState, slot *slot) action {
	possibleActions := make([]action, 0)
	opponentSlot := sbs.getOtherSlots(slot)[0] // only works for single battles for now
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted {
			continue
		}
		possibleActions = append(possibleActions, &swapAction{
			old: slot.mon,
			new: mon,
		})
	}
	for _, move := range slot.mon.Moves {
		possibleActions = append(possibleActions, &moveAction{
			mon:  slot.mon,
			slot: opponentSlot,
			move: move,
		})
	}

	return t.ai.evaluateActions(sbs, possibleActions)
}
