package main

import (
	"log"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
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
	return ma.userSlot.mon.effectiveSpeed()
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.Fainted {
		return
	}

	ma.userSlot.firstTurn = false

	ma.userSlot.suckerPunch = ma.move.Name == "sucker-punch"

	if _, ok := ma.userSlot.mon.Ailments["freeze"]; ok {
		if roll(1, 5) {
			log.Printf("%s thawed out", ma.userSlot.mon.Base.Name)
			delete(ma.userSlot.mon.Ailments, "freeze")
		} else {
			log.Printf("%s is frozen", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if turns, ok := ma.userSlot.mon.Ailments["sleep"]; ok {
		if turns <= 0 {
			log.Printf("%s woke up", ma.userSlot.mon.Base.Name)
			delete(ma.userSlot.mon.Ailments, "sleep")
		} else {
			ma.userSlot.mon.Ailments["sleep"] -= 1
			log.Printf("%s is asleep", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if _, ok := ma.userSlot.mon.Ailments["paralysis"]; ok {
		if roll(1, 4) {
			log.Printf("%s is paralysed", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if ma.Flinch {
		log.Printf("%s flinched", ma.userSlot.mon.Base.Name)
		return
	}

	ma.move.PP--

	if turns, ok := ma.userSlot.mon.Ailments["confusion"]; ok {
		if turns > 0 {
			ma.userSlot.mon.Ailments["confusion"] -= 1
			log.Printf("%s is confused", ma.userSlot.mon.Base.Name)
			if roll(1, 3) {
				damage := calculateDamage(ma.userSlot.mon, ma.userSlot.mon, &confusionMove, false, false)
				log.Printf("%s hit itself in confusion for %d damage", ma.userSlot.mon.Base.Name, damage)
				ma.userSlot.mon.Hp -= int(damage)
				if ma.userSlot.mon.Hp <= 0 {
					monFainted(bs, ma.userSlot)
				}
				return
			}
		} else {
			delete(ma.userSlot.mon.Ailments, "confusion")
			log.Printf("%s snapped out of confusion", ma.userSlot.mon.Base.Name)
		}
	}

	if _, ok := multipleTurnMoves[ma.move.Name]; ok {
		if ma.userSlot.invulnerableAction == nil {
			ma.move.PP++
			ma.userSlot.invulnerableAction = ma
			log.Printf("%s used %s and became invulnerable", ma.userSlot.mon.Base.Name, ma.move.Name)
			return
		}
		ma.userSlot.invulnerableAction = nil
	}

	if ma.userSlot.suckerPunch {
		targetMove := bs.getActions().getMoveActionBy(ma.targetSlot.mon)
		if targetMove == nil || targetMove.move.Class == "status" {
			log.Printf("%s used sucker punch but it failed", ma.userSlot.mon.Base.Name)
			return
		}
	}

	target := ma.targetSlot.mon
	if ma.move.Accuracy > 0 && !accuracyRoll(ma.userSlot.mon, target, ma.move.Accuracy) {
		log.Printf("%s's move %s missed", ma.userSlot.mon.Base.Name, ma.move.Name)
		return
	}

	log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)

	if ma.move.Name == "struggle" {
		ma.userSlot.mon.changeHp(-(ma.userSlot.mon.Stats["hp"] / 4))
	}

	if ma.targetSlot.protected || ma.targetSlot.invulnerableAction != nil {
		log.Printf("but it failed")
		return
	}

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			ma.applyStatusMove(bs, ma.userSlot.mon)
		} else {
			ma.applyStatusMove(bs, target)
		}
	} else {
		ma.applyDamageMove(bs)
	}

	if ma.userSlot.mon.Hp <= 0 {
		monFainted(bs, ma.userSlot)
		return
	}

	ma.userSlot.mon.Item.checkTrigger(true, leppaBerryEvent{
		move: ma.move,
	})
}

func (ma *moveAction) applyStatusMove(bs battleState, target *Pokemon) {

	if _, ok := protectMoves[ma.move.Name]; ok {
		ma.userSlot.resolveProtect()
	}

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
		ma.userSlot.mon.changeHp(change)
		log.Printf("%s healed for %d", target.Base.Name, change)
	}

	for stat, change := range ma.move.StatChanges {
		target.Stages[stat] = max(-6, min(6, target.Stages[stat]+change))
		log.Printf("%s's %s changed by %d stages (%d)", target.Base.Name, stat, change, target.Stages[stat])
	}
}

func (ma *moveAction) applySwagger(target *Pokemon) {
	target.changeStatStage("attack", 2)
	log.Printf("%s's attack changed by 2 stages (%d)", target.Base.Name, target.Stages["attack"])
	target.applyAilment("confusion", ma.move)
}

func (ma *moveAction) applyDamageMove(bs battleState) {
	target := ma.targetSlot.mon
	if target.Fainted {
		return
	}

	hits := 1
	if ma.move.MinHits > 0 {
		hits = determineHits(ma.move)
	}

	if ok := ma.resolveDamage(bs); !ok {
		return
	}

	for i := 1; i < hits; i++ {
		if ok := ma.resolveDamage(bs); !ok {
			return
		}
	}
}

func (ma *moveAction) resolveDamage(bs battleState) bool {
	user := ma.userSlot.mon
	target := ma.targetSlot.mon

	crit := roll(1, critRateMap[ma.move.CritRate])

	damage := calculateDamage(user, target, ma.move, crit, false)
	if damage == 0 {
		log.Printf("it does not affect %s", target.Base.Name)
		return false
	}

	target.Item.checkTrigger(true, resistBerryEvent{
		typeName: ma.move.Type,
	})
	target.Item.checkTrigger(true, focusSashEvent{
		damage:  &damage,
		consume: true,
	})
	user.Item.checkTrigger(true, gemEvent{
		typeName: ma.move.Type,
	})

	damage = min(damage, target.Hp)
	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if crit {
		log.Printf("it was a critical hit!")
	}
	target.changeHp(-damage)
	if target.Hp <= 0 {
		monFainted(bs, ma.targetSlot)
	}

	if ma.move.Drain != 0 {
		change := damage * ma.move.Drain / 100
		if change == 0 {
			if ma.move.Drain > 0 {
				change = 1
			}
			change = -1
		}

		user.changeHp(change)
		if change >= 0 {
			log.Printf("%s healed for %d", user.Base.Name, change)
		} else {
			log.Printf("%s took recoil for %d", user.Base.Name, -change)
			if user.Hp <= 0 {
				monFainted(bs, ma.userSlot)
			}
		}
	}

	if ma.move.AilmentChance > 0 && !target.Fainted && roll(ma.move.AilmentChance, 100) {
		target.applyAilment(ma.move.Ailment, ma.move)
	}

	if ma.move.FlinchChance > 0 && !target.Fainted && roll(ma.move.FlinchChance, 100) {
		targetMove := bs.getActions().getMoveActionBy(target)
		if targetMove != nil {
			targetMove.Flinch = true
		}
	}

	if ma.move.StatChance > 0 && roll(ma.move.StatChance, 100) {
		var mon *Pokemon
		switch ma.move.Category {
		case "damage-raise":
			mon = user
		case "damage-lower":
			mon = target
		}

		for stat, change := range ma.move.StatChanges {
			mon.Stages[stat] = max(-6, min(6, mon.Stages[stat]+change))
			log.Printf("%s's %s changed by %d stages (%d)", mon.Base.Name, stat, change, mon.Stages[stat])
		}
	}

	if _, ok := pivotMoves[ma.move.Name]; ok {
		if trainer := bs.getTrainer(ma.userSlot); trainer.canReplace(bs) {
			bs.injectReplaceAction(ma.userSlot, trainer, true)
		}
	}

	if target.Fainted {
		return false
	}

	return true
}
