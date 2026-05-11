package main

import (
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

func calculateDamage(user *pokemon.Pokemon, target *pokemon.Pokemon, move *pokeapi.BaseMove, crit bool, maxRoll bool) int {
	numerator := 1
	denominator := 1

	applyType := func(mult float64) {
		switch mult {
		case 0:
			numerator = 0
			denominator = 1
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

	stab := user.HasType(move.Type)

	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.EffectiveStat("attack", crit)
		defensiveStat = target.EffectiveStat("defense", crit)
	} else {
		offensiveStat = user.EffectiveStat("special-attack", crit)
		defensiveStat = target.EffectiveStat("special-defense", crit)
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

	if move.Class == "physical" && user.HasAilment("burn") {
		denominator *= 2
	}

	if !maxRoll {
		randFactor := rand.Intn(16) + 85
		numerator *= randFactor
		denominator *= 100
	}

	damage = max(1, damage*numerator/denominator)

	return damage
}

func roll(numerator int, denominator int) bool {
	return rand.Intn(denominator) < numerator
}

func accuracyRoll(user *pokemon.Pokemon, target *pokemon.Pokemon, moveAccuracy int) bool {
	accNum, accDen := user.AccuracyFraction()
	evNum, evDen := target.EvasionFraction()
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
