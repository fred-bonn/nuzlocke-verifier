package main

import (
	"math/rand"
)

type trainer struct {
	pokemonParty []*Pokemon
	player       bool
	ai           ai
	lost         bool
}

func (t *trainer) nextAction(bs battleState, slot *slot) action {
	opponentSlot := bs.getOtherSlots(slot)[0] // only works for single battles for now

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

	if len(possibleActions) == 0 {
		possibleActions = append(possibleActions, &moveAction{
			userSlot:   slot,
			targetSlot: opponentSlot,
			move:       &struggleMove,
		})
	}

	action, score := t.ai.evaluateActions(bs, possibleActions)
	if score > 0 {
		return action
	}
	if roll(1, 2) || slot.mon.Hp <= slot.mon.Stats["hp"]/2 || !bs.getTrainer(slot).canReplace(bs) || slot.isTrapped() {
		return action
	}

	var possibleMons []*Pokemon
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

func (t *trainer) selectSwitchIn(bs battleState, slot *slot) *Pokemon {
	var possibleMons []*Pokemon
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
	evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int)
	evaluteSwitchIns(bs battleState, mons []*Pokemon, opponentSlot *slot) *Pokemon
}

type randomAi struct{}

func (ra randomAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	return actions[rand.Intn(len(actions))], 1
}

func (ra randomAi) evaluteSwitchIns(bs battleState, mons []*Pokemon, opponentSlot *slot) *Pokemon {
	return mons[rand.Intn(len(mons))]
}
