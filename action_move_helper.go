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

var confusionMove = Move{
	Name:  "confusion",
	Type:  noType,
	Power: 40,
	Class: physicalClass,
}

var struggleMove = Move{
	Name:  "struggle",
	Type:  noType,
	Power: 50,
	Class: physicalClass,
}

func calculateDamage(user, target *pokemon, move *Move, crit *bool, weather weatherState, maxRoll, forScoring, pursuit bool) int {
	if f, ok := typeImmunityAbilities[target.ability]; ok && user.ability != moldBreakerAbility && f(target, move.Type, forScoring) {
		return 0
	}

	numerator := 1
	denominator := 1
	moveType := move.Type
	power := move.Power
	var offensiveStat, defensiveStat int
	if move.Class == physicalClass {
		offensiveStat = user.effectiveStat(attack, *crit)
		defensiveStat = target.effectiveStat(defense, *crit)
	} else {
		offensiveStat = user.effectiveStat(specialAttack, *crit)
		defensiveStat = target.effectiveStat(specialDefense, *crit)
		if weather == sandstormWeather && target.hasType(rockType) {
			defensiveStat = defensiveStat * 3 / 2
		}
	}

	if f, ok := typeConvertingAbilities[user.ability]; ok {
		f(&moveType, &power)
	}
	numerator, denominator = target.applyMoveType(numerator, denominator, moveType)
	if weather != noneWeather {
		if f, ok := weatherFuncs[weather]; ok {
			f(&numerator, &denominator, moveType)
		}
		switch weather {
		case sunWeather:
			if user.ability == solarPowerAbility && move.Class == specialClass {
				offensiveStat = offensiveStat * 3 / 2
			}
		case sandstormWeather:
			if user.ability == sandForceAbility && (moveType == rockType || moveType == groundType || moveType == steelType) {
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
	case "seismic toss", "night shade":
		*crit = false
		return user.level
	case "sonic boom":
		*crit = false
		return 20
	case "dragon rage":
		*crit = false
		return 40
	case "endeavor":
		*crit = false
		return target.hp - user.hp
	case "super fang":
		*crit = false
		return max(1, target.hp/2)
	}

	if move.Name == "acrobatics" && (user.item.consumed || user.item.state == flyingGem) {
		power *= 2
	} else if move.Name == "wake up slap" && target.hasAilment(sleepAilment) != nil {
		power *= 2
	} else if move.Name == "venoshock" && (target.hasAilment(poisonAilment) != nil || target.hasAilment(toxicAilment) != nil) {
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
	} else if move.Name == "knock off" && !target.item.consumed {
		power *= 2
	}

	if user.ability == technicianAbility && move.Power <= 60 {
		power = power * 3 / 2
	} else if t, ok := pinchAbilities[user.ability]; ok && t == moveType && user.hp*3 <= user.maxHP() {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.flashFire && moveType == fireType {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.ability == hustleAbility && move.Class == physicalClass {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.ability == mercilessAbility {
		if a := target.hasAilment(poisonAilment); a != nil {
			*crit = true
		} else if a := target.hasAilment(toxicAilment); a != nil {
			*crit = true
		}
	}

	if target.ability.blocksCrits() && (forScoring || user.ability != moldBreakerAbility) {
		*crit = false
	}

	if user.hasType(moveType) {
		numerator *= 3
		denominator *= 2
	}

	if target.ability == drySkinAbility && moveType == fireType {
		power = power * 5 / 4
	}

	if *crit {
		if user.ability == sniperAbility {
			numerator *= 3
			denominator *= 2
		}
		numerator *= 3
		denominator *= 2
	}

	if move.Class == physicalClass && user.hasAilment(burnAilment) != nil {
		denominator *= 2
	}

	user.checkItemTrigger(false, gemEvent{
		pokemonType: moveType,
		power:       &power,
	})

	user.checkItemTrigger(false, choiceItemEvent{
		move: move,
		stat: &offensiveStat,
	})

	user.checkItemTrigger(false, moveBoostingEvent{
		power:       &power,
		pokemonType: moveType,
	})

	if !maxRoll {
		numerator *= rand.Intn(16) + 85
		denominator *= 100
	}

	damage := ((((2*user.level)/5)+2)*power*offensiveStat)/defensiveStat/50 + 2
	damage = damage * numerator / denominator

	target.checkItemTrigger(false, resistBerryEvent{
		pokemonType: moveType,
		damage:      &damage,
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

func accuracyRoll(bs battleState, user *pokemon, target *pokemon, move *Move) bool {
	if user.ability == noGuardAbility || target.ability == noGuardAbility {
		return true
	} else if move.Name == "toxic" && user.hasType(poisonType) {
		return true
	} else if move.Name == "thunder wave" && user.hasType(electricType) {
		return true
	}

	moveAccuracy := move.Accuracy
	if user.ability == hustleAbility && move.Class == physicalClass {
		moveAccuracy = moveAccuracy * 80 / 100
	}

	accNum, accDen := user.accuracyFraction()
	evNum, evDen := target.evasionFraction(user.ability == keenEyeAbility)
	numerator := moveAccuracy * accNum * evNum
	denominator := 100 * accDen * evDen
	if user.ability == compoundEyesAbility {
		numerator *= 13
		denominator *= 10
	}

	if bs.getWeather() != noneWeather {
		switch bs.getWeather() {
		case hailWeather:
			if target.ability == snowCloakAbility {
				numerator *= 4
				denominator *= 5
			}
		case sandstormWeather:
			if target.ability == sandVeilAbility {
				numerator *= 4
				denominator *= 5
			}
		}
	}

	return roll(numerator, denominator)
}

func determineHits(move *Move) int {
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

func determineCrit(user, target *pokemon, move *Move) *bool {
	rate := determineCritRate(user, move)

	return new(roll(1, critRateMap[rate]))
}

func determineCritRate(user *pokemon, move *Move) int {
	if user.laserFocus {
		return 3
	}

	rate := move.CritRate
	if user.item.state == scopeLens {
		rate++
	}
	if user.ability == superLuckAbility {
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
	vprintf("%s fainted!", slot.mon.base.Name)
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
