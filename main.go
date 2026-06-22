package main

import (
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/spf13/pflag"
)

var verbose = pflag.BoolP("verbose", "v", false, "verbose logging")

func main() {
	weather := pflag.IntP("weather", "w", 0, "weather\n 0: None\n 1: Rain\n 2: Sun\n 3: Sandstorm\n 4: Hail")
	pflag.Parse()
	if *weather < 0 || *weather > 4 {
		log.Fatalf("error: weather (-w) must be tween 0 and 4")
	}

	args := pflag.Args()
	if len(args) != 2 {
		log.Fatalf("error: missing arguments: usage: <executable> <player_showdown> <opponent_showdown> <flags>")
	}

	cfg := &config{
		client: pokeapi.NewClient(),
	}

	playerParty, err := cfg.validateInput(args[0])
	if err != nil {
		log.Fatalf("error: failed validating input '%s': %s", args[0], err)
	}
	opponentParty, err := cfg.validateInput(args[1])
	if err != nil {
		log.Fatalf("error: failed validating input '%s': %s", args[1], err)
	}

	sbs := initSingleBattleState(
		trainer{
			ai:     rnbAi{},
			player: true,
			field:  map[string]struct{}{},
		},
		trainer{
			ai:    rnbAi{},
			field: map[string]struct{}{},
		},
		playerParty,
		opponentParty,
		weatherState(*weather),
	)

	sbs.execute()
}
