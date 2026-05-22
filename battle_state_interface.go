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
	mon       *Pokemon
	firstTurn bool
}

func (s *slot) setMon(new *Pokemon) {
	s.mon.SwitchReset()
	s.firstTurn = true
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

func resolveEndOfTurn(bs battleState) {
	for _, slot := range bs.getAllSlots() {
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
