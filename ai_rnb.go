package main

import (
	"log"
	"math/rand"
)

type rnbAi struct{}

func (rnb rnbAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	scores := make([]int, len(actions))
	damage := make([]int, len(actions))
	kills := make([]bool, len(actions))
	highestDamageIndex := -1
	canHighestKill := false

	for i, a := range actions {
		if a.move.PP <= 0 {
			damage[i] = -1
			scores[i] = -64
			continue
		}
		if a.move.Class == "status" {
			damage[i] = -1
			scores[i], _ = a.scoreActionMove(bs)
			continue
		}
		if a.move.Name == "rollout" {
			damage[i] = -1
			scores[i] = 7
			continue
		}
		if a.move.Name == "fake-out" {
			if a.userSlot.firstTurn {
				scores[i] = 9
			} else {
				damage[i] = -1
				scores[i] = -64
				continue
			}
		}
		if a.move.Name == "first-impression" && !a.userSlot.firstTurn {
			damage[i] = -1
			scores[i] = -64
			continue
		}
		if a.move.Ailment == "trap" {
			damage[i] = -1
			if _, ok := a.targetSlot.mon.Ailments["trap"]; !ok {
				scores[i] = 6 + 2*rollInt(1, 5)
			}
			continue
		}
		if a.move.Name == "sucker-punch" && a.userSlot.suckerPunch && roll(1, 2) {
			scores[i] = -20
		}

		damage[i], kills[i] = a.scoreActionMove(bs)
		if damage[i] == 0 {
			damage[i] = -1
			scores[i] = -64
		}

		canHighestKill = canHighestKill || kills[i]
		if !canHighestKill && (highestDamageIndex == -1 || damage[highestDamageIndex] < damage[i]) {
			highestDamageIndex = i
		}
	}

	for i, a := range actions {
		if damage[i] == -1 {
			continue
		}

		// add additional scoring to damaging moves
		if a.move.Name == "fell-stinger" && kills[i] {
			if a.userSlot.mon.isFasterThan(a.targetSlot.mon) {
				scores[i] = 21 + 2*rollInt(1, 5)
			} else {
				scores[i] = 15 + 2*rollInt(1, 5)
			}
		} else if a.move.Name == "acid-spray" {
			scores[i] = 6
		} else if a.move.Name == "future-sight" {
			// needs 8 score instead if faster and dead to target
			scores[i] = 6
		} else if a.move.Name == "pursuit" {
			if kills[i] {
				scores[i] = 10
			} else {
				if a.targetSlot.mon.Hp < a.targetSlot.mon.Stats["hp"]*20/100 {
					scores[i] = 10
				} else if a.targetSlot.mon.Hp < a.targetSlot.mon.Stats["hp"]*40/100 {
					scores[i] = 8 * rollInt(1, 2)
				}
			}
			if a.userSlot.mon.isFasterThan(a.targetSlot.mon) {
				scores[i] += 3
			}
		}

		// add score if fast dead and the move has priority
		if a.move.Priority > 0 && !a.userSlot.mon.isFasterThan(a.targetSlot.mon) {
			for _, move := range a.targetSlot.mon.Moves {
				if move.PP <= 0 || move.Class == "status" {
					continue
				}
				if a.targetSlot.mon.LockedMove != nil && a.targetSlot.mon.LockedMove != move {
					continue
				}
				dmg := calculateDamage(a.targetSlot.mon, a.userSlot.mon, move, move.CritRate >= 4, true)
				if a.userSlot.mon.Hp <= dmg {
					scores[i] += 11
					break
				}
			}
		}

		// if the highest damaging move kills then we only have to consider moves that can kill
		if canHighestKill {
			if !kills[i] {
				continue
			}
			if a.move.Priority > 0 || a.userSlot.mon.isFasterThan(a.targetSlot.mon) {
				scores[i] += 12 + 2*rollInt(1, 5)
			} else {
				scores[i] += 9 + 2*rollInt(1, 5)
			}
			continue
		}

		// if no moves kills then the highest damaging moves gets additional score
		if highestDamageIndex == i {
			scores[i] += 6 + 2*rollInt(1, 5)
			continue
		}

		// moves from this point that gets a base score if and only if it neither kills or is highest damage
		if _, ok := speedControlMoves[a.move.Name]; ok && !a.userSlot.mon.isFasterThan(a.targetSlot.mon) {
			scores[i] = 6
			continue
		}

		if c, ok := offenseControlMoves[a.move.Name]; ok {
			if a.targetSlot.mon.hasMoveClass(c) {
				scores[i] = 6
			} else {
				scores[i] = 5
			}
			continue
		}
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

func (rnb rnbAi) evaluteSwitchIns(bs battleState, mons []*Pokemon, opponentSlot *slot) *Pokemon {
	if len(mons) == 1 {
		return mons[0]
	}

	scores := make([]int, len(mons))
	opponent := opponentSlot.mon

	for i, mon := range mons {
		if mon.Base.Name == "ditto" || mon.Base.Name == "wobbufet" {
			scores[i] = 2
			continue
		}

		outspeeds := mon.isFasterThan(opponent)

		monDamage := calculateMaxDamage(mon, opponent, false)
		opponentDamage := calculateMaxDamage(opponent, mon, false)

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
		log.Printf("%s switch in: outspeeds=%t, monDamage=%d, opponentDamage=%d, killsOpponent=%t, monKilled=%t, monDamagePercent=%d, opponentDamagePercent=%d, score=%d\n", mon.Base.Name, outspeeds, monDamage, opponentDamage, killsOpponent, monKilled, monDamagePercent, opponentDamagePercent, scores[i])
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

func calculateMaxDamage(user, target *Pokemon, checkChoice bool) int {
	var maxDmg, dmg int
	rolls := 1
	for _, move := range user.Moves {
		if move.PP <= 0 || move.Class == "status" {
			continue
		}
		if checkChoice && user.LockedMove != nil && user.LockedMove != move {
			continue
		}

		rolls = 1
		if move.MaxHits == 5 {
			rolls = 3
		} else if move.MaxHits > 0 {
			rolls = move.MaxHits
		}
		for i := 0; i < rolls; i++ {
			dmg += calculateDamage(user, target, move, move.CritRate >= 4, true)
		}

		target.Item.checkTrigger(false, focusSashEvent{
			damage: &dmg,
		})

		if dmg > maxDmg {
			maxDmg = dmg
		}
		dmg = 0
	}
	return maxDmg
}
