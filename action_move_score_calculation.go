package main

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

	return damageRoll, damageRoll >= ma.targetSlot.mon.Hp
}

func (ma *moveAction) scoreStatusMove(bs battleState) int {
	return 6
}
