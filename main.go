package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fred-bonn/nuzlocke_verifier/internal/pokeapi"
)

type Config struct {
	client pokeapi.Client
}

func main() {
	cfg := Config{
		client: pokeapi.NewClient(),
	}

	pokemon, err := cfg.client.FetchPokemon("ivysaur")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	data, err := json.Marshal(pokemon)
	if err != nil {
		fmt.Println("Error marshaling Pokemon data:", err)
		return
	}

	writeToFile("data/pokemon/ivysaur.json", data)
	fmt.Printf("Pokemon: %+v\n", data)
}

func writeToFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}
