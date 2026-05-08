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
	Flinch     bool
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

	if ma.Flinch {
		log.Printf("%s flinched", ma.userSlot.mon.Base.Name)
		return
	}

	if turns, ok := ma.userSlot.mon.Ailments["confusion"]; ok {
		if turns > 0 {
			ma.userSlot.mon.Ailments["confusion"] -= 1
			log.Printf("%s is confused", ma.userSlot.mon.Base.Name)
			if roll(1, 3) {
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
	if ma.move.Accuracy > 0 && !accuracyRoll(ma.userSlot.mon, target, ma.move.Accuracy) {
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
		ma.applySwagger(target)
		return
	}

	if ma.move.StatChance != 100 && !roll(ma.move.StatChance, 100) {
		log.Printf("missed")
		return
	}

	if ma.move.Heal > 0 {
		change := target.Stats["hp"] * ma.move.Heal / 100
		ma.userSlot.mon.ChangeHp(change)
		log.Printf("%s healed for %d", target.Base.Name, change)
	}

	for stat, change := range ma.move.StatChanges {
		target.Stages[stat] = max(-6, min(6, target.Stages[stat]+change))
		log.Printf("%s's %s changed by %d stages (%d)", target.Base.Name, stat, change, target.Stages[stat])
	}
}

func (ma *moveAction) applySwagger(target *pokemon.Pokemon) {
	target.ChangeStatStage("attack", 2)
	log.Printf("%s's attack changed by 2 stages (%d)", target.Base.Name, target.Stages["attack"])
	if ok := target.ApplyAilment("confusion"); ok {
		log.Printf("%s became afflicted with confusion", target.Base.Name)
	}
}

func (ma *moveAction) applyDamageMove(bs battleState) {
	target := ma.targetSlot.mon
	if target.Fainted {
		return
	}

	critRate := ma.move.CritRate
	crit := roll(1, critRateMap[critRate])

	damage := calculateDamage(ma.userSlot.mon, ma.targetSlot.mon, ma.move, crit, false)
	if damage == 0 {
		log.Printf("it does not affect %s", target.Base.Name)
	}
	damage = min(damage, target.Hp)

	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if crit {
		log.Printf("it was a critical hit!")
	}

	target.ChangeHp(-damage)
	if target.Hp <= 0 {
		ma.monFainted(bs, ma.targetSlot)
	}

	if ma.move.Drain != 0 {
		change := max(1, damage*ma.move.Drain/100)
		ma.userSlot.mon.ChangeHp(change)
		if change >= 0 {
			log.Printf("%s healed for %d", ma.userSlot.mon.Base.Name, change)
		} else {
			log.Printf("%s took recoil for for %d", ma.userSlot.mon.Base.Name, -change)
			if ma.userSlot.mon.Hp <= 0 {
				ma.monFainted(bs, ma.userSlot)
			}
		}
	}

	if ma.move.AilmentChance > 0 && !target.Fainted && roll(ma.move.AilmentChance, 100) {
		if ok := target.ApplyAilment(ma.move.Ailment); ok {
			log.Printf("%s became afflicted with %s", target.Base.Name, ma.move.Ailment)
		}
	}

	if ma.move.FlinchChance > 0 && !target.Fainted && roll(ma.move.FlinchChance, 100) {
		target_ma := bs.getActions().getMoveActionBy(target)
		if target_ma != nil {
			target_ma.Flinch = true
		}
	}

	if _, ok := pivotMoves[ma.move.Name]; ok {
		if trainer := bs.getTrainer(ma.userSlot); trainer.canReplace(bs) {
			bs.injectReplaceAction(ma.userSlot, trainer, true)
		}
	}
}

func (ma *moveAction) monFainted(bs battleState, slot *slot) {
	slot.mon.Fainted = true
	bs.injectReplaceAction(slot, bs.getTrainer(slot), false)
	log.Printf("%s fainted!", slot.mon.Base.Name)
}
