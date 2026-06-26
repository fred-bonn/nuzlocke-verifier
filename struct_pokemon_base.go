package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

type basePokemon struct {
	Id     int
	Name   string
	Height int
	Weight int
	Types  []string
	Stats  map[string]int
}

func toPokemon(pj pokeapi.PokemonJSON) basePokemon {
	types := make([]string, len(pj.Types))
	for i, t := range pj.Types {
		types[i] = t.Type.Name
	}

	stats := make(map[string]int, 6)
	for _, s := range pj.Stats {
		stats[s.Stat.Name] = s.BaseStat
	}

	return basePokemon{
		Id:     pj.Id,
		Name:   pj.Name,
		Height: pj.Height,
		Weight: pj.Weight,
		Types:  types,
		Stats:  stats,
	}
}
