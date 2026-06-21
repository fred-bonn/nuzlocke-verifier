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
	getWeather() weatherState
	setWeather(weatherState, int)
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
		for _, ailment := range slot.mon.ailments {
			switch ailment.name {
			case "burn":
				takeResidualDamage(bs, slot, ailment.name, 1, 16)
			case "poison":
				takeResidualDamage(bs, slot, ailment.name, 1, 8)
			case "toxic":
				ailment.turns++
				takeResidualDamage(bs, slot, ailment.name, ailment.turns, 16)
			case "trap":
				ailment.turns--
				takeResidualDamage(bs, slot, ailment.name, 1, 8)
				if ailment.turns <= 0 {
					log.Printf("%s was freed", slot.mon.base.Name)
					delete(slot.mon.ailments, ailment.name)
				}
			case "leech-seed":
				log.Printf("%s leeched health from %s", ailment.afflictedBy.mon.base.Name, slot.mon.base.Name)
				dmg := takeResidualDamage(bs, slot, ailment.name, 1, 8)
				ailment.afflictedBy.mon.changeHpBy(dmg)
			case "yawn":
				ailment.turns--
				if ailment.turns == 0 {
					slot.mon.applyAilment("sleep", nil, ailment.afflictedBy)
					delete(slot.mon.ailments, "yawn")
				}
			}
		}

		// resolve end of turn effects of weather
		w := bs.getWeather()
		if w != None {
			if w.affectsMon(slot.mon) {
				takeResidualDamage(bs, slot, w.String(), 1, 16)
			}
			w.activateMonAbility(bs, slot)
		}

		// reset protect counter if the slot was not protected this turn
		if !slot.protected {
			slot.protectTurns = 0
		} else {
			slot.protected = false
		}

		if slot.mon.ability == "harvest" && roll(1, 2) && strings.HasSuffix(slot.mon.item.name, "berry") {
			log.Printf("%s harvested its %s", slot.mon.base.Name, slot.mon.item.name)
			slot.mon.item.consumed = false
			slot.mon.checkItemTrigger(true, nil)
		} else if slot.mon.ability == "speed-boost" && !slot.firstTurn {
			slot.mon.changeStatStageBy("speed", 1, false)
		}

		if slot.mon.item.name == "leftovers" {
			log.Printf("%s restored health from leftovers", slot.mon.base.Name)
			slot.mon.changeHpBy(slot.mon.maxHP() / 16)
		}

		slot.mon.laserFocus = false
	}
}

func takeResidualDamage(bs battleState, slot *slot, effect string, num, den int) int {
	if slot.mon.fainted {
		return 0
	}

	log.Printf("%s took damage from %s", slot.mon.base.Name, effect)
	change := slot.mon.maxHP() * num / den
	slot.mon.changeHpBy(-change)
	if slot.mon.hp <= 0 {
		slot.mon.fainted = true
		injectReplaceAction(bs, slot, false)
		log.Printf("%s fainted!", slot.mon.base.Name)
	}
	return change
}

func resolveOnEntry(bs battleState) {
	for _, slot := range bs.getAllSlots() {
		if f, ok := onSwitchAbilities[slot.mon.ability]; ok {
			f(slot, bs, true)
		}
	}
}
