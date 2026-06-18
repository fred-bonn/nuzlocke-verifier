package main

import (
	"log"
)

type switchAction struct {
	oldSlot *slot
	new     *Pokemon
}

func (sa *switchAction) invoke(bs battleState) {
	for _, slot := range bs.getOtherSlots(sa.oldSlot) {
		if ailment := slot.mon.hasAilment("infatuation"); ailment != nil && ailment.AfflictedBy == sa.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
		if ailment := slot.mon.hasAilment("trap"); ailment != nil && ailment.AfflictedBy == sa.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
	}
	if f, ok := onSwitchAbilities[sa.oldSlot.mon.Ability]; ok {
		f(sa.oldSlot, bs, false)
	}

	if a, ok := bs.getActions().queue.fetchBy(fetchPursuitMiddleware(sa.oldSlot.mon.Base.Name)); ok {
		p, _ := a.(*moveAction)
		p.pursuit = true
		p.invoke(bs)
	}

	log.Printf("switched %s for %s", sa.oldSlot.mon.Base.Name, sa.new.Base.Name)
	sa.oldSlot.setMon(bs, sa.new)
	if f, ok := onSwitchAbilities[sa.oldSlot.mon.Ability]; ok {
		f(sa.oldSlot, bs, true)
	}
}

func (sa *switchAction) prio() int {
	return 10
}

func (sa *switchAction) speed(bs battleState) int {
	return sa.oldSlot.mon.effectiveSpeed(bs)
}
