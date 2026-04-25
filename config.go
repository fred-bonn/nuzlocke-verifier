package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type config struct {
	client pokeapi.Client
}

func (cfg *config) loadShowdown(mons []parser.ParsedPokemon) ([]pokemon.Pokemon, error) {
	var res []pokemon.Pokemon

	for _, mon := range mons {
		var moves []pokeapi.BaseMove

		cleanedMonName := cleanName(mon.Name)
		basePokemon, err := cfg.loadPokemon(cleanedMonName)
		if err != nil {
			return nil, err
		}

		for _, moveName := range mon.Moves {
			cleanedMoveName := cleanName(moveName)
			baseMove, err := cfg.loadMove(cleanedMoveName)
			if err != nil {
				return nil, err
			}

			moves = append(moves, baseMove)
		}

		finalPokemon, err := pokemon.InitializePokemon(basePokemon, mon.Level, mon.IVs, mon.Nature, moves, mon.HP, mon.Status)
		if err != nil {
			return nil, err
		}

		res = append(res, finalPokemon)
	}

	return res, nil
}

func (cfg *config) loadPokemon(name string) (pokeapi.BasePokemon, error) {
	var pokemon pokeapi.BasePokemon

	data, err := os.ReadFile(fmt.Sprintf("data/pokemon/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Pokemon struct
		err = json.Unmarshal(data, &pokemon)
		if err != nil {
			return pokeapi.BasePokemon{}, fmt.Errorf("failed unmarshaling '%s' Pokemon data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
	} else {
		// Otherwise, fetch the Pokemon data from the API
		pokemon, err = cfg.client.FetchPokemon(name)
		if err != nil {
			return pokeapi.BasePokemon{}, fmt.Errorf("failed fetching Pokemon '%s': %w", name, err)
		}
		fmt.Printf("Fetched '%s' from API\n", name)

		// Save the fetched Pokemon data using the internal Pokemon struct to a file for future use
		data, err = json.Marshal(pokemon)
		if err != nil {
			return pokeapi.BasePokemon{}, fmt.Errorf("failed marshaling Pokemon JSON data '%s' to file: %w", name, err)
		}
		writeToFile(fmt.Sprintf("data/pokemon/%s.json", name), data)
	}

	return pokemon, nil
}

func (cfg *config) loadMove(name string) (pokeapi.BaseMove, error) {
	var move pokeapi.BaseMove

	if strings.HasPrefix(name, "hidden-power") {
		// If the move is Hidden Power, generate it
		var err error
		move, err = generateHiddenPower(name)
		if err != nil {
			return pokeapi.BaseMove{}, err
		}
		return move, nil
	}

	data, err := os.ReadFile(fmt.Sprintf("data/moves/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Move struct
		err = json.Unmarshal(data, &move)
		if err != nil {
			return pokeapi.BaseMove{}, fmt.Errorf("failed unmarshaling Move '%s' data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
	} else {
		// Otherwise, fetch the Move data from the API
		move, err = cfg.client.FetchMove(name)
		if err != nil {
			return pokeapi.BaseMove{}, fmt.Errorf("failed fetching Move '%s': %w", name, err)
		}
		fmt.Printf("Fetched '%s' from API\n", name)

		// Save the fetched Move data using the internal Move struct to a file for future use
		data, err = json.Marshal(move)
		if err != nil {
			return pokeapi.BaseMove{}, fmt.Errorf("failed marshaling Move JSON data '%s' to file: %w", name, err)
		}
		writeToFile(fmt.Sprintf("data/moves/%s.json", name), data)
	}

	return move, nil
}

func generateHiddenPower(name string) (pokeapi.BaseMove, error) {
	parts := strings.Split(name, "-")
	if len(parts) != 3 {
		return pokeapi.BaseMove{}, fmt.Errorf("type not specified for hidden power")
	}
	if _, ok := pokemon.TypeChart[parts[2]]; !ok {
		return pokeapi.BaseMove{}, fmt.Errorf("invalid type for hidden power")
	}

	move := pokeapi.BaseMove{
		Name:     "hidden-power",
		Type:     parts[2],
		Power:    60,
		Accuracy: 100,
		Class:    "special",
	}

	return move, nil
}

func writeToFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func cleanName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "’", "")
	return name
}
