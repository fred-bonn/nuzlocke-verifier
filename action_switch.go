package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type switchAction struct {
	old *pokemon.Pokemon
	new *pokemon.Pokemon
}

func (sa *switchAction) invoke(bs battleState) {
	bs.setMon(sa.old, sa.new)
	sa.old.ResetStages()
	log.Printf("switched %s for %s", sa.old.Base.Name, sa.new.Base.Name)
}

func (sa *switchAction) prio() int {
	return 10
}

func (sa *switchAction) speed() int {
	return sa.old.EffectiveStat("speed")
}
