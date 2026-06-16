package main

import (
	"log"
	"math/rand"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
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

var confusionMove = pokeapi.BaseMove{
	Name:  "confusion",
	Type:  "no-type",
	Power: 40,
	Class: "physical",
}

var struggleMove = pokeapi.BaseMove{
	Name:  "struggle",
	Type:  "no-type",
	Power: 50,
	Class: "physical",
}

func calculateDamage(user, target *Pokemon, move *pokeapi.BaseMove, crit *bool, maxRoll, forScoring, pursuit bool) int {
	if f, ok := typeImmunityAbilities[target.Ability]; ok && f(target, move.Type, forScoring) && user.Ability != "mold-breaker" {
		return 0
	}

	numerator := 1
	denominator := 1
	moveType := move.Type
	power := move.Power
	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.effectiveStat("attack", *crit)
		defensiveStat = target.effectiveStat("defense", *crit)
	} else {
		offensiveStat = user.effectiveStat("special-attack", *crit)
		defensiveStat = target.effectiveStat("special-defense", *crit)
	}

	if f, ok := typeConvertingAbilities[user.Ability]; ok {
		f(&moveType, &power)
	}

	applyType := func(mult float64) {
		if target.isGrounded() && target.hasType("flying") && moveType == "ground" {
			return
		}
		switch mult {
		case 0:
			numerator = 0
		case 0.5:
			denominator *= 2
		case 1:
		case 2:
			numerator *= 2
		}
	}

	applyType(pokemon.GetEffectiveness(moveType, target.Base.Types[0]))
	if len(target.Base.Types) > 1 {
		applyType(pokemon.GetEffectiveness(moveType, target.Base.Types[1]))
	}
	if numerator == 0 {
		return 0
	}

	if move.Name == "psywave" {
		if maxRoll {
			return user.Level
		}
		return (user.Level * (rand.Intn(100) + 51)) / 100
	}
	if move.Name == "seismic-toss" || move.Name == "night-shade" {
		return user.Level
	}
	if move.Name == "sonic-boom" {
		return 20
	}
	if move.Name == "dragon-rage" {
		return 40
	} else if move.Name == "acrobatics" && (user.Item.consumed || user.Item.name == "flying-gem") {
		power *= 2
	} else if move.Name == "wake-up-slap" && target.hasAilment("sleep") != nil {
		power *= 2
	} else if move.Name == "venoshock" && (target.hasAilment("poison") != nil || target.hasAilment("toxic") != nil) {
		power *= 2
	} else if move.Name == "hex" && target.hasNonVolatileAilment() {
		power *= 2
	} else if move.Name == "flail" {
		res := int(48 * (float64(user.Hp) / float64(user.maxHP())))
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
	}

	if user.Ability == "technician" && move.Power <= 60 {
		power = power * 3 / 2
	} else if t, ok := pinchAbilities[user.Ability]; ok && t == moveType && user.Hp*3 <= user.maxHP() {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.FlashFire && moveType == "fire" {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.Ability == "hustle" && move.Class == "physical" {
		offensiveStat = offensiveStat * 3 / 2
	} else if user.Ability == "merciless" {
		if a := target.hasAilment("poison"); a != nil {
			*crit = true
		} else if a := target.hasAilment("toxic"); a != nil {
			*crit = true
		}
	}

	if _, ok := critBlockingAbilities[target.Ability]; ok && (forScoring || user.Ability != "mold-breaker") {
		*crit = false
	}

	if user.hasType(moveType) {
		numerator *= 3
		denominator *= 2
	}

	if target.Ability == "dry-skin" && moveType == "fire" {
		power = power * 5 / 4
	}

	if *crit {
		if user.Ability == "sniper" {
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

	damage := ((((2*user.Level)/5)+2)*power*offensiveStat)/defensiveStat/50 + 2
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

func accuracyRoll(user *Pokemon, target *Pokemon, move *pokeapi.BaseMove) bool {
	if user.Ability == "no-guard" || target.Ability == "no-guard" {
		return true
	} else if move.Name == "toxic" && user.hasType("poison") {
		return true
	} else if move.Name == "thunder-wave" && user.hasType("electric") {
		return true
	}

	moveAccuracy := move.Accuracy
	if user.Ability == "hustle" && move.Class == "physical" {
		moveAccuracy = moveAccuracy * 80 / 100
	}

	accNum, accDen := user.accuracyFraction()
	evNum, evDen := target.evasionFraction(user.Ability == "keen-eye")
	numerator := moveAccuracy * accNum * evNum
	denominator := 100 * accDen * evDen
	if user.Ability == "compound-eyes" {
		numerator *= 13
		denominator *= 10
	}

	return roll(numerator, denominator)
}

func determineHits(move *pokeapi.BaseMove) int {
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

func determineCrit(user, target *Pokemon, move *pokeapi.BaseMove) *bool {
	rate := determineCritRate(user, move)

	return new(roll(1, critRateMap[rate]))
}

func determineCritRate(user *Pokemon, move *pokeapi.BaseMove) int {
	if user.LaserFocus {
		return 3
	}

	rate := move.CritRate
	if user.Item.name == "scope-lens" {
		rate++
	}
	if user.Ability == "super-luck" {
		rate++
	}
	if user.FocusEnergy {
		rate += 2
	}

	return rate
}

func monFainted(bs battleState, slot *slot, pursuit bool) {
	if slot.mon.Fainted {
		return
	}

	slot.mon.Fainted = true
	if !pursuit {
		injectReplaceAction(bs, slot, false)
	}
	log.Printf("%s fainted!", slot.mon.Base.Name)
}
