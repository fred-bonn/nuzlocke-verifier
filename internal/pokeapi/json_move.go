package pokeapi

type MoveJSON struct {
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
		Heal         int `json:"healing"`
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
