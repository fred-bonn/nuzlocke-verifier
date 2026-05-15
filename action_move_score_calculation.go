package main

func (ma *moveAction) score(bs battleState) (int, bool) {
	if ma.move.Class == "status" {
		return 6, false
	}

	damageRoll := calculateDamage(ma.userSlot.mon, ma.targetSlot.mon, ma.move, ma.move.CritRate >= 4, false)

	return damageRoll, damageRoll >= ma.targetSlot.mon.Hp
}
