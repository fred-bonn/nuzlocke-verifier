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
	opponentSlot := bs.getOtherSlots(slot)[0] // only works for single battles for now
	willSwitch := roll(1, 3) && slot.mon.Hp > slot.mon.Stats["hp"]/2 && bs.getTrainer(slot).canReplace(bs)
	if willSwitch && !slot.isTrapped() {
		var possibleMons []*pokemon.Pokemon
		for _, mon := range t.pokemonParty {
			if mon == slot.mon || mon.Fainted || bs.getActions().containstSwitchTo(mon) {
				continue
			}
			possibleMons = append(possibleMons, mon)
		}
		return &switchAction{
			oldSlot: slot,
			new:     t.ai.evaluteSwitchIns(bs, possibleMons, opponentSlot),
		}
	}

	possibleActions := make([]*moveAction, 0)
	for _, move := range slot.mon.Moves {
		if move.PP <= 0 {
			continue
		}
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

	return t.ai.evaluteSwitchIns(bs, possibleMons, bs.getOtherSlots(slot)[0]) // only works for single battles for now)
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
	evaluateActions(bs battleState, actions []*moveAction) *moveAction
	evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon
}

type randomAi struct{}

func (ra randomAi) evaluateActions(bs battleState, actions []*moveAction) *moveAction {
	return actions[rand.Intn(len(actions))]
}

func (ra randomAi) evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon {
	return mons[rand.Intn(len(mons))]
}
