package main

import "math/rand"

type ai interface {
	evaluateActions(sbs *singleBattleState, actions []action) action
}

type randomAi struct{}

func (ra randomAi) evaluateActions(sbs *singleBattleState, actions []action) action {
	return actions[rand.Intn(len(actions))]
}
