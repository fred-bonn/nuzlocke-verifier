package main

import (
	"slices"
)

type action interface {
	invoke(bs battleState)
	prio() int
	speed(bs battleState) int
}

func rollInt(numerator int, denominator int) int {
	if roll(numerator, denominator) {
		return 1
	}
	return 0
}

type PriorityQueue[T any] []T

func (q *PriorityQueue[T]) push(a T) {
	*q = append(*q, a)
}

func (q *PriorityQueue[T]) pop() (T, bool) {
	l := len(*q)
	if l == 0 {
		var zero T
		return zero, false
	}

	a := (*q)[l-1]
	*q = (*q)[:l-1]

	return a, true
}

func (q *PriorityQueue[T]) insertAt(a T, cmp func(T, T) bool) {
	for i := 0; i < len(*q); i++ {
		if cmp(a, (*q)[i]) {
			*q = slices.Insert(*q, i, a)
			return
		}
	}

	q.push(a)
}

func (q PriorityQueue[T]) sortBy(f func(a, b T) int) bool {
	if f == nil {
		return false
	}

	slices.SortFunc(q, f)

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

type ActionQueue struct {
	queue PriorityQueue[action]
}

func (a *ActionQueue) sort(bs battleState) {
	cmp := func(b, c action) int {
		if am, ok := b.(*moveAction); ok {
			if am.move.Name == "pursuit" {
				if _, ok := c.(*switchAction); ok {
					am.pursuit = true
					return 1
				}
			}
		}

		if b.prio() > c.prio() {
			return 1
		}
		if b.prio() < c.prio() {
			return -1
		}
		if b.speed(bs) > c.speed(bs) {
			return 1
		}
		if b.speed(bs) < c.speed(bs) {
			return -1
		}
		return rollInt(1, 2)*2 - 1
	}

	a.queue.sortBy(cmp)
}

func (a *ActionQueue) containstSwitchTo(mon *Pokemon) bool {
	for _, action := range a.queue {
		if sa, ok := action.(*switchAction); ok && sa.new == mon {
			return true
		}
	}
	return false
}

func (a *ActionQueue) getMoveActionBy(mon *Pokemon) *moveAction {
	for _, action := range a.queue {
		if ma, ok := action.(*moveAction); ok && mon == ma.userSlot.mon {
			return ma
		}
	}
	return nil
}
