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

	_, err := loadPokemon(&cfg, "ivysaur")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = loadPokemon(&cfg, "pidgey")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	_, err = loadPokemon(&cfg, "mewtwo")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}
