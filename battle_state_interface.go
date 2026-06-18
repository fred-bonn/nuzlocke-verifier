package main

import (
	"log"
	"strings"
)

type battleState interface {
	execute()
	gatherActions()
	getAllSlots() []*slot
	getOtherSlots(slot *slot) []*slot
	getOpponentSlot(slot *slot) *slot
	getActions() *actionQueue
}

func injectReplaceAction(bs battleState, slot *slot, midTurn bool) {
	bs.getActions().queue.push(&replaceAction{
		oldSlot: slot,
		trainer: slot.trainer,
		midTurn: midTurn,
	})
	bs.getActions().sort(bs)
}

func resolveEndOfTurn(bs battleState) {
	for _, slot := range bs.getAllSlots() {
		// resolve end of return effects from ailments and statuses
		for _, ailment := range slot.mon.Ailments {
			switch ailment.Name {
			case "burn":
				takeResidualDamage(bs, slot, ailment, 1, 16)
			case "poison":
				takeResidualDamage(bs, slot, ailment, 1, 8)
			case "toxic":
				ailment.Turns++
				takeResidualDamage(bs, slot, ailment, ailment.Turns, 16)
			case "trap":
				ailment.Turns--
				takeResidualDamage(bs, slot, ailment, 1, 8)
				if ailment.Turns <= 0 {
					log.Printf("%s was freed", slot.mon.Base.Name)
					delete(slot.mon.Ailments, ailment.Name)
				}
			case "leech-seed":
				log.Printf("%s leeched health from %s", ailment.AfflictedBy.mon.Base.Name, slot.mon.Base.Name)
				dmg := takeResidualDamage(bs, slot, ailment, 1, 8)
				ailment.AfflictedBy.mon.changeHpBy(dmg)
			case "yawn":
				ailment.Turns--
				if ailment.Turns == 0 {
					slot.mon.applyAilment("sleep", nil, ailment.AfflictedBy)
					delete(slot.mon.Ailments, "yawn")
				}
			}
		}

		// reset protect counter if the slot was not protected this turn
		if !slot.protected {
			slot.protectTurns = 0
		} else {
			slot.protected = false
		}

		if slot.mon.Ability == "harvest" && roll(1, 2) && strings.HasSuffix(slot.mon.Item.name, "berry") {
			log.Printf("%s harvested its %s", slot.mon.Base.Name, slot.mon.Item.name)
			slot.mon.Item.consumed = false
			slot.mon.checkItemTrigger(true, nil)
		} else if slot.mon.Ability == "speed-boost" && !slot.firstTurn {
			slot.mon.changeStatStageBy("speed", 1, false)
		}

		if slot.mon.Item.name == "leftovers" {
			log.Printf("%s restored health from leftovers", slot.mon.Base.Name)
			slot.mon.changeHpBy(slot.mon.maxHP() / 16)
		}

		slot.mon.LaserFocus = false
	}
}

func takeResidualDamage(bs battleState, slot *slot, ailment *Ailment, num, den int) int {
	if slot.mon.Fainted {
		return 0
	}

	log.Printf("%s took damage from %s", slot.mon.Base.Name, ailment.Name)
	change := slot.mon.maxHP() * num / den
	slot.mon.changeHpBy(-change)
	if slot.mon.Hp <= 0 {
		slot.mon.Fainted = true
		injectReplaceAction(bs, slot, false)
		log.Printf("%s fainted!", slot.mon.Base.Name)
	}
	return change
}

func resolveOnEntry(bs battleState) {
	for _, slot := range bs.getAllSlots() {
		if f, ok := onSwitchAbilities[slot.mon.Ability]; ok {
			f(slot, bs, true)
		}
	}
}
