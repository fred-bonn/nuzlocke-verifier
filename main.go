package main

import (
	"fmt"
	"log"

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
	path := "./showdown_test_file.txt"
	res, err := parser.ReadShowdownFile(path)
	if err != nil {
		log.Fatalf("error: parsing '%s' failed: %v", err)
	}

	if len(res) == 0 {
		log.Fatalf("error: file '%s' contained no pokemon", path)
	}

	for i := range res {
		fmt.Println(res[i])
	}

	mons, err := loadShowdonw(cfg, res)
	if err != nil {
		log.Fatalf("error: loading showdown file '%s' failed: %v", path, err)
	}

	for i := range mons {
		fmt.Println(mons[i].Ailments)
	}
}
