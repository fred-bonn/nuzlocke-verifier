package main

import "math/rand"

type ai interface {
	evaluateActions(sbs battleState, actions []action) action
}

type randomAi struct{}

func (ra randomAi) evaluateActions(sbs battleState, actions []action) action {
	return actions[rand.Intn(len(actions))]
}
