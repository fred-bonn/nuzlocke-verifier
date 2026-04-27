package main

import (
	"log"
	"math/rand"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type action interface {
	invoke(bs battleState)
	prio() int
	speed() int
}

type swapAction struct {
	old *pokemon.Pokemon
	new *pokemon.Pokemon
}

func (sa *swapAction) invoke(sbs battleState) {
	sbs.setMon(sa.old, sa.new)
	sa.old.ResetStages()
	log.Printf("swapped %s for %s", sa.old.Base.Name, sa.new.Base.Name)
}

func (sa *swapAction) prio() int {
	return 6
}

func (sa *swapAction) speed() int {
	return sa.old.EffectiveStat("speed")
}

type moveAction struct {
	mon  *pokemon.Pokemon
	slot int
	move pokeapi.BaseMove
}

func (ma *moveAction) prio() int {
	return ma.move.Priority
}

func (ma *moveAction) speed() int {
	return ma.mon.EffectiveStat("speed")
}

func (ma *moveAction) invoke(sbs battleState) {
	if ma.mon.Fainted {
		return
	}

	target := sbs.getMon(ma.slot)
	hitChance := target.EffectiveEvasion() * ma.mon.EffectiveAccuracy()
	if ma.move.Accuracy != 0 && hitChance < 1.0 {
		if !roll(hitChance) {
			log.Printf("%s's move missed", ma.mon.Base.Name)
			return
		}
	}

	log.Printf("%s used %s", ma.mon.Base.Name, ma.move.Name)

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			applyStatusMove(ma.mon, ma.move)
		} else {
			applyStatusMove(target, ma.move)
		}
	} else {
		applyDamageMove(target, ma.mon, ma.move)
	}
}

func roll(chance float32) bool {
	return rand.Float32() < chance
}

func applyStatusMove(target *pokemon.Pokemon, move pokeapi.BaseMove) {
	if move.StatChance == 100 || roll(float32(move.StatChance/100)) {
		for stat, change := range move.StatChanges {
			target.Stages[stat] += change
			log.Printf("%s's %s changed by %d stages (%d)", target.Base.Name, stat, change, target.Stages[stat])
		}
	}
}

var critRateMap = map[int]float32{
	0: 16.0,
	1: 8.0,
	2: 2.0,
	3: 1.0,
}

func applyDamageMove(target *pokemon.Pokemon, mon *pokemon.Pokemon, move pokeapi.BaseMove) {
	crit := roll(1.0 / critRateMap[move.CritRate])
	stab := mon.HasType(move.Type)

	var offensiveStat, defensiveStat int
	if move.Category == "physical" {
		offensiveStat = mon.EffectiveStat("attack")
		defensiveStat = target.EffectiveStat("defense")
	} else {
		offensiveStat = mon.EffectiveStat("special-attack")
		defensiveStat = target.EffectiveStat("special-defense")
	}

	base := ((2*mon.Level)/5 + 2)
	damage := (base * move.Power * offensiveStat) / defensiveStat
	damage = damage / 50
	damage += 2

	if stab {
		damage = int(float32(damage) * 1.5)
	}
	damage = int(float32(damage) * pokemon.GetEffectiveness(move.Type, target.Base.Types[0]))
	if len(target.Base.Types) > 1 {
		damage = int(float32(damage) * pokemon.GetEffectiveness(move.Type, target.Base.Types[1]))
	}
	if crit {
		damage = int(float32(damage) * 1.5)
	}
	randFactor := rand.Intn(16) + 85
	damage = damage * randFactor / 100

	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if crit {
		log.Printf("it was a critical hit!")
	}
	target.Hp -= int(damage)
	if target.Hp <= 0 {
		target.Hp = 0
		target.Fainted = true
		log.Printf("%s fainted!", target.Base.Name)
	}
}
