package main

import (
	"math/rand"
)

type Ailment struct {
	Name        string
	Turns       int
	AfflictedBy *slot
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
}

func generateAilment(ailment string, afflictedBy *slot) *Ailment {
	res := Ailment{
		Name:        ailment,
		AfflictedBy: afflictedBy,
	}
	switch ailment {
	case "sleep":
		res.Turns = rand.Intn(3) + 1
	case "confusion":
		res.Turns = rand.Intn(4) + 1
	default:
		res.Turns = 0
	}
	return &res
}

func generateTrap(low, high int, mon *slot) *Ailment {
	return &Ailment{
		Name:        "trap",
		Turns:       rand.Intn(high-low+1) + low,
		AfflictedBy: mon,
	}
}
