package pokemon

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type MoveBalance struct {
	Power    *int
	Accuracy *int
	PP       *int
}

func (mb MoveBalance) Apply(m *pokeapi.BaseMove) {
	if mb.Power != nil {
		m.Power = *mb.Power
	}
	if mb.Accuracy != nil {
		m.Accuracy = *mb.Accuracy
	}
	if mb.PP != nil {
		m.PP = *mb.PP
	}
}

var MoveBalanceMap = map[string]*MoveBalance{
	"absorb": {
		Power: new(40),
	},
	"pin-missile": {
		Accuracy: new(100),
	},
}
