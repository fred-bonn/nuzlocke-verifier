package pokeapi

type BasePokemon struct {
	Id     int
	Name   string
	Height int
	Weight int
	Types  []string
	Stats  map[string]int
}

type BaseMove struct {
	Name          string
	Type          string
	Power         int
	Accuracy      int
	Class         string
	Priority      int
	CritRate      int
	Drain         int
	Heal          int
	FlinchChange  int
	Contact       bool
	Ailentment    string
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
