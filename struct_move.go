package main

import (
	"encoding/json"
	"os"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type moveClass int

const (
	physicalClass moveClass = iota
	specialClass
	statusClass
	noneClass
)

func stringToMoveClass(s string) moveClass {
	switch s {
	case "physical":
		return physicalClass
	case "special":
		return specialClass
	case "status":
		return statusClass
	default:
		elogFatalf("error: %s is not a valid move class", s)
		return noneClass
	}
}

type Move struct {
	Name          string
	Type          pokemonType
	Power         int
	Accuracy      int
	PP            int
	MaxPP         int
	Class         moveClass
	Priority      int
	CritRate      int
	Drain         int
	Heal          int
	FlinchChance  int
	Contact       bool
	Ailment       ailmentState
	AilmentChance int
	MaxHits       int
	MinHits       int
	MaxTurns      int
	MinTurns      int
	StatChance    int
	StatChanges   map[string]int
	Target        string
	Category      string
}

var contactMoves map[string]any

func toMove(mj pokeapi.MoveJSON) Move {
	isContact := false
	statChanges := make(map[string]int)
	for _, sc := range mj.StatChanges {
		statChanges[sc.Stat.Name] = sc.Change
	}
	var statChance int
	var ailmentChance int
	if mj.DamageClass.Name == "status" {
		statChance = 100
		ailmentChance = 100
	} else {
		statChance = mj.Meta.StatChance
		ailmentChance = mj.Meta.AilmentChance
	}

	if contactMoves == nil {
		initContactMoves()
	}

	_, isContact = contactMoves[mj.Name]

	return Move{
		Name:          mj.Name,
		Type:          stringToPokemonType(mj.Type.Name),
		Power:         mj.Power,
		Accuracy:      mj.Accuracy,
		PP:            mj.PP,
		MaxPP:         mj.PP,
		Class:         stringToMoveClass(mj.DamageClass.Name),
		Priority:      mj.Priority,
		CritRate:      mj.Meta.CritRate,
		Drain:         mj.Meta.Drain,
		Heal:          mj.Meta.Heal,
		FlinchChance:  mj.Meta.FlinchChance,
		Contact:       isContact,
		Ailment:       stringToAilmentState(mj.Meta.Ailment.Name),
		AilmentChance: ailmentChance,
		MaxHits:       mj.Meta.MaxHits,
		MinHits:       mj.Meta.MinHits,
		MaxTurns:      mj.Meta.MaxTurns,
		MinTurns:      mj.Meta.MinTurns,
		StatChance:    statChance,
		StatChanges:   statChanges,
		Target:        mj.Target.Name,
		Category:      mj.Meta.Category.Name,
	}
}

func initContactMoves() error {
	var moves []string

	data, err := os.ReadFile("./contact_moves.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &moves)
	if err != nil {
		return err
	}

	contactMoves = make(map[string]any, len(moves))
	for _, move := range moves {
		contactMoves[move] = struct{}{}
	}

	return nil
}
