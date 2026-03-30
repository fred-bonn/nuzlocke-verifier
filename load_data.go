package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

func loadPokemon(cfg *Config, name string) (pokeapi.Pokemon, error) {
	data, err := os.ReadFile(fmt.Sprintf("data/pokemon/%s.json", name))
	// If the file exists and is read successfully, unmarshal it into a Pokemon struct
	if err == nil {
		var pokemon pokeapi.Pokemon
		err = json.Unmarshal(data, &pokemon)
		if err != nil {
			return pokeapi.Pokemon{}, fmt.Errorf("error unmarshaling Pokemon data: %w", err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
		return pokemon, nil
	}

	// Otherwise, fetch the Pokemon data from the API
	pokemon, err := cfg.client.FetchPokemon(name)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("error fetching Pokemon: %w", err)
	}
	fmt.Printf("Fetched '%s' from API\n", name)

	// Save the fetched Pokemon data using the internal Pokemon struct to a file for future use
	data, err = json.Marshal(pokemon)
	if err != nil {
		return pokeapi.Pokemon{}, fmt.Errorf("error marshaling Pokemon JSON data to file: %w", err)
	}
	writeToFile(fmt.Sprintf("data/pokemon/%s.json", name), data)

	return pokemon, nil
}

func writeToFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}
