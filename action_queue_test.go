package main

import "testing"

func TestPriorityQueuePush(t *testing.T) {
	tests := map[string]struct {
		queue priorityQueue[int]
		push  []int
		want  priorityQueue[int]
	}{
		"empty":          {queue: priorityQueue[int]{}, push: []int{1, 2, 3}, want: priorityQueue[int]{1, 2, 3}},
		"three elements": {queue: priorityQueue[int]{1, 2, 3}, push: []int{1, 2, 3}, want: priorityQueue[int]{1, 2, 3, 1, 2, 3}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, element := range tc.push {
				tc.queue.push(element)
			}
			if n := len(tc.queue); n != len(tc.want) {
				t.Fatalf("%s: tc.queue.len() = %d, want %d", name, n, len(tc.want))
			}
			for i, r := range tc.want {
				if r != tc.queue[i] {
					t.Errorf("%s: mismatch at index %d, %d != %d", name, i, r, tc.queue[i])
				}
			}
		})
	}
}

func TestPriorityQueuePop(t *testing.T) {
	tests := map[string]struct {
		queue    priorityQueue[int]
		want     int
		wantOk   bool
		wantLeft priorityQueue[int]
	}{
		"empty":          {queue: priorityQueue[int]{}, want: 0, wantOk: false, wantLeft: priorityQueue[int]{}},
		"three elements": {queue: priorityQueue[int]{1, 2, 3}, want: 3, wantOk: true, wantLeft: priorityQueue[int]{1, 2}},
		"one element":    {queue: priorityQueue[int]{1}, want: 1, wantOk: true, wantLeft: priorityQueue[int]{}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, ok := tc.queue.pop()
			if tc.wantOk && !ok {
				t.Fatalf("%s: queue is empty when it wasn't expected to be", name)
			} else if !tc.wantOk && ok {
				t.Fatalf("%s: queue is not empty when it was expected to be", name)
			}
			if ok && res != tc.want {
				t.Fatalf("%s: got %d, want %d", name, res, tc.want)
			}
			if len(tc.queue) != len(tc.wantLeft) {
				t.Fatalf("%s: tc.queue.len() = %d, want %d", name, len(tc.queue), len(tc.wantLeft))
			}
			for i, r := range tc.wantLeft {
				if r != tc.queue[i] {
					t.Errorf("%s: mismatch at index %d, %d != %d", name, i, r, tc.queue[i])
				}
			}
		})
	}
}

func TestPriorityQueueSort(t *testing.T) {
	tests := map[string]struct {
		queue priorityQueue[int]
		want  priorityQueue[int]
	}{
		"no change": {queue: priorityQueue[int]{1, 2, 3}, want: priorityQueue[int]{1, 2, 3}},
		"sort":      {queue: priorityQueue[int]{2, 3, 1}, want: priorityQueue[int]{1, 2, 3}},
		"reverse":   {queue: priorityQueue[int]{5, 4, 3, 2, 1}, want: priorityQueue[int]{1, 2, 3, 4, 5}},
	}

	f := func(a int, b int) int {
		if a < b {
			return -1
		}
		return 1
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.queue.sortBy(f)

			for i, r := range tc.want {
				if r != tc.queue[i] {
					t.Errorf("%s: %d != %d", name, r, tc.queue[i])
				}
			}
		})
	}
}

func TestPriorityQueueInsert(t *testing.T) {
	tests := map[string]struct {
		queue priorityQueue[int]
		input int
		want  priorityQueue[int]
	}{
		"start":  {queue: priorityQueue[int]{1, 2, 3}, input: 1, want: priorityQueue[int]{1, 1, 2, 3}},
		"end":    {queue: priorityQueue[int]{1, 2, 3}, input: 5, want: priorityQueue[int]{1, 2, 3, 5}},
		"middle": {queue: priorityQueue[int]{1, 2, 3}, input: 2, want: priorityQueue[int]{1, 2, 2, 3}},
	}

	f := func(a int, b int) bool {
		return a < b
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.queue.insertAt(tc.input, f)
			if len(tc.queue) != len(tc.want) {
				t.Fatalf("%s: mismatching array lengths,  %d != %d", name, len(tc.queue), len(tc.want))
			}

			for i, r := range tc.want {
				if r != tc.queue[i] {
					t.Errorf("%s: mismatch at index %d,  %d != %d", name, i, r, tc.queue[i])
				}
			}
		})
	}
}

func TestPriorityQueueFetchBy(t *testing.T) {
	tests := map[string]struct {
		queue    priorityQueue[int]
		matcher  func(int) bool
		want     int
		wantOk   bool
		wantLeft priorityQueue[int]
	}{
		"found first": {
			queue:    priorityQueue[int]{1, 2, 3},
			matcher:  func(v int) bool { return v == 1 },
			want:     1,
			wantOk:   true,
			wantLeft: priorityQueue[int]{2, 3},
		},
		"found middle": {
			queue:    priorityQueue[int]{1, 2, 3},
			matcher:  func(v int) bool { return v == 2 },
			want:     2,
			wantOk:   true,
			wantLeft: priorityQueue[int]{1, 3},
		},
		"not found": {
			queue:    priorityQueue[int]{1, 2, 3},
			matcher:  func(v int) bool { return v == 4 },
			want:     0,
			wantOk:   false,
			wantLeft: priorityQueue[int]{1, 2, 3},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, ok := tc.queue.fetchBy(tc.matcher)
			if ok != tc.wantOk {
				t.Fatalf("%s: got ok=%v, want %v", name, ok, tc.wantOk)
			}
			if ok && res != tc.want {
				t.Fatalf("%s: got %d, want %d", name, res, tc.want)
			}
			if len(tc.queue) != len(tc.wantLeft) {
				t.Fatalf("%s: mismatching array lengths, %d != %d", name, len(tc.queue), len(tc.wantLeft))
			}
			for i, r := range tc.wantLeft {
				if r != tc.queue[i] {
					t.Errorf("%s: mismatch at index %d, %d != %d", name, i, r, tc.queue[i])
				}
			}
		})
	}
}
