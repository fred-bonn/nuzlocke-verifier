package main

import (
	"math/rand"
)

var critRateMap = map[int]int{
	0: 16,
	1: 8,
	2: 2,
	3: 1,
	4: 1,
	5: 1,
	6: 1,
}

var confusionMove = move{
	Name:  "confusion",
	Type:  "no-type",
	Power: 40,
	Class: "physical",
}

var struggleMove = move{
	Name:  "struggle",
	Type:  "no-type",
	Power: 50,
	Class: "physical",
}

func calculateDamage(user, target *pokemon, move *move, crit *bool, weather weatherState, maxRoll, forScoring, pursuit bool) int {
	if f, ok := typeImmunityAbilities[target.ability]; ok && user.ability != "mold-breaker" && f(target, move.Type, forScoring) {
		return 0
	}

	numerator := 1
	denominator := 1
	moveType := move.Type
	power := move.Power
	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.effectiveStat(Attack, *crit)
		defensiveStat = target.effectiveStat(Defense, *crit)
	} else {
		offensiveStat = user.effectiveStat(SpecialAttack, *crit)
		defensiveStat = target.effectiveStat(SpecialDefense, *crit)
		if weather == Sandstorm && target.hasType("rock") {
			defensiveStat = defensiveStat * 3 / 2
		}
	}

	if f, ok := typeConvertingAbilities[user.ability]; ok {
		f(&moveType, &power)
	}
	numerator, denominator = target.applyMoveType(numerator, denominator, moveType)
	if weather != None {
		if f, ok := weatherFuncs[weather]; ok {
			f(&numerator, &denominator, moveType)
		}
		switch weather {
		case Sun:
			if user.ability == "solar-power" && move.Class == "special" {
				offensiveStat = offensiveStat * 3 / 2
			}
		case Sandstorm:
			if user.ability == "sand-force" && (moveType == "rock" || moveType == "ground" || moveType == "steel") {
				power = power * 13 / 10
			}
		}
	}
	if numerator == 0 {
		return 0
	}

	switch move.Name {
	case "psywave":
		if maxRoll {
			return user.level
		}
		*crit = false
		return (user.level * (rand.Intn(100) + 51)) / 100
	case "seismic-toss", "night-shade":
		*crit = false
		return user.level
	case "sonic-boom":
		*crit = false
		return 20
	case "dragon-rage":
		*crit = false
		return 40
	case "endeavor":
		*crit = false
		return target.hp - user.hp
	case "super-fang":
		*crit = false
		return max(1, target.hp/2)
	}

	if move.Name == "acrobatics" && (user.item.consumed || user.item.name == "flying-gem") {
		power *= 2
	} else if move.Name == "wake-up-slap" && target.hasAilment("sleep") != nil {
		power *= 2
	} else if move.Name == "venoshock" && (target.hasAilment("poison") != nil || target.hasAilment("toxic") != nil) {
		power *= 2
	} else if move.Name == "hex" && target.hasNonVolatileAilment() {
		power *= 2
	} else if move.Name == "flail" || move.Name == "reversal" {
		res := int(48 * (float64(user.hp) / float64(user.maxHP())))
		if res <= 1 {
			power = 200
		} else if res <= 4 {
			power = 150
		} else if res <= 9 {
			power = 100
		} else if res <= 16 {
			power = 80
		} else if res <= 32 {
			power = 40
		} else {
			power = 20
		}
	} else if move.Name == "pursuit" && pursuit {
		power *= 2
	} else if move.Name == "knock-off" && !target.item.consumed {
		power *= 2
	}

	if user.ability == "technician" && move.Power <= 60 {
		power = power * 3 / 2
	} else if t, ok := pinchAbilities[user.ability]; ok && t == moveType && user.hp*3 <= user.maxHP() {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.flashFire && moveType == "fire" {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.ability == "hustle" && move.Class == "physical" {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.ability == "merciless" {
		if a := target.hasAilment("poison"); a != nil {
			*crit = true
		} else if a := target.hasAilment("toxic"); a != nil {
			*crit = true
		}
	}

	if _, ok := critBlockingAbilities[target.ability]; ok && (forScoring || user.ability != "mold-breaker") {
		*crit = false
	}

	if user.hasType(moveType) {
		numerator *= 3
		denominator *= 2
	}

	if target.ability == "dry-skin" && moveType == "fire" {
		power = power * 5 / 4
	}

	if *crit {
		if user.ability == "sniper" {
			numerator *= 3
			denominator *= 2
		}
		numerator *= 3
		denominator *= 2
	}

	if move.Class == "physical" && user.hasAilment("burn") != nil {
		denominator *= 2
	}

	user.checkItemTrigger(false, gemEvent{
		typeName: moveType,
		power:    &power,
	})

	user.checkItemTrigger(false, choiceItemEvent{
		move: move,
		stat: &offensiveStat,
	})

	user.checkItemTrigger(false, moveBoostingEvent{
		power:    &power,
		typeName: moveType,
	})

	if !maxRoll {
		numerator *= rand.Intn(16) + 85
		denominator *= 100
	}

	damage := ((((2*user.level)/5)+2)*power*offensiveStat)/defensiveStat/50 + 2
	damage = damage * numerator / denominator

	target.checkItemTrigger(false, resistBerryEvent{
		typeName: moveType,
		damage:   &damage,
	})

	damage = max(1, damage)

	return damage
}

func roll(numerator int, denominator int) bool {
	return rand.Intn(denominator) < numerator
}

func rollInt(numerator int, denominator int) int {
	if roll(numerator, denominator) {
		return 1
	}
	return 0
}

func accuracyRoll(bs battleState, user *pokemon, target *pokemon, move *move) bool {
	if user.ability == "no-guard" || target.ability == "no-guard" {
		return true
	} else if move.Name == "toxic" && user.hasType("poison") {
		return true
	} else if move.Name == "thunder-wave" && user.hasType("electric") {
		return true
	}

	moveAccuracy := move.Accuracy
	if user.ability == "hustle" && move.Class == "physical" {
		moveAccuracy = moveAccuracy * 80 / 100
	}

	accNum, accDen := user.accuracyFraction()
	evNum, evDen := target.evasionFraction(user.ability == "keen-eye")
	numerator := moveAccuracy * accNum * evNum
	denominator := 100 * accDen * evDen
	if user.ability == "compound-eyes" {
		numerator *= 13
		denominator *= 10
	}

	if bs.getWeather() != None {
		switch bs.getWeather() {
		case Hail:
			if target.ability == "snow-cloak" {
				numerator *= 4
				denominator *= 5
			}
		case Sandstorm:
			if target.ability == "sand-veil" {
				numerator *= 4
				denominator *= 5
			}
		}
	}

	return roll(numerator, denominator)
}

func determineHits(move *move) int {
	if move.MaxHits == 5 && move.MinHits == 2 {
		r := rand.Intn(100) + 1
		if r <= 35 {
			return 2
		} else if r <= 70 {
			return 3
		} else if r <= 85 {
			return 4
		} else {
			return 5
		}
	}
	return move.MaxHits
}

func determineCrit(user, target *pokemon, move *move) *bool {
	rate := determineCritRate(user, move)

	return new(roll(1, critRateMap[rate]))
}

func determineCritRate(user *pokemon, move *move) int {
	if user.laserFocus {
		return 3
	}

	rate := move.CritRate
	if user.item.name == "scope-lens" {
		rate++
	}
	if user.ability == "super-luck" {
		rate++
	}
	if user.focusEnergy {
		rate += 2
	}

	return rate
}

func monFainted(bs battleState, slot *slot, pursuit bool) {
	if slot.mon.fainted {
		return
	}

	slot.mon.fainted = true
	if !pursuit {
		injectReplaceAction(bs, slot, false)
	}
	vlogf("%s fainted!", slot.mon.base.Name)
}

func fetchPursuitMiddleware(name string) func(a action) bool {
	return func(a action) bool {
		ma, ok := a.(*moveAction)
		if !ok {
			return false
		}
		if ma.move.Name != "pursuit" {
			return false
		}
		if ma.targetSlot.mon.base.Name != name {
			return false
		}
		return true
	}
}
