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

func (t *trainer) nextAction(bs battleState, slot *slot) action {
	possibleActions := make([]action, 0)
	opponentSlot := bs.getOtherSlots(slot)[0] // only works for single battles for now
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted || bs.getActions().containstSwitchTo(mon) {
			continue
		}
		possibleActions = append(possibleActions, &switchAction{
			oldSlot: slot,
			new:     mon,
		})
	}
	for _, move := range slot.mon.Moves {
		possibleActions = append(possibleActions, &moveAction{
			userSlot:   slot,
			targetSlot: opponentSlot,
			move:       &move,
		})
	}

	return t.ai.evaluateActions(bs, possibleActions)
}

func (t *trainer) selectSwitchIn(bs battleState, slot *slot) *pokemon.Pokemon {
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

	return t.ai.evaluteSwitchIns(bs, possibleMons)
}

func (t *trainer) canReplace(bs battleState) bool {
	count := 0
	for _, mon := range t.pokemonParty {
		if !mon.Fainted {
			count++
		}
		if count > 1 {
			return true
		}
	}
	return false
}

type ai interface {
	evaluateActions(bs battleState, actions []action) action
	evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon) *pokemon.Pokemon
}

type randomAi struct{}

func (ra randomAi) evaluateActions(bs battleState, actions []action) action {
	return actions[rand.Intn(len(actions))]
}

func (ra randomAi) evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon) *pokemon.Pokemon {
	return mons[rand.Intn(len(mons))]
}
