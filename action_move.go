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
	return ma.userSlot.mon.EffectiveStat("speed", false)
}

func (ma *moveAction) invoke(bs battleState) {
	if ma.userSlot.mon.Fainted {
		return
	}

	if turns, ok := ma.userSlot.mon.Ailments["confusion"]; ok {
		if turns > 0 {
			ma.userSlot.mon.Ailments["confusion"] -= 1
			log.Printf("%s is confused", ma.userSlot.mon.Base.Name)
			if roll(1.0 / 3.0) {
				log.Printf("%s hit itself in confusion", ma.userSlot.mon.Base.Name)
				// implement self damage
				return
			}
		} else {
			delete(ma.userSlot.mon.Ailments, "confusion")
			log.Printf("%s snapped out of confusion", ma.userSlot.mon.Base.Name)
		}
	}

	target := bs.getMon(ma.targetSlot)
	hitChance := target.EffectiveEvasion() * ma.userSlot.mon.EffectiveAccuracy()
	if ma.move.Accuracy != 0 {
		if !roll(float32(ma.move.Accuracy) / 100.0 * hitChance) {
			log.Printf("%s's move %s missed", ma.userSlot.mon.Base.Name, ma.move.Name)
			return
		}
	}

	log.Printf("%s used %s", ma.userSlot.mon.Base.Name, ma.move.Name)

	if ma.move.Class == "status" {
		if strings.HasPrefix(ma.move.Target, "user") {
			ma.applyStatusMove(bs, ma.userSlot.mon)
		} else {
			ma.applyStatusMove(bs, target)
		}
	} else {
		ma.applyDamageMove(bs)
	}
}

func roll(chance float32) bool {
	return rand.Float32() < chance
}

func (ma *moveAction) applyStatusMove(bs battleState, target *pokemon.Pokemon) {
	if ma.move.Category == "swagger" {
		ma.applySwagger(bs, target)
		return
	}

	if ma.move.StatChance != 100 && !roll(float32(ma.move.StatChance/100)) {
		log.Printf("missed")
		return
	}

	for stat, change := range ma.move.StatChanges {
		target.Stages[stat] = max(-6, min(6, target.Stages[stat]+change))
		log.Printf("%s's %s changed by %d stages (%d)", target.Base.Name, stat, change, target.Stages[stat])
	}
}

func (ma *moveAction) applySwagger(bs battleState, target *pokemon.Pokemon) {
	err := target.ChangeStatStage("attack", 2)
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = target.ApplyAilment("confusion")
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("%s's attack changed by 2 stages (%d)", target.Base.Name, target.Stages["attack"])
	log.Printf("%s became confused", target.Base.Name)
}

var critRateMap = map[int]float32{
	0: 16.0,
	1: 8.0,
	2: 2.0,
	3: 1.0,
}

func (ma *moveAction) applyDamageMove(bs battleState) {
	crit := roll(1.0 / critRateMap[ma.move.CritRate])
	damage := ma.calculateDamage(bs, crit, false)
	target := ma.targetSlot.mon

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

	if _, ok := pivotMoves[ma.move.Name]; ok {
		trainer := bs.getTrainer(ma.userSlot)
		if trainer.canReplace(bs) {
			bs.injectReplaceAction(ma.userSlot, trainer, true)
		}
	}
}

func (ma *moveAction) calculateDamage(bs battleState, crit bool, max bool) int {
	user := ma.userSlot.mon
	target := ma.targetSlot.mon
	move := ma.move
	stab := ma.userSlot.mon.HasType(ma.move.Type)

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
