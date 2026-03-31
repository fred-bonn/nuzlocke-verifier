package main

import (
	"fmt"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type Config struct {
	client pokeapi.Client
}

func main() {
	cfg := Config{
		client: pokeapi.NewClient(),
	}

	bp, err := loadPokemon(&cfg, "mewtwo")
	if err != nil {
		fmt.Printf("Error loading Pokémon: %v\n", err)
		return
	}

	p, err := initializePokemon(bp, 20, []int{31, 31, 31, 31, 31, 31}, "hardy", []Move{}, 25, "burn")
	if err != nil {
		fmt.Printf("Error initializing Pokémon: %v\n", err)
		return
	}
	fmt.Printf("%s\n", p.String())

}
