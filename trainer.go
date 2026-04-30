package main

import (
	"math/rand"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type trainer struct {
	pokemonParty []*pokemon.Pokemon
	player       bool
	ai           ai
	lost         bool
}

func (t *trainer) nextAction(sbs battleState, slot *slot) action {
	possibleActions := make([]action, 0)
	opponentSlot := sbs.getOtherSlots(slot)[0] // only works for single battles for now
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted {
			continue
		}
		possibleActions = append(possibleActions, &switchAction{
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

func (t *trainer) selectSwitchIn(sbs battleState, slot *slot) *pokemon.Pokemon {
	var possibleMons []*pokemon.Pokemon
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted {
			continue
		}
		possibleMons = append(possibleMons, mon)
	}

	if len(possibleMons) == 0 {
		t.lost = true
		return nil
	}

	return t.ai.evaluteSwitchIns(sbs, possibleMons)
}

type ai interface {
	evaluateActions(sbs battleState, actions []action) action
	evaluteSwitchIns(sbs battleState, mons []*pokemon.Pokemon) *pokemon.Pokemon
}

type randomAi struct{}

func (ra randomAi) evaluateActions(sbs battleState, actions []action) action {
	return actions[rand.Intn(len(actions))]
}

func (ra randomAi) evaluteSwitchIns(sbs battleState, mons []*pokemon.Pokemon) *pokemon.Pokemon {
	return mons[rand.Intn(len(mons))]
}
