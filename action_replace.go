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

	if ra.midTurn {
		f := func(a action) bool {
			ma, ok := a.(*moveAction)
			if !ok {
				return false
			}
			if ma.move.Name != "pursuit" {
				return false
			}
			if ma.targetSlot.mon.Base.Name != ra.oldSlot.mon.Base.Name {
				return false
			}
			return true
		}

		if a, ok := bs.getActions().queue.fetchBy(f); ok {
			p, _ := a.(*moveAction)
			if p.targetSlot.mon.Base.Name == ra.oldSlot.mon.Base.Name {
				p.pursuit = true
				p.invoke(bs)
				if ra.oldSlot.mon.Fainted {
					injectReplaceAction(bs, ra.oldSlot, false)
					return
				}
			}
		}
	}

	ra.oldSlot.setMon(mon)
	log.Printf("%s was sent out", mon.Base.Name)
	if f, ok := onSwitchAbilities[ra.oldSlot.mon.Ability]; ok {
		f(ra.oldSlot, bs, true)
	}
}
