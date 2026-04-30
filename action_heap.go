package main

import "math/rand"

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
