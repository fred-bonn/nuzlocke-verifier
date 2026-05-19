package main

import (
	"log"
	"math/rand"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type rnbAi struct{}

func (rnb rnbAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	scores := make([]int, len(actions))
	damage := make([]int, len(actions))
	var target *pokemon.Pokemon
	var user *pokemon.Pokemon
	fastDeadTo := make(map[string]bool)
	kills := false
	highestDamageIndex := -1
	canHighestKill := false

	for i, action := range actions {
		if action.move.Class == "status" {
			scores[i], _ = action.scoreActionMove(bs)
			continue
		}

		target = action.targetSlot.mon
		user = action.userSlot.mon
		isFastDead := false
		if deadTo, ok := fastDeadTo[target.Base.Name]; ok {
			isFastDead = deadTo && action.move.Priority > 0
		} else {
			fastDeadTo[target.Base.Name] = false
			if !user.IsFasterThan(target) {
				for _, move := range target.Moves {
					if move.PP > 0 && move.Class != "status" {
						dmg := calculateDamage(target, user, &move, move.CritRate >= 4, true)
						if user.Hp <= dmg {
							isFastDead = true
							fastDeadTo[target.Base.Name] = true
							break
						}
					}
				}
			}
		}

		damage[i], kills = action.scoreActionMove(bs)

		if damage[i] == 0 {
			scores[i] = -64
			continue
		}

		if action.move.Name == "fake-out" {
			if action.userSlot.firstTurn {
				scores[i] = 9
			} else {
				scores[i] = -64
				continue
			}
		} else if action.move.Name == "first-impression" && !action.userSlot.firstTurn {
			scores[i] = -64
			continue
		} else if action.move.Priority > 0 && isFastDead && target.IsFasterThan(user) {
			scores[i] = 9
		}

		if kills {
			canHighestKill = true
			highestDamageIndex = i
			if action.move.Priority > 0 || user.IsFasterThan(target) {
				scores[i] += 12 + 2*rollInt(1, 5)
			} else {
				scores[i] += 9 + 2*rollInt(1, 5)
			}
			continue
		} else if _, ok := target.Ailments["trap"]; !ok && action.move.Ailment == "trap" {
			scores[i] = 6 + 2*rollInt(1, 5)
		}

		if highestDamageIndex == -1 || (!canHighestKill && damage[i] > damage[highestDamageIndex]) {
			highestDamageIndex = i
		}
	}

	if highestDamageIndex != -1 && !canHighestKill {
		scores[highestDamageIndex] += 6 + 2*rollInt(1, 5)
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

	resultIndex := rand.Intn(len(bestIndices))

	return actions[bestIndices[resultIndex]], scores[bestIndices[resultIndex]]
}

func (rnb rnbAi) evaluteSwitchIns(bs battleState, mons []*pokemon.Pokemon, opponentSlot *slot) *pokemon.Pokemon {
	scores := make([]int, len(mons))
	opponent := opponentSlot.mon

	for i, mon := range mons {
		if mon.Base.Name == "ditto" || mon.Base.Name == "wobbufet" {
			scores[i] = 2
			continue
		}

		outspeeds := mon.IsFasterThan(opponent)

		monDamage := calculateMaxDamage(mon, opponent)
		opponentDamage := calculateMaxDamage(opponent, mon)

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

func calculateMaxDamage(user, target *pokemon.Pokemon) int {
	var maxDmg int
	for _, move := range user.Moves {
		if move.PP > 0 && move.Class != "status" {
			dmg := 0
			rolls := 1
			if move.MaxHits == 5 {
				rolls = 3
			} else if move.MaxHits > 0 {
				rolls = move.MaxHits
			}
			for i := 0; i < rolls; i++ {
				dmg += calculateDamage(user, target, &move, move.CritRate >= 4, false)
			}

			if dmg > maxDmg {
				maxDmg = dmg
			}
		}
	}
	return maxDmg
}
