package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

func main() {
	cfg := &config{
		client: pokeapi.NewClient(),
	}
	path := "./showdown_test_file.txt"
	res, err := parser.ReadShowdownFile(path)
	if err != nil {
		log.Fatalf("error: parsing '%s' failed: %v", path, err)
	}

	if len(res) == 0 {
		log.Fatalf("error: file '%s' contained no pokemon", path)
	}

	_, err = cfg.loadShowdown(res)
	if err != nil {
		log.Fatalf("error: loading showdown file '%s' failed: %v", path, err)
	}
}
