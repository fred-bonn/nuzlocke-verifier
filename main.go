package main

import (
	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type Config struct {
	client pokeapi.Client
}

func main() {
	parser.Parse()

	/*
		args := os.Args[1:]

		if len(args) == 0 {
			fmt.Println("Usage: nuzlocke-verifier [options]")
			fmt.Println("Options:")
			fmt.Println("...")
			return
		}

		cfg := Config{
			client: pokeapi.NewClient(),
		}

		path := args[0]

		pokemon, err := parseInputFile(&cfg, path)
		if err != nil {
			fmt.Printf("error parsing input file: %v", err)
		}
		fmt.Println("---")
		for i := range pokemon {
			fmt.Println(pokemon[i].String())
			fmt.Println("---")
		}
	*/
}
