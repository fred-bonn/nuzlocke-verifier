package main

import (
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestStringToAilmentState(t *testing.T) {
	tests := map[string]struct {
		input string
		want  ailmentState
	}{
		"paralysis":   {input: "paralysis", want: paralysisAilment},
		"poison":      {input: "poison", want: poisonAilment},
		"toxic":       {input: "toxic", want: toxicAilment},
		"burn":        {input: "burn", want: burnAilment},
		"freeze":      {input: "freeze", want: freezeAilment},
		"sleep":       {input: "sleep", want: sleepAilment},
		"infatuation": {input: "infatuation", want: infatuationAilment},
		"confusion":   {input: "confusion", want: confusionAilment},
		"trap":        {input: "trap", want: trapAilment},
		"bound":       {input: "bound", want: boundAilment},
		"leech seed":  {input: "leech seed", want: leechSeedAilment},
		"yawn":        {input: "yawn", want: yawnAilment},
		"unknown":     {input: "", want: noneAilment},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := stringToAilmentState(tc.input); got != tc.want {
				t.Errorf("stringToAilmentState(%s) = %s, want %s", tc.input, got, tc.want)
			}
		})
	}
}

func TestIsNonVolatileStatus(t *testing.T) {
	tests := map[string]struct {
		state ailmentState
		want  bool
	}{
		"paralysis":   {state: paralysisAilment, want: true},
		"none":        {state: noneAilment, want: false},
		"sleep":       {state: sleepAilment, want: true},
		"infatuation": {state: infatuationAilment, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.state.isNonVolatileStatus(); got != tc.want {
				t.Errorf("tc.%s.isNonVolatileStatus() = %t, want %t", tc.state, got, tc.want)
			}
		})
	}
}

func TestAilmentIterators(t *testing.T) {
	tests := map[string]struct {
		iterator func(func(ailmentState) bool)
		want     bool
	}{
		"non volatile": {iterator: nonVolatileStatuses, want: true},
		"volatile":     {iterator: volatileStatuses, want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for state := range tc.iterator {
				if got := state.isNonVolatileStatus(); got != tc.want {
					t.Errorf("%s := range %s, tc.%s.isNonVolatileStatus() = %t, want %t", state, FunctionName(tc.iterator), state, got, tc.want)
				}
			}
		})
	}
}

const (
	generateAilmentIterations int = 500
)

func TestGenerateAilment(t *testing.T) {
	tests := map[string]struct {
		state    ailmentState
		want     ailmentState
		minTurns int
		maxTurns int
	}{
		"sleep":     {state: sleepAilment, want: sleepAilment, minTurns: 1, maxTurns: 3},
		"confusion": {state: confusionAilment, want: confusionAilment, minTurns: 1, maxTurns: 4},
		"yawn":      {state: yawnAilment, want: yawnAilment, minTurns: 2, maxTurns: 2},
		"paralysis": {state: paralysisAilment, want: paralysisAilment, minTurns: 0, maxTurns: 0},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i := 0; i < generateAilmentIterations; i++ {
				if got := generateAilment(tc.state, nil); got.state != tc.want || got.turns < tc.minTurns || got.turns > tc.maxTurns {
					t.Fatalf("generateAilment(%s, _) = {state: %s, turns: %d, _}, want {state: %s, turns: %d-%d, _}", tc.state, got.state, got.turns, tc.want, tc.minTurns, tc.maxTurns)
				}
			}
		})
	}
}

func TestGenerateTrap(t *testing.T) {
	tests := map[string]struct {
		minTurns int
		maxTurns int
	}{
		"1-5":   {minTurns: 1, maxTurns: 5},
		"2-2":   {minTurns: 2, maxTurns: 2},
		"0-1":   {minTurns: 0, maxTurns: 1},
		"0-0":   {minTurns: 0, maxTurns: 0},
		"99-99": {minTurns: 99, maxTurns: 99},
		"0-10":  {minTurns: 0, maxTurns: 10},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for i := 0; i < generateAilmentIterations; i++ {
				if got := generateTrap(tc.minTurns, tc.maxTurns, nil); got.state != trapAilment || got.turns < tc.minTurns || got.turns > tc.maxTurns {
					t.Fatalf("generateTrap(%d, %d, _) = {state: %s, turns: %d, _}, want {state: %s, turns: %d-%d, _}", tc.minTurns, tc.maxTurns, got.state, got.turns, trapAilment, tc.minTurns, tc.maxTurns)
				}
			}
		})
	}
}

func FunctionName(fn interface{}) string {
	full := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	full = path.Base(full)
	if i := strings.LastIndex(full, "."); i != -1 {
		return full[i+1:]
	}
	return full
}
