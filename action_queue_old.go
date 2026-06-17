package main

import (
	"math/rand"
)

type ActionQueueOld []action

func (aq ActionQueueOld) Len() int {
	return len(aq)
}

func (aq ActionQueueOld) Less(i, j int) bool {
	if am, ok := aq[i].(*moveAction); ok {
		if am.move.Name == "pursuit" {
			if _, ok := aq[j].(*switchAction); ok {
				am.pursuit = true
				return false
			}
		}
	}

	if aq[i].prio() < aq[j].prio() {
		return true
	} else if aq[i].prio() > aq[j].prio() {
		return false
	} else if aq[i].speed(nil) < aq[j].speed(nil) {
		return true
	} else if aq[i].speed(nil) > aq[j].speed(nil) {
		return false
	}

	return (rand.Int() % 2) == 0
}

func (aq ActionQueueOld) Swap(i, j int) {
	aq[i], aq[j] = aq[j], aq[i]
}

func (aq *ActionQueueOld) Push(a any) {
	*aq = append(*aq, a.(action))
}

func (aq *ActionQueueOld) Pop() any {
	if aq.Len() == 0 {
		return nil
	}
	action := (*aq)[0]
	*aq = (*aq)[1:]
	return action
}
