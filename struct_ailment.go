package main

import (
	"math/rand"
)

type ailmentState int

const (
	paralysisAilment ailmentState = iota
	poisonAilment
	toxicAilment
	burnAilment
	freezeAilment
	sleepAilment
	infatuationAilment
	confusionAilment
	trapAilment
	boundAilment
	leechSeedAilment
	yawnAilment
	noneAilment
)

func stringToAilmentState(s string) ailmentState {
	switch s {
	case "paralysis":
		return paralysisAilment
	case "poison":
		return poisonAilment
	case "toxic":
		return toxicAilment
	case "burn":
		return burnAilment
	case "freeze":
		return freezeAilment
	case "sleep":
		return sleepAilment
	case "infatuation":
		return infatuationAilment
	case "confusion":
		return confusionAilment
	case "trap":
		return trapAilment
	case "bound":
		return boundAilment
	case "leech-seed":
		return leechSeedAilment
	case "yawn":
		return yawnAilment
	default:
		return noneAilment
	}
}

func (as ailmentState) String() string {
	switch as {
	case paralysisAilment:
		return "paralysis"
	case poisonAilment:
		return "poison"
	case toxicAilment:
		return "toxic"
	case burnAilment:
		return "burn"
	case freezeAilment:
		return "freeze"
	case sleepAilment:
		return "sleep"
	case infatuationAilment:
		return "infatuation"
	case confusionAilment:
		return "confusion"
	case trapAilment:
		return "trap"
	case boundAilment:
		return "bound"
	case leechSeedAilment:
		return "leech-seed"
	case yawnAilment:
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
	paralysisAilment: {},
	poisonAilment:    {},
	toxicAilment:     {},
	burnAilment:      {},
	freezeAilment:    {},
	sleepAilment:     {},
}

var volatileStatuses = map[ailmentState]struct{}{
	infatuationAilment: {},
	confusionAilment:   {},
	trapAilment:        {},
	boundAilment:       {},
	leechSeedAilment:   {},
	yawnAilment:        {},
}

func generateAilment(as ailmentState, afflictedBy *slot) *ailment {
	res := ailment{
		state:       as,
		afflictedBy: afflictedBy,
	}
	switch as {
	case sleepAilment:
		res.turns = rand.Intn(3) + 1
	case confusionAilment:
		res.turns = rand.Intn(4) + 1
	case yawnAilment:
		res.turns = 2
	}
	return &res
}

func generateTrap(low, high int, mon *slot) *ailment {
	return &ailment{
		state:       trapAilment,
		turns:       rand.Intn(high-low+1) + low,
		afflictedBy: mon,
	}
}
