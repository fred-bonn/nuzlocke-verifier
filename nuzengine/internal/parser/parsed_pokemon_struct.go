package parser

import "fmt"

type ParsedPokemon struct {
	Name    string
	Item    string
	Level   int
	Nature  string
	Ability string
	Status  string
	HP      int
	IVs     map[string]int
	Moves   []string
}

func (p ParsedPokemon) String() string {
	return fmt.Sprintf("Name: %s, Item: %s, Level: %d, Nature: %s, Ability: %s Status: %s, HP: %d, IVs: %v, Moves: %v", p.Name, p.Item, p.Level, p.Nature, p.Ability, p.Status, p.HP, p.IVs, p.Moves)
}
