package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

type BasePokemon struct {
	Id     int
	Name   string
	Height int
	Weight int
	Types  []pokemonType
	Stats  map[string]int
}

func toPokemon(pj pokeapi.PokemonJSON) BasePokemon {
	types := make([]pokemonType, len(pj.Types))
	for i, t := range pj.Types {
		types[i] = stringToPokemonType(t.Type.Name)
	}

	stats := make(map[string]int, 6)
	for _, s := range pj.Stats {
		stats[s.Stat.Name] = s.BaseStat
	}

	return BasePokemon{
		Id:     pj.Id,
		Name:   pj.Name,
		Height: pj.Height,
		Weight: pj.Weight,
		Types:  types,
		Stats:  stats,
	}
}
