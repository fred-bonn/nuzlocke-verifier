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
	flinch     bool
	pursuit    bool
}

func (ma *moveAction) prio() int {
	bonus := 0
	if ma.userSlot.mon.ability == "prankster" && ma.move.Class == "status" {
		bonus++
	}

	return ma.move.Priority + bonus
}

func (ma *moveAction) speed(bs battleState) int {
	return ma.userSlot.mon.effectiveSpeed(bs)
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.fainted {
		return
	}

	if (ma.move.Name == "fake-out" || ma.move.Name == "first-impression") && !ma.userSlot.firstTurn {
		log.Printf("%s used %s", ma.userSlot.mon.base.Name, ma.move.Name)
		log.Printf("but it failed")
		return
	}

	ma.userSlot.firstTurn = false

	ma.userSlot.suckerPunch = ma.move.Name == "sucker-punch"

	if _, ok := ma.userSlot.mon.ailments["freeze"]; ok {
		if _, ok := selfThawingMoves[ma.move.Name]; ok || roll(1, 5) {
			log.Printf("%s thawed out", ma.userSlot.mon.base.Name)
			delete(ma.userSlot.mon.ailments, "freeze")
		} else {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is frozen", ma.userSlot.mon.base.Name)
			return
		}
	}

	if sleep := ma.userSlot.mon.hasAilment("sleep"); sleep != nil {
		if sleep.turns <= 0 {
			log.Printf("%s woke up", ma.userSlot.mon.base.Name)
			delete(ma.userSlot.mon.ailments, "sleep")
		} else {
			if ma.userSlot.mon.ability == "early-bird" {
				sleep.turns -= 2
			} else {
				sleep.turns--
			}
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is asleep", ma.userSlot.mon.base.Name)
			return
		}
	}

	if confusion := ma.userSlot.mon.hasAilment("confusion"); confusion != nil {
		if confusion.turns > 0 {
			confusion.turns -= 1
			log.Printf("%s is confused", ma.userSlot.mon.base.Name)
			if roll(1, 3) {
				damage := calculateDamage(ma.userSlot.mon, ma.userSlot.mon, &confusionMove, new(false), false, false, false)
				ma.userSlot.invulnerableAction = nil
				log.Printf("%s hit itself in confusion for %d damage", ma.userSlot.mon.base.Name, damage)
				ma.userSlot.mon.hp -= int(damage)
				if ma.userSlot.mon.hp <= 0 {
					monFainted(bs, ma.userSlot, false)
				}
				return
			}
		} else {
			delete(ma.userSlot.mon.ailments, "confusion")
			log.Printf("%s snapped out of confusion", ma.userSlot.mon.base.Name)
		}
	}

	if paralysis := ma.userSlot.mon.hasAilment("paralysis"); paralysis != nil {
		if roll(1, 4) {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is paralysed", ma.userSlot.mon.base.Name)
			return
		}
	}

	if infatuation := ma.userSlot.mon.hasAilment("infatuation"); infatuation != nil {
		if roll(1, 2) {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is infatuated with %s", ma.userSlot.mon.base.Name, infatuation.afflictedBy.mon.base.Name)
			return
		}
	}

	if ma.flinch {
		log.Printf("%s flinched", ma.userSlot.mon.base.Name)
		return
	}

	ma.move.PP--

	if _, ok := multipleTurnMoves[ma.move.Name]; ok {
		if ma.userSlot.invulnerableAction == nil {
			ma.move.PP++
			ma.userSlot.invulnerableAction = ma
			log.Printf("%s used %s and became invulnerable", ma.userSlot.mon.base.Name, ma.move.Name)
			return
		}
		ma.userSlot.invulnerableAction = nil
	}

	if ma.userSlot.suckerPunch {
		targetMove := bs.getActions().getMoveActionBy(ma.targetSlot.mon)
		if targetMove == nil || targetMove.move.Class == "status" {
			log.Printf("%s used sucker punch but it failed", ma.userSlot.mon.base.Name)
			return
		}
	}

	target := ma.targetSlot.mon
	if ma.move.Accuracy > 0 && !ma.pursuit && !accuracyRoll(ma.userSlot.mon, target, ma.move) {
		log.Printf("%s's move %s missed", ma.userSlot.mon.base.Name, ma.move.Name)
		return
	}

	log.Printf("%s used %s", ma.userSlot.mon.base.Name, ma.move.Name)

	if ma.move.Name == "struggle" {
		ma.userSlot.mon.changeHpBy(-(ma.userSlot.mon.maxHP() / 4))
	}

	if ma.targetSlot.protected || ma.targetSlot.invulnerableAction != nil {
		log.Printf("but it failed")
		return
	}

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			ma.applyStatusMove(bs, ma.userSlot.mon, false)
		} else {
			ma.applyStatusMove(bs, target, true)
		}
	} else {
		ma.applyDamageMove(bs)
	}

	if ma.userSlot.mon.hp <= 0 {
		monFainted(bs, ma.userSlot, false)
		return
	}

	ma.userSlot.mon.checkItemTrigger(true, leppaBerryEvent{
		move: ma.move,
	})

	if ma.move.Type == "fire" {
		delete(ma.targetSlot.mon.ailments, "freeze")
	}
}

