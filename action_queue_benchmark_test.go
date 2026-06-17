package main

import (
	"container/heap"
	"testing"
)

type benchAction struct {
	priority int
	spd      int
}

func (ba *benchAction) invoke(bs battleState) {}
func (ba *benchAction) prio() int             { return ba.priority }
func (ba *benchAction) speed() int            { return ba.spd }

func newActionQueue(actions ...action) *ActionQueue {
	q := &ActionQueue{}
	*q = append(*q, actions...)
	return q
}

func newEmptyActionQueue() *ActionQueue {
	res := make(ActionQueue, 0, 5)
	return &res
}

// BenchmarkActionQueueHeapInit measures the cost of initializing the heap for a turn:
// creating a new heap, initializing it, and pushing actions for execution.
func BenchmarkActionQueueHeapInit(b *testing.B) {
	actions := []action{
		&benchAction{priority: 1, spd: 50},
		&benchAction{priority: 5, spd: 30},
		&benchAction{priority: 3, spd: 70},
		&benchAction{priority: 2, spd: 40},
		&benchAction{priority: 4, spd: 60},
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := newEmptyActionQueue()
		heap.Init(q)
		for _, a := range actions {
			heap.Push(q, a)
		}
	}
}

// BenchmarkActionQueueHeapDrain measures the cost of draining the heap after initialization.
// Reuses the same heap allocation across iterations to focus on the drain cost.
func BenchmarkActionQueueHeapDrain(b *testing.B) {
	actions := []action{
		&benchAction{priority: 1, spd: 50},
		&benchAction{priority: 5, spd: 30},
		&benchAction{priority: 3, spd: 70},
		&benchAction{priority: 2, spd: 40},
		&benchAction{priority: 4, spd: 60},
	}

	q := newEmptyActionQueue()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Push actions to the heap
		for _, a := range actions {
			heap.Push(q, a)
		}
		// Drain the heap
		for q.Len() > 0 {
			_ = heap.Pop(q).(action)
		}
	}
}
