package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

func (ma *moveAction) scoreActionMove(bs battleState) (int, bool) {
	if ma.move.Class == "status" {
		return ma.scoreStatusMove(bs), false
	}

	damageRoll := 0
	critRate := determineCritRate(ma.userSlot.mon, ma.move)
	rolls := 1
	if ma.move.MaxHits == 5 {
		rolls = 3
	} else if ma.move.MaxHits > 0 {
		rolls = ma.move.MaxHits
	}
	for i := 0; i < rolls; i++ {
		damageRoll += calculateDamage(ma.userSlot.mon, ma.targetSlot.mon, ma.move, new(critRate >= 3), false, true, false)
	}

	ma.targetSlot.mon.checkItemTrigger(false, focusSashEvent{
		damage: &damageRoll,
	})
	if ma.targetSlot.mon.Ability == "sturdy" && ma.targetSlot.mon.Hp == ma.targetSlot.mon.maxHP() {
		damageRoll = min(damageRoll, ma.targetSlot.mon.Hp-1)
	}

	return damageRoll, damageRoll >= ma.targetSlot.mon.Hp
}

func (ma *moveAction) scoreStatusMove(bs battleState) int {
	if ma.move.Category == "heal" {
		if ma.userSlot.mon.Hp > ma.userSlot.mon.maxHP()*85/100 {
			return -64
		}
		if ma.shouldMonHeal(bs) {
			return 7
		}
		return 5
	}

	if _, ok := powderMoves[ma.move.Name]; ok && (ma.targetSlot.mon.hasType("grass") || ma.targetSlot.mon.Ability == "overcoat") {
		return -64
	}
	if _, ok := paralysisMoves[ma.move.Name]; ok {
		return ma.scoreParalysisMove(bs)
	}
	if _, ok := sleepMoves[ma.move.Name]; ok {
		return ma.scoreSleepMove(bs)
	}
	if _, ok := protectMoves[ma.move.Name]; ok {
		return ma.scoreProtectMove(bs)
	}

	switch ma.move.Name {
	case "sticky-web":
		if ma.targetSlot.hasFieldEffect(ma.move.Name) {
			return -64
		}
		if ma.userSlot.firstTurn {
			return 9 + 3*rollInt(3, 4)
		}
		return 6 + 3*rollInt(3, 4)
	case "stealth-rock", "spikes", "toxic-spikes":
		if ma.targetSlot.hasFieldEffect(ma.move.Name) {
			return -64
		}
		if ma.userSlot.firstTurn {
			return 8 + rollInt(3, 4)
		}
		return 6 + rollInt(3, 4)
	case "attract":
		if ma.targetSlot.mon.hasAilment("infatuation") != nil || ma.targetSlot.mon.Ability == "oblivious" {
			return -64
		}
	case "leech-seed":
		if ma.targetSlot.mon.hasAilment("leech-seed") != nil || ma.targetSlot.mon.hasType("grass") {
			return -64
		}
	case "toxic":
		return ma.scoreToxic()
	case "focus-energy", "laser-focus":
		return ma.scoreCritStatus()
	case "belly-drum":
		return ma.scoreBellyDrum()
	}

	return 6
}

func (ma *moveAction) shouldMonHeal(bs battleState) bool {
	if ma.userSlot.mon.hasAilment("toxic") != nil {
		return false
	}

	maxDmg := calculateMaxDamage(ma.targetSlot.mon, ma.userSlot.mon, true)
	if maxDmg >= ma.userSlot.mon.maxHP()*ma.move.Heal/100 {
		return false
	}

	if ma.userSlot.mon.isFasterThan(bs, ma.targetSlot.mon) {
		if maxDmg < min(ma.userSlot.mon.maxHP(), ma.userSlot.mon.Hp+ma.userSlot.mon.maxHP()*ma.move.Heal/100) {
			return true
		} else {
			if ma.userSlot.mon.Hp < ma.userSlot.mon.maxHP()*40/100 {
				return true
			} else if ma.userSlot.mon.Hp <= ma.userSlot.mon.maxHP()*66/100 {
				return roll(1, 2)
			}
		}
	} else {
		if ma.userSlot.mon.Hp < ma.userSlot.mon.maxHP()*50/100 {
			return true
		} else if ma.userSlot.mon.Hp <= ma.userSlot.mon.maxHP()*70/100 {
			return roll(3, 4)
		}
	}

	return false
}

func (ma *moveAction) scoreParalysisMove(bs battleState) int {
	target := ma.targetSlot.mon
	user := ma.userSlot.mon

	if target.hasNonVolatileAilment() || target.hasType("electric") || target.Ability == "limber" {
		return -64
	}

	score := 6
	if target.isFasterThan(bs, user) && user.effectiveSpeed(bs) > target.effectiveSpeed(bs)/4 {
		score++
	} else if user.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
		return m.Name == "hex" || m.FlinchChance > 0
	}) {
		score++
	} else if target.hasAilment("confusion") != nil {
		score++
	} else if target.hasAilment("infatuation") != nil {
		score++
	}

	return score + rollInt(1, 2)
}

