package main

import (
	"log"
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
	mon          *Pokemon
	firstTurn    bool
	suckerPunch  bool
	protected    bool
	protectTurns int
}

func (s *slot) setMon(new *Pokemon) {
	s.mon.SwitchReset()
	s.firstTurn = true
	s.suckerPunch = false
	s.mon.grounded = false
	s.mon = new
}

func (s *slot) isTrapped() bool {
	if _, ok := s.mon.Ailments["bound"]; ok {
		return true
	}
	if _, ok := s.mon.Ailments["trap"]; ok {
		return true
	}
	return false
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
		for ailment := range slot.mon.Ailments {
			switch ailment {
			case "burn":
				takeResidualDamage(bs, slot, ailment, 1, 16)
			case "poison":
				takeResidualDamage(bs, slot, ailment, 1, 8)
			case "toxic":
				slot.mon.Ailments[ailment]++
				takeResidualDamage(bs, slot, ailment, slot.mon.Ailments[ailment], 16)
			case "trap":
				slot.mon.Ailments[ailment]--
				takeResidualDamage(bs, slot, ailment, 1, 8)
				if slot.mon.Ailments[ailment] == 0 {
					log.Printf("%s was freed", slot.mon.Base.Name)
					delete(slot.mon.Ailments, ailment)
				}
			}
		}

		// reset protect counter if the slot was not protected this turn
		if !slot.protected {
			slot.protectTurns = 0
		} else {
			slot.protected = false
		}
	}
}

func takeResidualDamage(bs battleState, slot *slot, ailment string, num, den int) {
	if slot.mon.Fainted {
		return
	}

	log.Printf("%s took damage from %s", slot.mon.Base.Name, ailment)
	change := slot.mon.Stats["hp"] * num / den
	slot.mon.ChangeHp(-change)
	if slot.mon.Hp <= 0 {
		slot.mon.Fainted = true
		bs.injectReplaceAction(slot, bs.getTrainer(slot), false)
		log.Printf("%s fainted!", slot.mon.Base.Name)
	}
}
