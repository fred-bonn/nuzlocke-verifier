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

	bp, err := loadMove(&cfg, "thunder-punch")
	if err != nil {
		fmt.Printf("Error loading Move: %v\n", err)
		return
	}
	fmt.Printf("%s\n", bp.Name)
}
