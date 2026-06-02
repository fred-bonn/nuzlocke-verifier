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

func calculateDamage(user, target *Pokemon, move *pokeapi.BaseMove, crit bool, maxRoll bool) int {
	numerator := 1
	denominator := 1

	if move.Name == "acrobatics" && (user.Item.consumed || user.Item.name == "flying-gem") {
		numerator *= 2
	}

	if move.Name == "flail" {
		res := int(48 * (float64(user.Hp) / float64(user.Stats["hp"])))
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

	applyType := func(mult float64) {
		if target.isGrounded() && target.hasType("flying") && move.Type == "ground" {
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

	applyType(pokemon.GetEffectiveness(move.Type, target.Base.Types[0]))
	if len(target.Base.Types) > 1 {
		applyType(pokemon.GetEffectiveness(move.Type, target.Base.Types[1]))
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

	stab := user.hasType(move.Type)

	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.effectiveStat("attack", crit)
		defensiveStat = target.effectiveStat("defense", crit)
	} else {
		offensiveStat = user.effectiveStat("special-attack", crit)
		defensiveStat = target.effectiveStat("special-defense", crit)
	}

	damage := ((((2*user.Level)/5)+2)*move.Power*offensiveStat)/defensiveStat/50 + 2

	if stab {
		numerator *= 3
		denominator *= 2
	}

	if crit {
		numerator *= 3
		denominator *= 2
	}

	if move.Class == "physical" && user.hasAilment("burn") {
		denominator *= 2
	}

	target.Item.checkTrigger(false, resistBerryEvent{
		typeName:    move.Type,
		denominator: &denominator,
	})

	user.Item.checkTrigger(false, gemEvent{
		typeName:    move.Type,
		denominator: &denominator,
		numerator:   &numerator,
	})

	user.Item.checkTrigger(false, choiceItemEvent{
		move:        move,
		denominator: &denominator,
		numerator:   &numerator,
	})

	if !maxRoll {
		numerator *= rand.Intn(16) + 85
		denominator *= 100
	}

	damage = max(1, damage*numerator/denominator)

	return damage
}

func roll(numerator int, denominator int) bool {
	return rand.Intn(denominator) < numerator
}

func accuracyRoll(user *Pokemon, target *Pokemon, moveAccuracy int) bool {
	accNum, accDen := user.accuracyFraction()
	evNum, evDen := target.evasionFraction()
	numerator := moveAccuracy * accNum * evNum
	denominator := 100 * accDen * evDen
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

func monFainted(bs battleState, slot *slot) {
	slot.mon.Fainted = true
	bs.injectReplaceAction(slot, bs.getTrainer(slot), false)
	log.Printf("%s fainted!", slot.mon.Base.Name)
}
