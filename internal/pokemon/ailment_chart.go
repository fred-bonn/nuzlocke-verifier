package pokemon

import (
	"math/rand"
)

const (
	INDEFINITE_AILMENT_DURATION = 200
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

func GenerateAilment(ailment string) int {
	switch ailment {
	case "sleep":
		return rand.Intn(3) + 1
	case "confusion":
		return rand.Intn(4) + 1
	default:
		return INDEFINITE_AILMENT_DURATION
	}
}

func GenerateTrap(low, high int) int {
	return rand.Intn(high-low+1) + low
}

func ValidAilment(ailment string) bool {
	if _, ok := validAilments[ailment]; ok {
		return true
	}
	return false
}
