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
	Drain         int
	Heal          int
	FlinchChange  int
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
