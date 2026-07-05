package main

import (
	"log"
	"os"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
	"github.com/spf13/pflag"
)

var verbose = pflag.BoolP("verbose", "v", false, "verbose logging")

func main() {
	weather := pflag.IntP("weather", "w", int(noneWeather), "weather\n 0: None (default)\n 1: Rain\n 2: Sun\n 3: Sandstorm\n 4: Hail")
	pflag.Parse()
	if *weather < 0 || *weather > 4 {
		log.Printf("error: weather (-w) must be between 0 and 4")
		os.Exit(1)
	}

	args := pflag.Args()
	if len(args) != 2 {
		log.Printf("error: missing arguments: usage: <executable> <player_showdown> <opponent_showdown> <flags>")
		os.Exit(1)
	}

	cfg := &config{
		client: pokeapi.NewClient(),
	}

	playerParty, err := cfg.validateInput(args[0])
	if err != nil {
		log.Printf("error: failed validating input '%s': %s", args[0], err)
		os.Exit(1)
	}
	opponentParty, err := cfg.validateInput(args[1])
	if err != nil {
		log.Printf("error: failed validating input '%s': %s", args[1], err)
		os.Exit(1)
	}

	sbs := initSingleBattleState(
		trainer{
			ai:           rnbAi{},
			player:       true,
			fieldEffects: make(map[fieldEffect]int),
		},
		trainer{
			ai:           rnbAi{},
			fieldEffects: make(map[fieldEffect]int),
		},
		playerParty,
		opponentParty,
		weatherState(*weather),
	)

	err = sbs.execute()
	if err != nil {
		log.Printf("error: battle state execute failed: %s", err)
		os.Exit(1)
	}
}
