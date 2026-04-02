package main

import (
	"fmt"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type Config struct {
	client pokeapi.Client
}

func main() {
	cfg := &Config{
		client: pokeapi.NewClient(),
	}

	res, err := parser.ReadShowdownFile("./showdown_test_file.txt")
	if err != nil {
		fmt.Println(err.Error())
	}
	for i := range res {
		loadPokemon(cfg, res[i].Name)
	}
}