func (ma *moveAction) applyStatusMove(bs battleState, target *pokemon, offensive bool) {
	if _, ok := protectMoves[ma.move.Name]; ok {
		ma.userSlot.resolveProtect()
		return
	}

	if ma.move.Category == "field-effect" {
		ma.targetSlot.applyFieldEffect(ma.move.Name)
		return
	}

	switch ma.move.Name {
	case "swagger":
		target.changeStatStageBy("attack", 2, false)
		target.applyAilment("confusion", ma.move, ma.userSlot)
		return
	case "focus-energy":
		target.focusEnergy = true
		return
	case "laser-focus":
		target.laserFocus = true
		return
	case "belly-drum":
		if target.hp*2 <= target.maxHP() {
			log.Printf("but it failed")
			return
		}

		log.Printf("%s took damage from belly drum", target.base.Name)
		target.changeHpBy(-(target.maxHP() / 2))
		target.changeStatStageBy("attack", 6, false)
	}

	if ma.move.Heal > 0 {
		change := target.maxHP() * ma.move.Heal / 100
		target.changeHpBy(change)
		log.Printf("%s healed for %d", target.base.Name, change)
	}

	if ma.move.Ailment != "none" {
		target.applyAilment(ma.move.Ailment, ma.move, ma.userSlot)
	}

	for stat, change := range ma.move.StatChanges {
		target.changeStatStageBy(stat, change, offensive)
	}
}

func (ma *moveAction) applyDamageMove(bs battleState) {
	target := ma.targetSlot.mon
	if target.fainted {
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

	crit := determineCrit(user, target, ma.move)

	damage := calculateDamage(user, target, ma.move, crit, false, false, ma.pursuit)
	if damage == 0 {
		log.Printf("it does not affect %s", target.base.Name)
		return false
	}

	target.checkItemTrigger(true, resistBerryEvent{
		typeName: ma.move.Type,
	})
	target.checkItemTrigger(true, focusSashEvent{
		damage: &damage,
	})
	if target.ability == "sturdy" && target.hp == target.maxHP() {
		damage = min(damage, target.hp-1)
	}
	user.checkItemTrigger(true, gemEvent{
		typeName: ma.move.Type,
	})

	damage = min(damage, target.hp)
	log.Printf("%s took %d damage", target.base.Name, int(damage))
	if *crit {
		log.Printf("it was a critical hit!")
	}
	target.changeHpBy(-damage)
	if target.hp <= 0 {
		monFainted(bs, ma.targetSlot, ma.pursuit)
	}

	if ma.move.Name == "bug-bite" && strings.HasSuffix(target.item.name, "berry") && !target.item.consumed {
		log.Printf("%s's %s was consumed by bug bite", target.base.Name, target.item.name)
		item, _ := registerItem(target.item.name, user)
		item.activate()
		target.item, _ = registerItem("", target)

	} else if ma.move.Name == "wake-up-slap" {
		if a := target.hasAilment("sleep"); a != nil {
			log.Printf("%s woke up", target.base.Name)
			delete(target.ailments, "sleep")
		}
	} else if ma.move.Name == "knock-off" && !target.item.consumed && target.ability != "sticky-hold" {
		log.Printf("%s had its %s knocked off", target.base.Name, target.item.name)
		target.item = &item{
			consumed: true,
		}
	}

	if ma.move.Drain != 0 {
		change := damage * ma.move.Drain / 100
		if change == 0 {
			if ma.move.Drain > 0 {
				change = 1
			} else {
				change = -1
			}
		}
		if target.ability == "liquid-ooze" {
			if change > 0 {
				change = -change
			}
		}

		user.changeHpBy(change)
		if change >= 0 {
			log.Printf("%s healed for %d", user.base.Name, change)
		} else {
			log.Printf("%s took recoil for %d", user.base.Name, -change)
			if user.hp <= 0 {
				monFainted(bs, ma.userSlot, false)
			}
		}
	}

	if f, ok := contactDefensiveAbilities[target.ability]; ok && ma.move.Contact {
		f(ma.userSlot, ma.targetSlot)
		if user.hp <= 0 {
			monFainted(bs, ma.userSlot, false)
		}
	} else if target.ability == "cotten-down" {
		for _, slot := range bs.getOtherSlots(ma.targetSlot) {
			slot.mon.changeStatStageBy("speed", -1, true)
		}
	} else if target.ability == "water-compaction" && ma.move.Type == "water" {
		target.changeStatStageBy("defense", 2, false)
	}
	if f, ok := contactOffensiveAbilities[user.ability]; ok && ma.move.Contact {
		f(ma.userSlot, ma.targetSlot)
	}

	sg := user.serenceGraceBonus()

	if ma.move.StatChance > 0 && ma.move.Category == "damage-raise" && roll(ma.move.StatChance*sg, 100) {
		for stat, change := range ma.move.StatChanges {
			user.changeStatStageBy(stat, change, false)
		}
	}

	if _, ok := pivotMoves[ma.move.Name]; ok {
		if ma.userSlot.trainer.canReplace(bs) {
			injectReplaceAction(bs, ma.userSlot, true)
		}
	}

	if target.fainted {
		return false
	}

	if target.ability == "shield-dust" {
		return true
	}

	if ma.move.AilmentChance > 0 && !target.fainted && roll(ma.move.AilmentChance*sg, 100) {
		target.applyAilment(ma.move.Ailment, ma.move, ma.userSlot)
	}

	if ma.move.FlinchChance > 0 && !target.fainted && target.ability != "inner-focus" && roll(ma.move.FlinchChance*sg, 100) {
		if targetMove := bs.getActions().getMoveActionBy(target); targetMove != nil {
			targetMove.flinch = true
		}
	}

	if ma.move.StatChance > 0 && ma.move.Category == "damage-lower" && roll(ma.move.StatChance*sg, 100) {
		for stat, change := range ma.move.StatChanges {
			target.changeStatStageBy(stat, change, true)
		}
	}

	return true
}
