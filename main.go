package main

import (
	"log"
	"os"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatalf("error: missing arguments: usage: %s <player_showdown> <opponent_showdown>", args[0])
	}

	cfg := &config{
		client: pokeapi.NewClient(),
	}

	playerParty, err := cfg.validateInput(args[1])
	if err != nil {
		log.Fatalf("error: failed validating input '%s': %s", args[1], err)
	}
	opponentParty, err := cfg.validateInput(args[2])
	if err != nil {
		log.Fatalf("error: failed validating input '%s': %s", args[2], err)
	}

	sbs, err := initSingleBattleState(trainer{
		pokemonParty: playerParty,
		ai:           randomAi{},
		player:       true,
	}, trainer{
		pokemonParty: opponentParty,
		ai:           randomAi{},
	})
	if err != nil {
		log.Fatalf("error: failed initializing battle state: %s", err)
	}

	sbs.execute()
}
