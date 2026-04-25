package pokeapi

import (
	"encoding/json"
	"os"
)

type moveJSON struct {
	Name  string `json:"name"`
	Power int    `json:"power"`
	PP    int    `json:"pp"`
	Type  struct {
		Name string `json:"name"`
	} `json:"type"`
	Accuracy    int `json:"accuracy"`
	Priority    int `json:"priority"`
	DamageClass struct {
		Name string `json:"name"`
	} `json:"damage_class"`
	Meta struct {
		Ailment struct {
			Name string `json:"name"`
		} `json:"ailment"`
		AilmentChance int `json:"ailment_chance"`
		Category      struct {
			Name string `json:"name"`
		} `json:"category"`
		CritRate     int `json:"crit_rate"`
		Drain        int `json:"drain"`
		FlinchChance int `json:"flinch_chance"`
		Heal         int `json:"heal"`
		MaxHits      int `json:"max_hits"`
		MaxTurns     int `json:"max_turns"`
		MinHits      int `json:"min_hits"`
		MinTurns     int `json:"min_turns"`
		StatChance   int `json:"stat_chance"`
	} `json:"meta"`
	StatChanges []struct {
		Change int `json:"change"`
		Stat   struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stat_changes"`
	Target struct {
		Name string `json:"name"`
	} `json:"target"`
}

var contactMoves map[string]any

func (mj moveJSON) toMove() BaseMove {
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
	}

	if contactMoves == nil {
		initContactMoves()
	}

	_, isContact = contactMoves[mj.Name]

	return BaseMove{
		Name:          mj.Name,
		Type:          mj.Type.Name,
		Power:         mj.Power,
		Accuracy:      mj.Accuracy,
		Class:         mj.DamageClass.Name,
		Priority:      mj.Priority,
		CritRate:      mj.Meta.CritRate,
		Drain:         mj.Meta.Drain,
		Heal:          mj.Meta.Heal,
		FlinchChange:  mj.Meta.FlinchChance,
		Contact:       isContact,
		Ailentment:    mj.Meta.Ailment.Name,
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
