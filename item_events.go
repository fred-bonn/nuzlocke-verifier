package main

import "github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"

type resistBerryEvent struct {
	typeName string
	damage   *int
}

type gemEvent struct {
	typeName string
	power    *int
}

type leppaBerryEvent struct {
	move *pokeapi.BaseMove
}

type choiceItemEvent struct {
	move *pokeapi.BaseMove
	stat *int
}

type focusSashEvent struct {
	damage *int
}

type moveBoostingEvent struct {
	power    *int
	typeName string
}
