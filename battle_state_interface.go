package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type battleState interface {
	setMon(old, new *pokemon.Pokemon)
	getMon(slot *slot) *pokemon.Pokemon
	getOtherSlots(slot *slot) []*slot
	injectReplaceAction(slot *slot, trainer *trainer)
	getTrainer(slot *slot) *trainer
	gatherActions()
	execute()
}

type slot struct {
	mon *pokemon.Pokemon
}
