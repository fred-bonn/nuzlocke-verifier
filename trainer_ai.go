package main

import "math/rand"

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
