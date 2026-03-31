package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type Config struct {
	client pokeapi.Client
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Usage: nuzlocke-verifier [options]")
		fmt.Println("Options:")
		fmt.Println("...")
		return
	}

	cfg := Config{
		client: pokeapi.NewClient(),
	}

	path := args[0]

	parseInputFile(&cfg, path)
}

func parseInputFile(cfg *Config, path string) ([]Pokemon, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	// var pokemons []Pokemon

	for i := 0; i < len(lines); {
		line := lines[i]
		name := cleanPokemonName(parsePokemonLine(line))
		if name == "" {
			return nil, fmt.Errorf("invalid pokemon name format: %q", line)
		}
		fmt.Printf("Parsed Pokémon: %s\n", name)
		i++

		line = lines[i]
		level := parseLevelLine(line)
		if level == 0 {
			return nil, fmt.Errorf("invalid level format: %q", line)
		}
		fmt.Printf("Parsed Level: %d\n", level)
		i++

		line = lines[i]
		nature := parseNatureLine(line)
		if nature == "" {
			return nil, fmt.Errorf("invalid nature format: %q", line)
		}
		fmt.Printf("Parsed Nature: %s\n", nature)
		i++

		// skip ability for now
		i++

		var moves []string
		for strings.HasPrefix(lines[i], "-") {
			line = lines[i]
			move := parseMoveLine(line)
			if move == "" {
				return nil, fmt.Errorf("invalid move format: %q", line)
			}
			fmt.Printf("Parsed Move: %s\n", move)
			moves = append(moves, move)
			i++
		}

		i++
	}

	return nil, nil
}

func parsePokemonLine(line string) string {
	parts := strings.Split(line, " @ ")
	if len(parts) != 2 {
		return ""
	}
	name := strings.TrimSpace(parts[0])
	return name
}

func parseLevelLine(line string) int {
	parts := strings.Split(line, ":")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	if parts[0] != "Level" {
		return 0
	}

	level, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	return level
}

func parseNatureLine(line string) string {
	parts := strings.Fields(line)

	if len(parts) != 2 {
		return ""
	}

	if parts[1] != "Nature" {
		return ""
	}

	nature := strings.ToLower(parts[0])

	_, ok := natureChart[nature]
	if !ok {
		return ""
	}

	return nature
}

func parseMoveLine(line string) string {
	return ""
}
