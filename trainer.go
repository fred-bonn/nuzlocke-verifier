package main

import (
	"log"
	"math/rand"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type trainer struct {
	pokemonParty []*pokemon.Pokemon
	player       bool
	ai           ai
	lost         bool
}

func (t *trainer) nextAction(bs battleState, slot *slot) action {
	opponentSlot := bs.getOtherSlots(slot)[0] // only works for single battles for now
	willSwitch := roll(1, 3) && slot.mon.Hp > slot.mon.Stats["hp"]/2 && bs.getTrainer(slot).canReplace(bs)
	if willSwitch && !slot.isTrapped() {
		var possibleMons []*pokemon.Pokemon
		for _, mon := range t.pokemonParty {
			if mon == slot.mon || mon.Fainted || bs.getActions().containstSwitchTo(mon) {
				continue
			}
			possibleMons = append(possibleMons, mon)
		}
		return &switchAction{
			oldSlot: slot,
			new:     t.ai.evaluteSwitchIns(bs, possibleMons, opponentSlot),
		}
	}

	possibleActions := make([]*moveAction, 0)
	for _, move := range slot.mon.Moves {
		if move.PP <= 0 {
			continue
		}
		possibleActions = append(possibleActions, &moveAction{
			userSlot:   slot,
			targetSlot: opponentSlot,
			move:       &move,
		})
	}

	return t.ai.evaluateActions(bs, possibleActions)
}

func (t *trainer) selectSwitchIn(bs battleState, slot *slot) *pokemon.Pokemon {
	var possibleMons []*pokemon.Pokemon
	for _, mon := range t.pokemonParty {
		if mon == slot.mon || mon.Fainted {
			continue
		}
		possibleMons = append(possibleMons, mon)
	}

	if len(possibleMons) == 0 {
		t.lost = true
		return nil
	}

	return t.ai.evaluteSwitchIns(bs, possibleMons, bs.getOtherSlots(slot)[0]) // only works for single battles for now)
}

func (t *trainer) canReplace(bs battleState) bool {
	count := 0
	for _, mon := range t.pokemonParty {
		if !mon.Fainted {
			count++
		}
		if count > 1 {
			return true
		}
	}
	return false
}

type ai interface {
	evaluateActions(bs battleState, actions []*moveAction) *moveAction
	evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon
}

type randomAi struct{}

func (ra randomAi) evaluateActions(bs battleState, actions []*moveAction) *moveAction {
	return actions[rand.Intn(len(actions))]
}

func (ra randomAi) evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon {
	return mons[rand.Intn(len(mons))]
}

type rnbAi struct{}

func (rnb rnbAi) evaluateActions(bs battleState, actions []*moveAction) *moveAction {
	scores := make([]int, len(actions))
	damage := make([]int, len(actions))
	kill := make([]bool, len(actions))
	highestDamageIndex := -1
	canHighestKill := false

	for i, action := range actions {
		if action.move.Class == "status" {
			scores[i], _ = action.score(bs)
			continue
		}

		damage[i], kill[i] = action.score(bs)
		if kill[i] {
			canHighestKill = true
			highestDamageIndex = i
			if action.move.Priority > 0 || action.userSlot.mon.EffectiveStat("speed", false) >= action.targetSlot.mon.EffectiveStat("speed", false) {
				scores[i] = 12 + 2*rollInt(1, 5)
			} else {
				scores[i] = 9 + 2*rollInt(1, 5)
			}
			continue
		}

		if highestDamageIndex == -1 || (!canHighestKill && damage[i] > damage[highestDamageIndex]) {
			highestDamageIndex = i
		}
	}

	if highestDamageIndex != -1 && !canHighestKill {
		scores[highestDamageIndex] = 6 + 2*rollInt(1, 5)
	}

	log.Println(scores)

	maxScore := scores[0]
	for _, score := range scores {
		if score > maxScore {
			maxScore = score
		}
	}

	var bestIndices []int
	for i, score := range scores {
		if score == maxScore {
			bestIndices = append(bestIndices, i)
		}
	}

	return actions[bestIndices[rand.Intn(len(bestIndices))]]
}

func (rnb rnbAi) evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon {
	scores := make([]int, len(mons))
	opponent := opponentSlot.mon

	for i, mon := range mons {
		if mon.Base.Name == "ditto" || mon.Base.Name == "wobbufet" {
			scores[i] = 2
			continue
		}

		monSpeed := mon.EffectiveStat("speed", false)
		opponentSpeed := opponent.EffectiveStat("speed", false)
		outspeeds := monSpeed >= opponentSpeed

		var monDamage, opponentDamage int
		for _, move := range mon.Moves {
			if move.PP > 0 && move.Class != "status" {
				dmg := calculateDamage(mon, opponent, &move, move.CritRate >= 4, true)
				if dmg > monDamage {
					monDamage = dmg
				}
			}
		}
		for _, move := range opponent.Moves {
			if move.PP > 0 && move.Class != "status" {
				dmg := calculateDamage(opponent, mon, &move, move.CritRate >= 4, true)
				if dmg > opponentDamage {
					opponentDamage = dmg
				}
			}
		}

		killsOpponent := monDamage >= opponent.Hp
		monKilled := opponentDamage >= mon.Hp

		monDamagePercent := monDamage * 100 / opponent.Hp
		opponentDamagePercent := opponentDamage * 100 / mon.Hp

		if outspeeds && killsOpponent {
			scores[i] = 5
		} else if !outspeeds && killsOpponent && !monKilled {
			scores[i] = 4
		} else if outspeeds && monDamagePercent > opponentDamagePercent {
			scores[i] = 3
		} else if !outspeeds && monDamagePercent > opponentDamagePercent {
			scores[i] = 2
		} else if outspeeds {
			scores[i] = 1
		} else if !outspeeds && monKilled {
			scores[i] = -1
		}
	}

	maxScore := scores[0]
	bestIndex := 0
	for i := 1; i < len(scores); i++ {
		if scores[i] > maxScore {
			maxScore = scores[i]
			bestIndex = i
		}
	}

	log.Println(scores[bestIndex])

	return mons[bestIndex]
}
