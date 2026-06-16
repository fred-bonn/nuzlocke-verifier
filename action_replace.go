package main

import "log"

type replaceAction struct {
	oldSlot *slot
	trainer *trainer
	midTurn bool
}

func (ra *replaceAction) prio() int {
	if ra.midTurn {
		return 10
	}
	return -10
}

func (ra *replaceAction) speed(bs battleState) int {
	return ra.oldSlot.mon.effectiveSpeed(bs)
}

func (ra *replaceAction) invoke(bs battleState) {
	mon := ra.trainer.selectSwitchIn(bs, ra.oldSlot)
	if mon == nil {
		return
	}

	for _, slot := range bs.getOtherSlots(ra.oldSlot) {
		if ailment := slot.mon.hasAilment("infatuation"); ailment != nil && ailment.AfflictedBy == ra.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
		if ailment := slot.mon.hasAilment("trap"); ailment != nil && ailment.AfflictedBy == ra.oldSlot {
			delete(slot.mon.Ailments, "infatuation")
		}
	}
	if f, ok := onSwitchAbilities[ra.oldSlot.mon.Ability]; ok {
		f(ra.oldSlot, bs, false)
	}

	ra.oldSlot.setMon(mon)
	log.Printf("%s was sent out", mon.Base.Name)
	if f, ok := onSwitchAbilities[ra.oldSlot.mon.Ability]; ok {
		f(ra.oldSlot, bs, true)
	}
}
