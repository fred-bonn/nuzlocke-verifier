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

	_, err := parseInputFile(&cfg, path)
	if err != nil {
		fmt.Printf("error parsing input file: %v", err)
	}
}

func parseInputFile(cfg *Config, path string) ([]Pokemon, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	lines := strings.Split(string(data), "\n")

	// var pokemons []Pokemon

	for i := 0; i < len(lines); {
		// extract pokemon name (and item)
		line := lines[i]

		name := parsePokemonLine(line)
		if name == "" {
			return nil, fmt.Errorf("invalid pokemon name format: %q", line)
		}

		name = cleanPokemonName(name)
		i++

		// extract pokemon level
		line = lines[i]

		level := parseLevelLine(line)
		if level == 0 {
			return nil, fmt.Errorf("invalid level format: %q", line)
		}

		i++

		// extract pokemon nature
		line = lines[i]

		nature := parseNatureLine(line)
		if nature == "" {
			return nil, fmt.Errorf("invalid nature format: %q", line)
		}

		i++

		// skip ability for now
		i++

		// extract IVs map
		ivs := make(map[string]int, 6)
		if strings.HasPrefix(lines[i], "IVs: ") {
			line = lines[i]

			err := parseIVsLine(line, ivs)
			if err != nil {
				return nil, fmt.Errorf("invalid IVs format: %q", line)
			}

			i++
		}

		// extract moves slice
		var moves []string
		for strings.HasPrefix(lines[i], "-") {
			line = lines[i]

			move := parseMoveLine(line)
			if move == "" {
				return nil, fmt.Errorf("invalid move format: %q", line)
			}

			move = cleanMoveName(move)
			moves = append(moves, move)
			i++

			if i >= len(lines) {
				break
			}
		}

		// skip the final empty line before next Pokemon
		i++
	}

	return nil, nil
}

func parsePokemonLine(line string) string {
	parts := strings.Split(line, " @ ")
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

func parseIVsLine(line string, ivs map[string]int) error {
	line = strings.TrimPrefix(line, "IVs: ")
	parts := strings.Split(line, "/")

	for i := range parts {
		fields := strings.Fields(parts[i])
		if len(fields) != 2 {
			return fmt.Errorf("invalid IV format: %d", len(fields))
		}

		iv, err := strconv.Atoi(fields[0])
		if err != nil {
			return fmt.Errorf("not a number: %s", fields[0])
		}

		name := strings.ToLower(fields[1])

		ivs[name] = iv
	}
	return nil
}

func parseMoveLine(line string) string {
	if line == "" {
		return ""
	}

	parts := strings.Split(line, "-")
	if parts[0] != "" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
