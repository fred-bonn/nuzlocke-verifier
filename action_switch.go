package main

import (
	"log"
)

type switchAction struct {
	oldSlot *slot
	new     *Pokemon
}

func (sa *switchAction) invoke(bs battleState) {
	log.Printf("switched %s for %s", sa.oldSlot.mon.Base.Name, sa.new.Base.Name)
	for _, slot := range bs.getOtherSlots(sa.oldSlot) {
		if ailment := slot.mon.hasAilment("infatuation"); ailment != nil && ailment.AfflictedBy == sa.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
		if ailment := slot.mon.hasAilment("trap"); ailment != nil && ailment.AfflictedBy == sa.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
		if slot.trainer != sa.oldSlot.trainer {
			log.Printf("test")
			slot.mon.Unnerved = sa.new.Ability == "unnerve"
			slot.mon.Item.checkTrigger(true, nil)
		}
	}
	sa.oldSlot.setMon(sa.new)
}

func (sa *switchAction) prio() int {
	return 10
}

func (sa *switchAction) speed() int {
	return sa.oldSlot.mon.effectiveSpeed()
}
