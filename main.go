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
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err.Error())
	}
}
