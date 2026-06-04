package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

func (ma *moveAction) scoreActionMove(bs battleState) (int, bool) {
	if ma.move.Class == "status" {
		return ma.scoreStatusMove(bs), false
	}

	damageRoll := 0
	rolls := 1
	if ma.move.MaxHits == 5 {
		rolls = 3
	} else if ma.move.MaxHits > 0 {
		rolls = ma.move.MaxHits
	}
	for i := 0; i < rolls; i++ {
		damageRoll += calculateDamage(ma.userSlot.mon, ma.targetSlot.mon, ma.move, ma.move.CritRate >= 4, false)
	}

	ma.targetSlot.mon.Item.checkTrigger(false, focusSashEvent{
		damage: &damageRoll,
	})

	return damageRoll, damageRoll >= ma.targetSlot.mon.Hp
}

func (ma *moveAction) scoreStatusMove(bs battleState) int {
	if ma.move.Category == "heal" {
		if ma.userSlot.mon.Hp > ma.userSlot.mon.Stats["hp"]*85/100 {
			return -64
		}
		if ma.shouldMonHeal(bs) {
			return 7
		}
		return 5
	}

	if _, ok := powderMoves[ma.move.Name]; ok && ma.targetSlot.mon.hasType("grass") {
		return -64
	}
	if _, ok := paralysisMoves[ma.move.Name]; ok {
		return ma.scoreParalysisMove()
	}
	if _, ok := protectMoves[ma.move.Name]; ok {
		return ma.scoreProtectMove(bs)
	}

	switch ma.move.Name {
	case "sticky-web":
		if ma.userSlot.firstTurn {
			return 9 + 3*rollInt(3, 4)
		}
		return 6 + 3*rollInt(3, 4)
	case "attract":
		if ma.targetSlot.mon.hasAilment("infatuation") != nil {
			return -64
		}
	case "leech-seed":
		if ma.targetSlot.mon.hasAilment("leech-seed") != nil || ma.targetSlot.mon.hasType("grass") {
			return -64
		}
	}

	return 6
}

func (ma *moveAction) shouldMonHeal(bs battleState) bool {
	if ma.userSlot.mon.hasAilment("toxic") != nil {
		return false
	}

	maxDmg := calculateMaxDamage(ma.targetSlot.mon, ma.userSlot.mon, true)
	if maxDmg >= ma.userSlot.mon.Stats["hp"]*ma.move.Heal/100 {
		return false
	}

	if ma.userSlot.mon.isFasterThan(ma.targetSlot.mon) {
		if maxDmg < min(ma.userSlot.mon.Stats["hp"], ma.userSlot.mon.Hp+ma.userSlot.mon.Stats["hp"]*ma.move.Heal/100) {
			return true
		} else {
			if ma.userSlot.mon.Hp < ma.userSlot.mon.Stats["hp"]*40/100 {
				return true
			} else if ma.userSlot.mon.Hp <= ma.userSlot.mon.Stats["hp"]*66/100 {
				return roll(1, 2)
			}
		}
	} else {
		if ma.userSlot.mon.Hp < ma.userSlot.mon.Stats["hp"]*50/100 {
			return true
		} else if ma.userSlot.mon.Hp <= ma.userSlot.mon.Stats["hp"]*70/100 {
			return roll(3, 4)
		}
	}

	return false
}

func (ma *moveAction) scoreParalysisMove() int {
	if ma.targetSlot.mon.hasNonVolatileAilment() || ma.targetSlot.mon.hasType("electric") {
		return -64
	}

	bonus := 0
	if ma.targetSlot.mon.isFasterThan(ma.userSlot.mon) && ma.userSlot.mon.effectiveSpeed() > ma.targetSlot.mon.effectiveSpeed()/4 {
		bonus++
	} else if ma.userSlot.mon.hasMovePredicate(func(m *pokeapi.BaseMove) bool {
		return m.Name == "hex" || m.FlinchChance > 0
	}) {
		bonus++
	} else if ma.targetSlot.mon.hasAilment("confusion") != nil {
		bonus++
	} else if ma.targetSlot.mon.hasAilment("infatuation") != nil {
		bonus++
	}

	return 6 + bonus + rollInt(1, 2)
}

func (ma *moveAction) scoreProtectMove(bs battleState) int {
	// still needs to return if user is dead to secondary damage, and minus score if other volatile status are active
	if ma.userSlot.protectTurns == 2 || (ma.userSlot.protectTurns == 1 && roll(1, 2)) {
		return -64
	}

	bonus := 0
	if ma.userSlot.firstTurn {
		bonus--
	}
	if ma.targetSlot.mon.hasAilment("poison") != nil || ma.targetSlot.mon.hasAilment("toxic") != nil || ma.targetSlot.mon.hasAilment("burn") != nil {
		bonus++
	}
	if ma.userSlot.mon.hasAilment("poison") != nil || ma.userSlot.mon.hasAilment("toxic") != nil || ma.userSlot.mon.hasAilment("burn") != nil {
		bonus -= 2
	}

	return 6 + bonus
}
