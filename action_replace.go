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
	return ra.oldSlot.mon.EffectiveStat("speed", false)
}

func (ra *replaceAction) invoke(bs battleState) {
	mon := ra.trainer.selectSwitchIn(bs, ra.oldSlot)
	if mon == nil {
		return
	}

	ra.oldSlot.setMon(mon)
	log.Printf("%s was sent out", mon.Base.Name)
}
