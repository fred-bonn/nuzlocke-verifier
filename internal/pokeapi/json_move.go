package pokeapi

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

func (mj moveJSON) toMove() BaseMove {
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

	return BaseMove{
		Name:          mj.Name,
		Type:          mj.Type.Name,
		Power:         mj.Power,
		Accuracy:      mj.Accuracy,
		Class:         mj.DamageClass.Name,
		Priority:      mj.Priority,
		Drain:         mj.Meta.Drain,
		Heal:          mj.Meta.Heal,
		FlinchChange:  mj.Meta.FlinchChance,
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
