package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokemon"
)

type action interface {
	invoke(bs *singleBattleState)
	prio() int
	speed() int
}

type swapAction struct {
	old *pokemon.Pokemon
	new *pokemon.Pokemon
}

func (sa *swapAction) invoke(sbs *singleBattleState) {
	sbs.setMon(sa.old, sa.new)
	log.Printf("swapped %s for %s", sa.old.Base.Name, sa.new.Base.Name)
}

func (sa *swapAction) prio() int {
	return 6
}

func (sa *swapAction) speed() int {
	return sa.old.Base.Stats["speed"]
}

type moveAction struct {
	mon  *pokemon.Pokemon
	slot int
	move pokeapi.BaseMove
}

func (ma *moveAction) prio() int {
	return ma.move.Priority
}

func (ma *moveAction) speed() int {
	return ma.mon.Stats["speed"]
}

func (ma *moveAction) invoke(sbs *singleBattleState) {
	if ma.mon.Fainted {
		return
	}
	log.Printf("%s doing move %s on %d", ma.mon.Base.Name, ma.move.Name, ma.slot)
}
