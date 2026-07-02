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
		"paralysis":   {paralysisAilment, true},
		"none":        {noneAilment, false},
		"sleep":       {sleepAilment, true},
		"infatuation": {infatuationAilment, false},
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
		iterator      func(func(ailmentState) bool)
		isNonVolatile bool
	}{
		"non volatile": {nonVolatileStatuses, true},
		"volatile":     {volatileStatuses, false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for state := range tc.iterator {
				if got := state.isNonVolatileStatus(); got != tc.isNonVolatile {
					t.Errorf("%s := range %s, tc.%s.isNonVolatileStatus() = %t, want %t", state, FunctionName(tc.iterator), state, got, tc.isNonVolatile)
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
		"sleep":     {sleepAilment, sleepAilment, 1, 3},
		"confusion": {confusionAilment, confusionAilment, 1, 4},
		"yawn":      {yawnAilment, yawnAilment, 2, 2},
		"paralysis": {paralysisAilment, paralysisAilment, 0, 0},
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
		"1-5":   {1, 5},
		"2-2":   {2, 2},
		"0-1":   {0, 1},
		"0-0":   {0, 0},
		"99-99": {99, 99},
		"0-10":  {0, 10},
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
