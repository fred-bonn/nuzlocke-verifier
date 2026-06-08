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

func (ra *replaceAction) speed() int {
	return ra.oldSlot.mon.effectiveSpeed()
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
		if slot.trainer != ra.oldSlot.trainer {
			slot.mon.Unnerved = mon.Ability == "unnerve"
			slot.mon.checkItemTrigger(true, nil)
		}
	}
	ra.oldSlot.setMon(mon)
	log.Printf("%s was sent out", mon.Base.Name)
}
