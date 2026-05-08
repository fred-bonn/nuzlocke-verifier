package pokemon

import (
	"math/rand"
)

var validAilments = map[string]struct{}{
	"paralysis": {},
	"poison":    {},
	"toxic":     {},
	"burn":      {},
	"freeze":    {},
	"confusion": {},
	"trap":      {},
	"bound":     {},
}

var nonVolatileStatuses = map[string]struct{}{
	"paralysis": {},
	"poison":    {},
	"toxic":     {},
	"burn":      {},
	"freeze":    {},
}

var volatileStatuses = map[string]struct{}{
	"confusion": {},
	"trap":      {},
	"bound":     {},
}

func GenerateAilment(ailment string) int {
	switch ailment {
	case "sleep":
		return rand.Intn(3) + 1
	case "confusion":
		return rand.Intn(4) + 1
	default:
		return 0
	}
}

func GenerateTrap(low, high int) int {
	return rand.Intn(high-low+1) + low
}
