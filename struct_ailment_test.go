package main

import "testing"

func TestStringToAilmentState(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ailmentState
	}{
		{name: "paralysis", input: "paralysis", want: Paralysis},
		{name: "poison", input: "poison", want: Poison},
		{name: "toxic", input: "toxic", want: Toxic},
		{name: "burn", input: "burn", want: Burn},
		{name: "freeze", input: "freeze", want: Freeze},
		{name: "sleep", input: "sleep", want: Sleep},
		{name: "infatuation", input: "infatuation", want: Infatuation},
		{name: "confusion", input: "confusion", want: Confusion},
		{name: "trap", input: "trap", want: Trap},
		{name: "bound", input: "bound", want: Bound},
		{name: "leech-seed", input: "leech-seed", want: LeechSeed},
		{name: "yawn", input: "yawn", want: Yawn},
		{name: "unknown", input: "unknown", want: Invalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToAilmentState(tt.input); got != tt.want {
				t.Fatalf("stringToAilmentState(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
