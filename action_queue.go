package main

import (
	"math/rand"
)

type action interface {
	invoke(bs battleState)
	prio() int
	speed() int
}

func rollInt(numerator int, denominator int) int {
	if roll(numerator, denominator) {
		return 1
	}
	return 0
}

type ActionQueue []action

func (aq ActionQueue) Len() int {
	return len(aq)
}

func (aq ActionQueue) Less(i, j int) bool {
	if aq[i].prio() < aq[j].prio() {
		return true
	} else if aq[i].prio() > aq[j].prio() {
		return false
	} else if aq[i].speed() < aq[j].speed() {
		return true
	} else if aq[i].speed() > aq[j].speed() {
		return false
	}
	return (rand.Int() % 2) == 0
}

func (aq ActionQueue) Swap(i, j int) {
	aq[i], aq[j] = aq[j], aq[i]
}

func (aq *ActionQueue) Push(a any) {
	*aq = append(*aq, a.(action))
}

func (aq *ActionQueue) Pop() any {
	if aq.Len() == 0 {
		return nil
	}
	action := (*aq)[0]
	*aq = (*aq)[1:]
	return action
}

func (aq *ActionQueue) containstSwitchTo(mon *Pokemon) bool {
	for _, a := range *aq {
		if sa, ok := a.(*switchAction); ok && sa.new == mon {
			return true
		}
	}
	return false
}

func (aq *ActionQueue) getMoveActionBy(mon *Pokemon) *moveAction {
	for _, a := range *aq {
		if ma, ok := a.(*moveAction); ok && mon == ma.userSlot.mon {
			return ma
		}
	}
	return nil
}
