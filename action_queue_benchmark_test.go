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
		q := newActionQueue(actions...)
		heap.Init(q)
	}
}

func BenchmarkActionQueueHeapPushPop(b *testing.B) {
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
		q := newActionQueue(actions...)
		heap.Init(q)
		heap.Push(q, &benchAction{priority: 6, spd: 55})
		_ = heap.Pop(q).(action)
	}
}

func BenchmarkActionQueueHeapDrain(b *testing.B) {
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
		q := newActionQueue(actions...)
		heap.Init(q)
		for q.Len() > 0 {
			_ = heap.Pop(q).(action)
		}
	}
}
