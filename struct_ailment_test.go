package main

import "testing"

func TestStringToAilmentState(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ailmentState
	}{
		{name: "paralysis", input: "paralysis", want: paralysisAilment},
		{name: "poison", input: "poison", want: poisonAilment},
		{name: "toxic", input: "toxic", want: toxicAilment},
		{name: "burn", input: "burn", want: burnAilment},
		{name: "freeze", input: "freeze", want: freezeAilment},
		{name: "sleep", input: "sleep", want: sleepAilment},
		{name: "infatuation", input: "infatuation", want: infatuationAilment},
		{name: "confusion", input: "confusion", want: confusionAilment},
		{name: "trap", input: "trap", want: trapAilment},
		{name: "bound", input: "bound", want: boundAilment},
		{name: "leech-seed", input: "leech-seed", want: leechSeedAilment},
		{name: "yawn", input: "yawn", want: yawnAilment},
		{name: "unknown", input: "unknown", want: noneAilment},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringToAilmentState(tt.input); got != tt.want {
				t.Fatalf("stringToAilmentState(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
