package main

import (
	"math/rand"
	"slices"
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

type PriorityQueue[T any] []T

func (q PriorityQueue[T]) len() int {
	return len(q)
}

func (q *PriorityQueue[T]) push(a T) {
	*q = append(*q, a)
}

func (q *PriorityQueue[T]) pop() (T, bool) {
	l := q.len()
	if l == 0 {
		var zero T
		return zero, false
	}

	a := (*q)[l-1]
	*q = (*q)[:l-1]

	return a, true
}

func (q *PriorityQueue[T]) insertAt(a T, cmp func(T, T) bool) {
	for i := 0; i < q.len(); i++ {
		if cmp(a, (*q)[i]) {
			*q = slices.Insert(*q, i, a)
			return
		}
	}

	q.push(a)
}

func (q *PriorityQueue[T]) sortBy(f func(a, b T) int) bool {
	if f == nil {
		return false
	}

	slices.SortFunc(*q, f)

	return true
}

func (q *PriorityQueue[T]) fetchBy(f func(T) bool) (T, bool) {
	for i, e := range *q {
		if f(e) {
			*q = append((*q)[:i], (*q)[i+1:]...)
			return e, true
		}
	}

	var zero T
	return zero, false
}

type ActionQueue []action

func (aq ActionQueue) Len() int {
	return len(aq)
}

func (aq ActionQueue) Less(i, j int) bool {
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

func (aq *ActionQueue) getSwitchActionBy(mon *Pokemon) *switchAction {
	for _, a := range *aq {
		if sa, ok := a.(*switchAction); ok && mon == sa.oldSlot.mon {
			return sa
		}
	}
	return nil
}
