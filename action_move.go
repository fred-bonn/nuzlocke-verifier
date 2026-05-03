package main

import (
	"log"
	"math/rand"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type moveAction struct {
	userSlot   *slot
	targetSlot *slot
	move       pokeapi.BaseMove
}

func (ma *moveAction) prio() int {
	return ma.move.Priority
}

func (ma *moveAction) speed() int {
	return ma.userSlot.mon.EffectiveStat("speed")
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.Fainted {
		return
	}

	target := bs.getMon(ma.targetSlot)
	hitChance := target.EffectiveEvasion() * ma.userSlot.mon.EffectiveAccuracy()
	if ma.move.Accuracy != 0 && hitChance < 1.0 {
		if !roll(hitChance) {
			log.Printf("%s's move missed", ma.userSlot.mon.Base.Name)
			return
		}
	}

	log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			ma.applyStatusMove(bs, ma.userSlot.mon, ma.move)
		} else {
			ma.applyStatusMove(bs, target, ma.move)
		}
	} else {
		ma.applyDamageMove(bs)
	}
}

func roll(chance float32) bool {
	return rand.Float32() < chance
}

func (ma *moveAction) applyStatusMove(bs battleState, target *pokemon.Pokemon, move pokeapi.BaseMove) {
	if move.StatChance == 100 || roll(float32(move.StatChance/100)) {
		for stat, change := range move.StatChanges {
			target.Stages[stat] = max(-6, min(6, target.Stages[stat]+change))
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

func (ma *moveAction) applyDamageMove(bs battleState) {
	user := ma.userSlot.mon
	target := ma.targetSlot.mon
	move := ma.move

	crit := roll(1.0 / critRateMap[move.CritRate])
	stab := user.HasType(move.Type)

	var offensiveStat, defensiveStat int
	if move.Class == "physical" {
		offensiveStat = user.EffectiveStat("attack")
		defensiveStat = target.EffectiveStat("defense")
	} else {
		offensiveStat = user.EffectiveStat("special-attack")
		defensiveStat = target.EffectiveStat("special-defense")
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

	randFactor := rand.Intn(16) + 85
	numerator *= randFactor
	denominator *= 100

	damage = damage * numerator / denominator
	if damage < 1 {
		damage = 1
	}

	log.Printf("%s took %d damage", target.Base.Name, int(damage))
	if crit {
		log.Printf("it was a critical hit!")
	}

	target.Hp -= int(damage)
	if target.Hp <= 0 {
		target.Hp = 0
		target.Fainted = true
		bs.injectReplaceAction(ma.targetSlot, bs.getTrainer(ma.targetSlot), false)
		log.Printf("%s fainted!", target.Base.Name)
	}

	if _, ok := pivotMoves[move.Name]; ok {
		trainer := bs.getTrainer(ma.userSlot)
		if trainer.canReplace(bs) {
			bs.injectReplaceAction(ma.userSlot, trainer, true)
		}
	}
}
