package main

import (
	"math/rand"
	"sort"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type battleState interface {
	setMon(old, new *pokemon.Pokemon)
	getMon(slot *slot) *pokemon.Pokemon
	getOtherSlots(slot *slot) []*slot
	gatherActions()
	execute()
}

type slot struct {
	mon *pokemon.Pokemon
}

func sortActions(actions []action) {
	sort.Slice(actions, func(i, j int) bool {
		if actions[i].prio() < actions[j].prio() {
			return false
		} else if actions[i].prio() > actions[j].prio() {
			return true
		} else if actions[i].speed() < actions[j].speed() {
			return false
		} else if actions[i].speed() > actions[j].speed() {
			return true
		}
		return (rand.Int() % 2) == 0
	})
}
