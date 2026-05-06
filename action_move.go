package main

import (
	"log"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type moveAction struct {
	userSlot   *slot
	targetSlot *slot
	move       *pokeapi.BaseMove
}

func (ma *moveAction) prio() int {
	return ma.move.Priority
}

func (ma *moveAction) speed() int {
	speed := ma.userSlot.mon.EffectiveStat("speed", false)
	if _, ok := ma.userSlot.mon.Ailments["paralysis"]; ok {
		return speed / 4
	}
	return speed
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.Fainted {
		return
	}

	if turns, ok := ma.userSlot.mon.Ailments["confusion"]; ok {
		if turns > 0 {
			ma.userSlot.mon.Ailments["confusion"] -= 1
			log.Printf("%s is confused", ma.userSlot.mon.Base.Name)
			if roll(1.0 / 3.0) {
				damage := calculateDamage(ma.userSlot.mon, ma.userSlot.mon, &confusionMove, false, false)
				log.Printf("%s hit itself in confusion for %d damage", ma.userSlot.mon.Base.Name, damage)
				ma.userSlot.mon.Hp -= int(damage)
				if ma.userSlot.mon.Hp <= 0 {
					ma.monFainted(bs, ma.userSlot)
				}
				return
			}
		} else {
			delete(ma.userSlot.mon.Ailments, "confusion")
			log.Printf("%s snapped out of confusion", ma.userSlot.mon.Base.Name)
		}
	}

	target := bs.getMon(ma.targetSlot)
	hitChance := target.EffectiveEvasion() * ma.userSlot.mon.EffectiveAccuracy()
	if ma.move.Accuracy != 0 && !roll(float32(ma.move.Accuracy)/100.0*hitChance) {
		log.Printf("%s's move %s missed", ma.userSlot.mon.Base.Name, ma.move.Name)
		return
	}

	log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			ma.applyStatusMove(bs, ma.userSlot.mon)
		} else {
			ma.applyStatusMove(bs, target)
		}
	} else {
		ma.applyDamageMove(bs)
	}
}

func (ma *moveAction) applyStatusMove(bs battleState, target *pokemon.Pokemon) {
	if ma.move.Category == "swagger" {
		ma.applySwagger(bs, target)
		return
	}

	if ma.move.StatChance != 100 && !roll(float32(ma.move.StatChance/100)) {
		log.Printf("missed")
		return
	}

	for stat, change := range ma.move.StatChanges {
		target.Stages[stat] = max(-6, min(6, target.Stages[stat]+change))
		log.Printf("%s's %s changed by %d stages (%d)", target.Base.Name, stat, change, target.Stages[stat])
	}
}

func (ma *moveAction) applySwagger(bs battleState, target *pokemon.Pokemon) {
	err := target.ChangeStatStage("attack", 2)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = target.ApplyAilment("confusion")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("%s's attack changed by 2 stages (%d)", target.Base.Name, target.Stages["attack"])
	log.Printf("%s became confused", target.Base.Name)
}

func (ma *moveAction) applyDamageMove(bs battleState) {
	crit := roll(1.0 / critRateMap[ma.move.CritRate])
	damage := calculateDamage(ma.userSlot.mon, ma.targetSlot.mon, ma.move, crit, false)
	target := ma.targetSlot.mon

	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if crit {
		log.Printf("it was a critical hit!")
	}

	target.Hp -= int(damage)
	if target.Hp <= 0 {
		ma.monFainted(bs, ma.targetSlot)
	}

	if _, ok := pivotMoves[ma.move.Name]; ok {
		trainer := bs.getTrainer(ma.userSlot)
		if trainer.canReplace(bs) {
			bs.injectReplaceAction(ma.userSlot, trainer, true)
		}
	}
}

func (ma *moveAction) monFainted(bs battleState, slot *slot) {
	slot.mon.Fainted = true
	bs.injectReplaceAction(slot, bs.getTrainer(slot), false)
	log.Printf("%s fainted!", slot.mon.Base.Name)
}
