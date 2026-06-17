package main

import "testing"

func TestPriorityQueuePush(t *testing.T) {
	tests := map[string]struct {
		queue    priorityQueue[int]
		push     []int
		wantLeft priorityQueue[int]
	}{
		"empty":          {priorityQueue[int]{}, []int{1, 2, 3}, priorityQueue[int]{1, 2, 3}},
		"three elements": {priorityQueue[int]{1, 2, 3}, []int{1, 2, 3}, priorityQueue[int]{1, 2, 3, 1, 2, 3}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, element := range tc.push {
				tc.queue.push(element)
			}
			if n := len(tc.queue); n != len(tc.wantLeft) {
				t.Fatalf("%s: tc.queue.len() = %d, want %d", name, n, len(tc.wantLeft))
			}
			for i, r := range tc.wantLeft {
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
		"empty":          {priorityQueue[int]{}, 0, false, priorityQueue[int]{}},
		"three elements": {priorityQueue[int]{1, 2, 3}, 3, true, priorityQueue[int]{1, 2}},
		"one element":    {priorityQueue[int]{1}, 1, true, priorityQueue[int]{}},
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
		"no change": {priorityQueue[int]{1, 2, 3}, priorityQueue[int]{1, 2, 3}},
		"sort":      {priorityQueue[int]{2, 3, 1}, priorityQueue[int]{1, 2, 3}},
		"reverse":   {priorityQueue[int]{5, 4, 3, 2, 1}, priorityQueue[int]{1, 2, 3, 4, 5}},
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
		"start":  {priorityQueue[int]{1, 2, 3}, 1, priorityQueue[int]{1, 1, 2, 3}},
		"end":    {priorityQueue[int]{1, 2, 3}, 5, priorityQueue[int]{1, 2, 3, 5}},
		"middle": {priorityQueue[int]{1, 2, 3}, 2, priorityQueue[int]{1, 2, 2, 3}},
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
