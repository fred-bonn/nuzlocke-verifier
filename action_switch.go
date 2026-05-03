package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type switchAction struct {
	oldSlot *slot
	new     *pokemon.Pokemon
}

func (sa *switchAction) invoke(bs battleState) {
	log.Printf("switched %s for %s", sa.oldSlot.mon.Base.Name, sa.new.Base.Name)
	bs.setMon(sa.oldSlot.mon, sa.new)
	sa.oldSlot.mon.ResetStages()
}

func (sa *switchAction) prio() int {
	return 10
}

func (sa *switchAction) speed() int {
	return sa.oldSlot.mon.EffectiveStat("speed", false)
}
