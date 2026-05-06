package main

import (
	"math/rand"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

var critRateMap = map[int]float32{
	0: 16.0,
	1: 8.0,
	2: 2.0,
	3: 1.0,
}

var confusionMove = pokeapi.BaseMove{
	Name:  "confusion",
	Type:  "no-type",
	Power: 40,
	Class: "physical",
}

func calculateDamage(user *pokemon.Pokemon, target *pokemon.Pokemon, move *pokeapi.BaseMove, crit bool, max bool) int {
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

	numerator := 1
	denominator := 1

	if stab {
		numerator *= 3
		denominator *= 2
	}

	if crit {
		numerator *= 3
		denominator *= 2
	}

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

	if !max {
		randFactor := rand.Intn(16) + 85
		numerator *= randFactor
		denominator *= 100
	}

	damage = damage * numerator / denominator
	if damage < 1 {
		damage = 1
	}

	return damage
}

func roll(chance float32) bool {
	return rand.Float32() < chance
}
