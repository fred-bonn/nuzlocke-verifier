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
	if ma.userSlot.mon.Ability == "prankster" && ma.move.Class == "status" {
		bonus++
	}

	return ma.move.Priority + bonus
}

func (ma *moveAction) speed(bs battleState) int {
	return ma.userSlot.mon.effectiveSpeed(bs)
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.Fainted {
		return
	}

	if (ma.move.Name == "fake-out" || ma.move.Name == "first-impression") && !ma.userSlot.firstTurn {
		log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)
		log.Printf("but it failed")
		return
	}

	ma.userSlot.firstTurn = false

	ma.userSlot.suckerPunch = ma.move.Name == "sucker-punch"

	if _, ok := ma.userSlot.mon.Ailments["freeze"]; ok {
		if _, ok := selfThawingMoves[ma.move.Name]; ok || roll(1, 5) {
			log.Printf("%s thawed out", ma.userSlot.mon.Base.Name)
			delete(ma.userSlot.mon.Ailments, "freeze")
		} else {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is frozen", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if sleep := ma.userSlot.mon.hasAilment("sleep"); sleep != nil {
		if sleep.Turns <= 0 {
			log.Printf("%s woke up", ma.userSlot.mon.Base.Name)
			delete(ma.userSlot.mon.Ailments, "sleep")
		} else {
			if ma.userSlot.mon.Ability == "early-bird" {
				sleep.Turns -= 2
			} else {
				sleep.Turns--
			}
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is asleep", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if confusion := ma.userSlot.mon.hasAilment("confusion"); confusion != nil {
		if confusion.Turns > 0 {
			confusion.Turns -= 1
			log.Printf("%s is confused", ma.userSlot.mon.Base.Name)
			if roll(1, 3) {
				damage := calculateDamage(ma.userSlot.mon, ma.userSlot.mon, &confusionMove, new(false), false, false, false)
				ma.userSlot.invulnerableAction = nil
				log.Printf("%s hit itself in confusion for %d damage", ma.userSlot.mon.Base.Name, damage)
				ma.userSlot.mon.Hp -= int(damage)
				if ma.userSlot.mon.Hp <= 0 {
					monFainted(bs, ma.userSlot, false)
				}
				return
			}
		} else {
			delete(ma.userSlot.mon.Ailments, "confusion")
			log.Printf("%s snapped out of confusion", ma.userSlot.mon.Base.Name)
		}
	}

	if paralysis := ma.userSlot.mon.hasAilment("paralysis"); paralysis != nil {
		if roll(1, 4) {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is paralysed", ma.userSlot.mon.Base.Name)
			return
		}
	}

	if infatuation := ma.userSlot.mon.hasAilment("infatuation"); infatuation != nil {
		if roll(1, 2) {
			ma.userSlot.invulnerableAction = nil
			log.Printf("%s is infatuated with %s", ma.userSlot.mon.Base.Name, infatuation.AfflictedBy.mon.Base.Name)
			return
		}
	}

	if ma.flinch {
		log.Printf("%s flinched", ma.userSlot.mon.Base.Name)
		return
	}

	ma.move.PP--

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
	if ma.move.Accuracy > 0 && !ma.pursuit && !accuracyRoll(ma.userSlot.mon, target, ma.move) {
		log.Printf("%s's move %s missed", ma.userSlot.mon.Base.Name, ma.move.Name)
		return
	}

	log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)

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

	if ma.userSlot.mon.Hp <= 0 {
		monFainted(bs, ma.userSlot, false)
		return
	}

	ma.userSlot.mon.checkItemTrigger(true, leppaBerryEvent{
		move: ma.move,
	})

	if ma.move.Type == "fire" {
		delete(ma.targetSlot.mon.Ailments, "freeze")
	}
}

func (ma *moveAction) applyStatusMove(bs battleState, target *Pokemon, offensive bool) {
	if _, ok := protectMoves[ma.move.Name]; ok {
		ma.userSlot.resolveProtect()
		return
	}

	switch ma.move.Name {
	case "swagger":
		target.changeStatStageBy("attack", 2, false)
		target.applyAilment("confusion", ma.move, ma.userSlot)
		return
	case "focus-energy":
		target.FocusEnergy = true
		return
	case "laser-focus":
		target.LaserFocus = true
		return
	case "belly-drum":
		if target.Hp*2 <= target.maxHP() {
			log.Printf("but it failed")
			return
		}
		target.changeHpBy(-(target.maxHP() / 2))
		target.changeStatStageBy("attack", 6, false)
	}

	if ma.move.Heal > 0 {
		change := target.maxHP() * ma.move.Heal / 100
		target.changeHpBy(change)
		log.Printf("%s healed for %d", target.Base.Name, change)
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

	crit := determineCrit(user, target, ma.move)

	damage := calculateDamage(user, target, ma.move, crit, false, false, ma.pursuit)
	if damage == 0 {
		log.Printf("it does not affect %s", target.Base.Name)
		return false
	}

	target.checkItemTrigger(true, resistBerryEvent{
		typeName: ma.move.Type,
	})
	target.checkItemTrigger(true, focusSashEvent{
		damage: &damage,
	})
	if target.Ability == "sturdy" && target.Hp == target.maxHP() {
		damage = min(damage, target.Hp-1)
	}
	user.checkItemTrigger(true, gemEvent{
		typeName: ma.move.Type,
	})

	damage = min(damage, target.Hp)
	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if *crit {
		log.Printf("it was a critical hit!")
	}
	if ma.move.Name == "bug-bite" && strings.HasSuffix(target.Item.name, "berry") && !target.Item.consumed {
		log.Printf("%s's %s was consumed by bug bite", target.Base.Name, target.Item.name)
		item, _ := registerItem(target.Item.name, user)
		item.activate()
		target.Item, _ = registerItem("", target)

	}
	target.changeHpBy(-damage)
	if target.Hp <= 0 {
		monFainted(bs, ma.targetSlot, ma.pursuit)
	}

	if ma.move.Name == "wake-up-slap" {
		if a := target.hasAilment("sleep"); a != nil {
			log.Printf("%s woke up", target.Base.Name)
			delete(target.Ailments, "sleep")
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
		if target.Ability == "liquid-ooze" {
			if change > 0 {
				change = -change
			}
		}

		user.changeHpBy(change)
		if change >= 0 {
			log.Printf("%s healed for %d", user.Base.Name, change)
		} else {
			log.Printf("%s took recoil for %d", user.Base.Name, -change)
			if user.Hp <= 0 {
				monFainted(bs, ma.userSlot, false)
			}
		}
	}

	if f, ok := contactDefensiveAbilities[target.Ability]; ok && ma.move.Contact {
		f(ma.userSlot, ma.targetSlot)
		if user.Hp <= 0 {
			monFainted(bs, ma.userSlot, false)
		}
	} else if target.Ability == "cotten-down" {
		for _, slot := range bs.getOtherSlots(ma.targetSlot) {
			slot.mon.changeStatStageBy("speed", -1, true)
		}
	}
	if f, ok := contactOffensiveAbilities[user.Ability]; ok && ma.move.Contact {
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

	if target.Fainted {
		return false
	}

	if target.Ability == "shield-dust" {
		return true
	}

	if ma.move.AilmentChance > 0 && !target.Fainted && roll(ma.move.AilmentChance*sg, 100) {
		target.applyAilment(ma.move.Ailment, ma.move, ma.userSlot)
	}

	if ma.move.FlinchChance > 0 && !target.Fainted && target.Ability != "inner-focus" && roll(ma.move.FlinchChance*sg, 100) {
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
