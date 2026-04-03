package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

func loadShowdown(cfg *Config, mons []parser.ParsedPokemon) ([]pokemon.Pokemon, error) {
	var res []pokemon.Pokemon

	for _, mon := range mons {
		var moves []pokeapi.BaseMove

		cleanedMonName := cleanPokemonName(mon.Name)
		basePokemon, err := loadPokemon(cfg, cleanedMonName)
		if err != nil {
			return nil, err
		}

		for _, moveName := range mon.Moves {
			cleanedMoveName := cleanMoveName(moveName)
			baseMove, err := loadMove(cfg, cleanedMoveName)
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

func loadPokemon(cfg *Config, name string) (pokeapi.BasePokemon, error) {
	data, err := os.ReadFile(fmt.Sprintf("data/pokemon/%s.json", name))
	// If the file exists and is read successfully, unmarshal it into a Pokemon struct
	if err == nil {
		var pokemon pokeapi.BasePokemon
		err = json.Unmarshal(data, &pokemon)
		if err != nil {
			return pokeapi.BasePokemon{}, fmt.Errorf("failed unmarshaling '%s' Pokemon data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
		return pokemon, nil
	}

	// Otherwise, fetch the Pokemon data from the API
	pokemon, err := cfg.client.FetchPokemon(name)
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

	return pokemon, nil
}

func loadMove(cfg *Config, name string) (pokeapi.BaseMove, error) {
	data, err := os.ReadFile(fmt.Sprintf("data/moves/%s.json", name))
	// If the file exists and is read successfully, unmarshal it into a Move struct
	if err == nil {
		var move pokeapi.BaseMove
		err = json.Unmarshal(data, &move)
		if err != nil {
			return pokeapi.BaseMove{}, fmt.Errorf("failed unmarshaling Move '%s' data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
		return move, nil
	}

	// Otherwise, fetch the Move data from the API
	move, err := cfg.client.FetchMove(name)
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
