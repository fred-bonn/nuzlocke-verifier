package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

type resistBerryEvent struct {
	typeName    string
	denominator *int
}

type gemEvent struct {
	typeName    string
	denominator *int
	numerator   *int
}

type leppaBerryEvent struct {
	move *pokeapi.BaseMove
}

type choiceItemEvent struct {
	move        *pokeapi.BaseMove
	denominator *int
	numerator   *int
}

type focusSashEvent struct {
	damage  *int
	consume bool
}
