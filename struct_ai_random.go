package main

import "math/rand"

type ai interface {
	evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int)
	evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon
}

type randomAi struct{}

type learningAi struct {
	policy map[string]int
}

func newLearningAi() *learningAi {
	return &learningAi{
		policy: map[string]int{
			"damage":  8,
			"status":  5,
			"healing": 6,
			"switch":  5,
		},
	}
}

func (la *learningAi) observeOutcome(kind string, delta int) {
	if la == nil {
		return
	}
	if la.policy == nil {
		la.policy = make(map[string]int)
	}
	if _, ok := la.policy[kind]; !ok {
		la.policy[kind] = 0
	}
	la.policy[kind] += delta
	if la.policy[kind] > 100 {
		la.policy[kind] = 100
	}
	if la.policy[kind] < -100 {
		la.policy[kind] = -100
	}
}

func (ra randomAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	return actions[rand.Intn(len(actions))], 1
}

func (ra randomAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	return mons[rand.Intn(len(mons))]
}

func (la *learningAi) evaluateActions(bs battleState, actions []*moveAction) (*moveAction, int) {
	if len(actions) == 0 {
		return nil, 0
	}
	if la == nil || la.policy == nil {
		return rnbAi{}.evaluateActions(bs, actions)
	}

	bestAction := actions[0]
	bestScore := la.scoreAction(bs, bestAction)
	for _, action := range actions[1:] {
		score := la.scoreAction(bs, action)
		if score > bestScore {
			bestAction = action
			bestScore = score
		} else if score == bestScore && rand.Intn(2) == 0 {
			bestAction = action
		}
	}
	return bestAction, bestScore
}

func (la *learningAi) evaluteSwitchIns(bs battleState, mons []*pokemon, opponentSlot *slot) *pokemon {
	if len(mons) == 0 {
		return nil
	}
	if len(mons) == 1 {
		return mons[0]
	}
	if la == nil || la.policy == nil {
		return rnbAi{}.evaluteSwitchIns(bs, mons, opponentSlot)
	}

	bestMon := mons[0]
	bestScore := la.scoreSwitchIn(bs, bestMon, opponentSlot)
	for _, mon := range mons[1:] {
		score := la.scoreSwitchIn(bs, mon, opponentSlot)
		if score > bestScore {
			bestMon = mon
			bestScore = score
		} else if score == bestScore && rand.Intn(2) == 0 {
			bestMon = mon
		}
	}
	return bestMon
}

func (la *learningAi) scoreAction(bs battleState, action *moveAction) int {
	baseScore, _ := action.scoreActionMove(bs)
	category := "damage"
	if action.move.Class == statusClass {
		if action.move.Category == "heal" {
			category = "healing"
		} else {
			category = "status"
		}
	}
	if action.targetSlot != nil && action.targetSlot.mon != nil && action.targetSlot.mon.hp <= 0 {
		return -100000
	}
	return baseScore + la.policy[category]
}

func (la *learningAi) scoreSwitchIn(bs battleState, mon *pokemon, opponentSlot *slot) int {
	if mon == nil || opponentSlot == nil || opponentSlot.mon == nil {
		return -100000
	}
	baseDamage := calculateMaxDamage(bs, mon, opponentSlot.mon, false)
	opponentDamage := calculateMaxDamage(bs, opponentSlot.mon, mon, false)
	score := la.policy["switch"]
	if mon.isFasterThan(bs, opponentSlot.mon) {
		score += 4
	}
	if baseDamage >= opponentSlot.mon.hp {
		score += 12
	}
	if opponentDamage >= mon.hp {
		score -= 8
	}
	return score + (baseDamage-opponentDamage)/10
}
