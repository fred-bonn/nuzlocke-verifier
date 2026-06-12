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

func calculateDamage(user, target *Pokemon, move *pokeapi.BaseMove, crit *bool, maxRoll, forScoring bool) int {
	if f, ok := typeImmunityAbilities[target.Ability]; ok && f(target, move.Type, forScoring) {
		return 0
	}

	numerator := 1
	denominator := 1
	moveType := move.Type

	if f, ok := typeConvertingAbilities[user.Ability]; ok {
		f(&moveType, &numerator, &denominator)
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
	}

	if move.Name == "acrobatics" && (user.Item.consumed || user.Item.name == "flying-gem") {
		numerator *= 2
	} else if move.Name == "wake-up-slap" && target.hasAilment("sleep") != nil {
		numerator *= 2
	} else if move.Name == "venoshock" && (target.hasAilment("poison") != nil || target.hasAilment("toxic") != nil) {
		numerator *= 2
	} else if move.Name == "hex" && target.hasNonVolatileAilment() {
		numerator *= 2
	}

	if move.Name == "flail" {
		res := int(48 * (float64(user.Hp) / float64(user.maxHP())))
		if res <= 1 {
			move.Power = 200
		} else if res <= 4 {
			move.Power = 150
		} else if res <= 9 {
			move.Power = 100
		} else if res <= 16 {
			move.Power = 80
		} else if res <= 32 {
			move.Power = 40
		} else {
			move.Power = 20
		}
	}

	if user.Ability == "technician" && move.Power <= 60 {
		numerator *= 3
		denominator *= 2
	} else if t, ok := pinchAbilities[user.Ability]; ok && t == moveType && user.Hp*3 <= user.maxHP() {
		numerator *= 3
		denominator *= 2
	} else if user.FlashFire && moveType == "fire" {
		numerator *= 3
		denominator *= 2
	} else if user.Ability == "hustle" && move.Class == "physical" {
		numerator *= 3
		denominator *= 2
	} else if user.Ability == "merciless" {
		if a := target.hasAilment("poison"); a != nil {
			*crit = true
		} else if a := target.hasAilment("toxic"); a != nil {
			*crit = true
		}
	}

	if user.hasType(moveType) {
		numerator *= 3
		denominator *= 2
	}

	if target.Ability == "dry-skin" && moveType == "fire" {
		numerator *= 5
		denominator *= 4
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

	target.checkItemTrigger(false, resistBerryEvent{
		typeName:    moveType,
		denominator: &denominator,
	})

	user.checkItemTrigger(false, gemEvent{
		typeName:    moveType,
		denominator: &denominator,
		numerator:   &numerator,
	})

	user.checkItemTrigger(false, choiceItemEvent{
		move:        move,
		denominator: &denominator,
		numerator:   &numerator,
	})

	if !maxRoll {
		numerator *= rand.Intn(16) + 85
		denominator *= 100
	}

	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.effectiveStat("attack", *crit)
		defensiveStat = target.effectiveStat("defense", *crit)
	} else {
		offensiveStat = user.effectiveStat("special-attack", *crit)
		defensiveStat = target.effectiveStat("special-defense", *crit)
	}

	damage := ((((2*user.Level)/5)+2)*move.Power*offensiveStat)/defensiveStat/50 + 2

	damage = max(1, damage*numerator/denominator)

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
	if _, ok := critBlockingAbilities[target.Ability]; ok && user.Ability != "mold-breaker" {
		return new(false)
	} else if user.LaserFocus {
		return new(true)
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

	return new(roll(1, critRateMap[rate]))
}

func monFainted(bs battleState, slot *slot) {
	if slot.mon.Fainted {
		return
	}

	slot.mon.Fainted = true
	bs.injectReplaceAction(slot, false)
	log.Printf("%s fainted!", slot.mon.Base.Name)
}
