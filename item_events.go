package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

type resistBerryEvent struct {
	typeName    string
	denominator *int
}

type leppaBerryEvent struct {
	move *pokeapi.BaseMove
}
