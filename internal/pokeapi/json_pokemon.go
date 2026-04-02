package pokeapi

// This file contains the structs used to unmarshal the JSON data from the PokeAPI.
type pokemonJSON struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		}
	}
}

// Converting PokeAPI JSON to internal Pokemon struct
func (pj pokemonJSON) toPokemon() BasePokemon {
	types := make([]string, len(pj.Types))
	for i, t := range pj.Types {
		types[i] = t.Type.Name
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
