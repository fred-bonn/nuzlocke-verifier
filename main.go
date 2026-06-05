package main

import (
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

var verbose = flag.Bool("v", true, "verbose logging")

func main() {
	flag.Parse()
	args := flag.Args()
	fmt.Println(args)
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
		},
		trainer{
			ai: rnbAi{},
		},
		playerParty,
		opponentParty,
	)

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	sbs.execute()
}