func (ma *moveAction) scoreSleepMove(bs battleState) int {
	target := ma.targetSlot.mon
	user := ma.userSlot.mon

	if _, ok := sleepBlockingAbilities[target.Ability]; ok {
		return -64
	}
	if a := target.hasAilment("yawn"); a != nil {
		return -64
	}
	if target.hasNonVolatileAilment() {
		return -64
	}

	isHex := func(m *pokeapi.BaseMove) bool {
		return m.Name == "hex"
	}

	score := 6
	maxDmg := calculateMaxDamage(target, ma.userSlot.mon, true)
	if maxDmg < user.Hp && roll(1, 2) {
		if user.hasMovePredicate(isHex) {
			score += 1
		} else {
			for _, slot := range bs.getOtherSlots(ma.userSlot) {
				if slot.trainer == ma.userSlot.trainer && slot.mon.hasMovePredicate(isHex) {
					score += 1
				}
			}
		}

		if user.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
			return m.Name == "dream-eater" || m.Name == "nightmare"
		}) && !target.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
			return m.Name == "snore" || m.Name == "sleep-talk"
		}) {
			score += 1
		}
	}

	return score
}

func (ma *moveAction) scoreToxic() int {
	target := ma.targetSlot.mon
	user := ma.userSlot.mon

	if target.hasNonVolatileAilment() {
		return -64
	}
	if target.Ability == "immunity" {
		return -64
	}
	if (target.hasType("poison") || target.hasType("steel")) && ma.userSlot.mon.Ability != "corrosion" {
		return -64
	}

	score := 6
	maxDmg := calculateMaxDamage(target, user, true)
	if maxDmg < user.Hp && roll(19, 50) {
		if !target.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
			return m.Class == "physical" || m.Class == "special"
		}) {
			score += 1
		}

		if user.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
			return m.Name == "hex" || m.Name == "venoshock"
		}) || user.Ability == "merciless" {
			score += 2
		} else {
			score += 1
		}
	}

	return score
}

func (ma *moveAction) scoreProtectMove(bs battleState) int {
	user := ma.userSlot.mon
	target := ma.targetSlot.mon

	if ma.userSlot.protectTurns == 2 || (ma.userSlot.protectTurns == 1 && roll(1, 2)) {
		return -64
	}
	if deadToSecondaryDamage(user, bs) {
		return -64
	}

	score := 6
	if ma.userSlot.firstTurn {
		score--
	}
	// still needs perish song and cursed
	if target.hasAilment("poison") != nil || target.hasAilment("toxic") != nil || target.hasAilment("burn") != nil || target.hasAilment("leech-seed") != nil || target.hasAilment("yawn") != nil || target.hasAilment("infatuation") != nil {
		score++
	}

	if user.hasAilment("poison") != nil || user.hasAilment("toxic") != nil || user.hasAilment("burn") != nil || user.hasAilment("leech-seed") != nil || user.hasAilment("yawn") != nil || user.hasAilment("infatuation") != nil {
		score -= 2
	}

	return score
}

func deadToSecondaryDamage(mon *Pokemon, bs battleState) bool {
	// update when weather is implemented
	dmg := 0
	if a := mon.hasAilment("burn"); a != nil {
		dmg += mon.maxHP() / 16
	} else if a := mon.hasAilment("poison"); a != nil {
		dmg += mon.maxHP() / 8
	} else if a := mon.hasAilment("toxic"); a != nil {
		dmg += (mon.maxHP() * (a.Turns + 1)) / 16
	}
	if a := mon.hasAilment("trap"); a != nil {
		dmg += mon.maxHP() / 8
	}

	return dmg >= mon.Hp
}

func (ma *moveAction) scoreCritStatus() int {
	user := ma.userSlot.mon
	if _, ok := critBlockingAbilities[ma.targetSlot.mon.Ability]; ok && user.Ability != "mold-breaker" {
		return -64
	}
	if ma.move.Name == "focus-energy" && user.FocusEnergy {
		return -64
	}

	if user.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
		return m.CritRate > 0
	}) {
		return 7
	} else if user.Item.name == "scope-lens" {
		return 7
	} else if user.Ability == "super-luck" || user.Ability == "sniper" {
		return 7
	}

	return 6
}

func (ma *moveAction) scoreBellyDrum() int {
	user := ma.userSlot.mon
	target := ma.targetSlot.mon

	if a := target.hasAilment("freeze"); a != nil && target.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
		if _, ok := selfThawingMoves[m.Name]; ok {
			return true
		}
		return false
	}) {
		return 9
	}
	if a := target.hasAilment("sleep"); a != nil {
		return 9
	}

	dmg := calculateMaxDamage(target, user, true)
	threshhold := user.maxHP() / 2
	if user.Item.name == "sitrus-berry" {
		threshhold += user.maxHP() / 4
	}
	if dmg < threshhold {
		return 8
	}

	return 4
}
