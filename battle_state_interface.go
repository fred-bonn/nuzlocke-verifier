package main

import (
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
	setWeather(weatherState)
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
			switch ailment.state {
			case Burn:
				takeResidualDamage(bs, slot, ailment.state.String(), 1, 16)
			case Poison:
				takeResidualDamage(bs, slot, ailment.state.String(), 1, 8)
			case Toxic:
				ailment.turns++
				takeResidualDamage(bs, slot, ailment.state.String(), ailment.turns, 16)
			case Trap:
				ailment.turns--
				takeResidualDamage(bs, slot, ailment.state.String(), 1, 8)
				if ailment.turns <= 0 {
					vlogf("%s was freed", slot.mon.base.Name)
					delete(slot.mon.ailments, ailment.state)
				}
			case LeechSeed:
				vlogf("%s leeched health from %s", ailment.afflictedBy.mon.base.Name, slot.mon.base.Name)
				dmg := takeResidualDamage(bs, slot, ailment.state.String(), 1, 8)
				ailment.afflictedBy.mon.changeHpBy(dmg)
			case Yawn:
				ailment.turns--
				if ailment.turns == 0 {
					slot.mon.applyAilment(Sleep, nil, ailment.afflictedBy)
					delete(slot.mon.ailments, ailment.state)
				}
			}
		}

		// resolve end of turn effects of weather
		if w := bs.getWeather(); w != NoneWeather {
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
			vlogf("%s harvested its %s", slot.mon.base.Name, slot.mon.item.name)
			slot.mon.item.consumed = false
			slot.mon.checkItemTrigger(true, nil)
		} else if slot.mon.ability == "speed-boost" && !slot.firstTurn {
			slot.mon.changeStatStageBy(Speed, 1, false)
		}

		if slot.mon.item.name == "leftovers" {
			change := slot.mon.maxHP() / 16
			vlogItem("%s restored %d health from leftovers", slot.mon.base.Name, change)
			slot.mon.changeHpBy(change)
		}

		slot.mon.laserFocus = false
	}
}

func takeResidualDamage(bs battleState, slot *slot, effect string, num, den int) int {
	if slot.mon.fainted {
		return 0
	}

	change := slot.mon.maxHP() * num / den
	vlogf("%s took %d damage from %s", slot.mon.base.Name, change, effect)
	slot.mon.changeHpBy(-change)
	if slot.mon.hp <= 0 {
		slot.mon.fainted = true
		injectReplaceAction(bs, slot, false)
		vlogf("%s fainted!", slot.mon.base.Name)
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
