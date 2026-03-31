package pokeapi

type BasePokemon struct {
	Id     int
	Name   string
	Height int
	Weight int
	Types  []string
	Stats  map[string]int
}
