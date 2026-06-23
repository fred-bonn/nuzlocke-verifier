package main

import (
	"math/rand"
)

type ailment struct {
	name        string
	turns       int
	afflictedBy *slot
}

var nonVolatileStatuses = map[string]struct{}{
	"paralysis": {},
	"poison":    {},
	"toxic":     {},
	"burn":      {},
	"freeze":    {},
	"sleep":     {},
}

var volatileStatuses = map[string]struct{}{
	"infatuation": {},
	"confusion":   {},
	"trap":        {},
	"bound":       {},
	"leech-seed":  {},
	"yawn":        {},
}

func generateAilment(ailmentName string, afflictedBy *slot) *ailment {
	res := ailment{
		name:        ailmentName,
		afflictedBy: afflictedBy,
	}
	switch ailmentName {
	case "sleep":
		res.turns = rand.Intn(3) + 1
	case "confusion":
		res.turns = rand.Intn(4) + 1
	case "yawn":
		res.turns = 2
	default:
		res.turns = 0
	}
	return &res
}

func generateTrap(low, high int, mon *slot) *ailment {
	return &ailment{
		name:        "trap",
		turns:       rand.Intn(high-low+1) + low,
		afflictedBy: mon,
	}
}
