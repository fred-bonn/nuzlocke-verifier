package main

import "log"

type replaceAction struct {
	oldSlot *slot
	trainer *trainer
}

func (ra *replaceAction) prio() int {
	return -10
}

func (ra *replaceAction) speed() int {
	return ra.oldSlot.mon.EffectiveStat("speed")
}

func (ra *replaceAction) invoke(sbs battleState) {
	mon := ra.trainer.selectSwitchIn(sbs, ra.oldSlot)
	if mon == nil {
		return
	}

	sbs.setMon(ra.oldSlot.mon, mon)
	log.Printf("%s was sent out", mon.Base.Name)
}
