package main

import "testing"

type benchAction struct {
	priority int
	spd      int
}

func (ba *benchAction) invoke(bs battleState)    {}
func (ba *benchAction) prio(bs battleState) int  { return ba.priority }
func (ba *benchAction) speed(bs battleState) int { return ba.spd }

type benchBattleState struct {
	actions *actionQueue
}

func (bs *benchBattleState) execute()                      {}
func (bs *benchBattleState) gatherActions()                {}
func (bs *benchBattleState) getAllSlots() []*slot          { return nil }
func (bs *benchBattleState) getOtherSlots(s *slot) []*slot { return nil }
func (bs *benchBattleState) getOpponentSlot(s *slot) *slot { return nil }
func (bs *benchBattleState) getActions() *actionQueue      { return bs.actions }
func (bs *benchBattleState) getWeather() weatherState      { return None }
func (bs *benchBattleState) setWeather(weatherState, int)  {}

func newEmptyActionQueue() *actionQueue {
	return &actionQueue{
		queue: make(priorityQueue[action], 0, 5),
	}
}

// BenchmarkActionQueueInit measures the cost of initializing the queue for a turn:
// creating a new queue, pushing actions, and then sorting for execution.
func BenchmarkActionQueueInit(b *testing.B) {
	actions := []action{
		&benchAction{priority: 1, spd: 50},
		&benchAction{priority: 5, spd: 30},
		&benchAction{priority: 3, spd: 70},
		&benchAction{priority: 2, spd: 40},
		&benchAction{priority: 4, spd: 60},
	}
	bs := &benchBattleState{}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q := newEmptyActionQueue()
		bs.actions = q
		for _, a := range actions {
			q.queue.push(a)
		}
		q.sort(bs)
	}
}

// BenchmarkActionQueueDrain measures the cost of draining the queue after sorting.
// Reuses the same queue allocation across iterations to focus on the drain cost.
func BenchmarkActionQueueDrain(b *testing.B) {
	actions := []action{
		&benchAction{priority: 2, spd: 50},
		&benchAction{priority: 5, spd: 30},
		&benchAction{priority: 3, spd: 70},
		&benchAction{priority: 2, spd: 40},
		&benchAction{priority: 4, spd: 60},
	}
	bs := &benchBattleState{}
	q := newEmptyActionQueue()
	bs.actions = q

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Push actions to the queue
		for _, a := range actions {
			q.queue.push(a)
		}
		q.sort(bs)
		// Drain the queue
		for len(q.queue) > 0 {
			q.queue.pop()
		}
	}
}
