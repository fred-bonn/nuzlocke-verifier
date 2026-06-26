package main

import (
	"math/rand"
)

type ailmentState int

const (
	Paralysis ailmentState = iota
	Poison
	Toxic
	Burn
	Freeze
	Sleep
	Infatuation
	Confusion
	Trap
	Bound
	LeechSeed
	Yawn
	NoneAilment
)

func stringToAilmentState(s string) ailmentState {
	switch s {
	case "paralysis":
		return Paralysis
	case "poison":
		return Poison
	case "toxic":
		return Toxic
	case "burn":
		return Burn
	case "freeze":
		return Freeze
	case "sleep":
		return Sleep
	case "infatuation":
		return Infatuation
	case "confusion":
		return Confusion
	case "trap":
		return Trap
	case "bound":
		return Bound
	case "leech-seed":
		return LeechSeed
	case "yawn":
		return Yawn
	default:
		return NoneAilment
	}
}

func (as ailmentState) String() string {
	switch as {
	case Paralysis:
		return "paralysis"
	case Poison:
		return "poison"
	case Toxic:
		return "toxic"
	case Burn:
		return "burn"
	case Freeze:
		return "freeze"
	case Sleep:
		return "sleep"
	case Infatuation:
		return "infatuation"
	case Confusion:
		return "confusion"
	case Trap:
		return "trap"
	case Bound:
		return "bound"
	case LeechSeed:
		return "leech-seed"
	case Yawn:
		return "yawn"
	default:
		return "invalid"
	}
}

type ailment struct {
	state       ailmentState
	turns       int
	afflictedBy *slot
}

var nonVolatileStatuses = map[ailmentState]struct{}{
	Paralysis: {},
	Poison:    {},
	Toxic:     {},
	Burn:      {},
	Freeze:    {},
	Sleep:     {},
}

var volatileStatuses = map[ailmentState]struct{}{
	Infatuation: {},
	Confusion:   {},
	Trap:        {},
	Bound:       {},
	LeechSeed:   {},
	Yawn:        {},
}

func generateAilment(as ailmentState, afflictedBy *slot) *ailment {
	res := ailment{
		state:       as,
		afflictedBy: afflictedBy,
	}
	switch as {
	case Sleep:
		res.turns = rand.Intn(3) + 1
	case Confusion:
		res.turns = rand.Intn(4) + 1
	case Yawn:
		res.turns = 2
	}
	return &res
}

func generateTrap(low, high int, mon *slot) *ailment {
	return &ailment{
		state:       Trap,
		turns:       rand.Intn(high-low+1) + low,
		afflictedBy: mon,
	}
}
