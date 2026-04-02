package main

import (
	"math/rand"
)

const (
	INDEFINITE_AILMENT_DURATION = 100
)

var validAilments = map[string]struct{}{
	"paralysis": {},
	"poison":    {},
	"toxic":     {},
	"burn":      {},
	"freeze":    {},
	"sleep":     {},
	"confusion": {},
}

func generateAilment(ailment string) int {
	switch ailment {
	case "sleep":
		return rand.Intn(3) + 1
	case "confusion":
		return rand.Intn(4) + 1
	default:
		return INDEFINITE_AILMENT_DURATION
	}
}
