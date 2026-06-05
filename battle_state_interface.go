package main

import (
	"log"
	"strings"
)

type battleState interface {
	getAllSlots() []*slot
	getOtherSlots(slot *slot) []*slot
	injectReplaceAction(slot *slot, trainer *trainer, midTurn bool)
	getTrainer(slot *slot) *trainer
	gatherActions()
	getActions() *ActionQueue
	execute()
}

type slot struct {
	mon                *Pokemon
	trainer            *trainer
	firstTurn          bool
	suckerPunch        bool
	protected          bool
	protectTurns       int
	invulnerableAction *moveAction
	unnerved           bool
}

func (s *slot) setMon(new *Pokemon) {
	s.mon.switchReset()
	s.firstTurn = true
	s.suckerPunch = false
	s.mon = new
}

func (s *slot) isTrapped() bool {
	return s.mon.hasAilment("trap") != nil || s.mon.hasAilment("bound") != nil
}

func (s *slot) resolveProtect() {
	denominator := 1
	for i := 0; i < s.protectTurns; i++ {
		denominator *= 3
	}
	if roll(1, denominator) {
		s.protected = true
		s.protectTurns++
	} else {
		log.Println("but it failed")
	}
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
			slot.mon.Item.checkTrigger(true, nil)
		} else if slot.mon.Ability == "speed-boost" && !slot.firstTurn {
			slot.mon.changeStatStageBy("speed", 1)
		}
	}
}

func takeResidualDamage(bs battleState, slot *slot, ailment *Ailment, num, den int) int {
	if slot.mon.Fainted {
		return 0
	}

	log.Printf("%s took damage from %s", slot.mon.Base.Name, ailment.Name)
	change := slot.mon.Stats["hp"] * num / den
	slot.mon.changeHpBy(-change)
	if slot.mon.Hp <= 0 {
		slot.mon.Fainted = true
		bs.injectReplaceAction(slot, bs.getTrainer(slot), false)
		log.Printf("%s fainted!", slot.mon.Base.Name)
	}
	return change
}
