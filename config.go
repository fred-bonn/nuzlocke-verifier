package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fred-bonn/nuz/internal/engine"
	"github.com/fred-bonn/nuz/internal/parser"
	"github.com/fred-bonn/nuz/internal/pokeapi"
)

type config struct {
	client pokeapi.Client
}

func (cfg *config) validateInput(trainerContent string) ([]*engine.Pokemon, error) {
	trainerPokemon, err := parser.ParseShowdown(trainerContent)
	if err != nil {
		return nil, fmt.Errorf("failed parsing showdown content: %w", err)
	}

	return cfg.validateParty(trainerPokemon)
}

func (cfg *config) validateParty(trainerPokemon []parser.ParsedPokemon) ([]*engine.Pokemon, error) {
	if len(trainerPokemon) == 0 {
		return nil, fmt.Errorf("showdown party is empty")
	}

	trainerParty, err := cfg.loadShowdown(trainerPokemon)
	if err != nil {
		return nil, fmt.Errorf("failed loading showdown file: %w", err)
	}

	return trainerParty, nil
}

func (cfg *config) loadShowdown(parsedPokemons []parser.ParsedPokemon) ([]*engine.Pokemon, error) {
	var res []*engine.Pokemon

	for _, parsedPokemon := range parsedPokemons {
		var moves []*engine.Move

		basePokemon, err := cfg.loadPokemon(apiName(parsedPokemon.Name))
		if err != nil {
			return nil, err
		}

		basePokemon.Name = cleanName(parsedPokemon.Name)

		for _, moveName := range parsedPokemon.Moves {
			baseMove, err := cfg.loadMove(apiName(moveName))
			if err != nil {
				return nil, err
			}

			if mb, ok := engine.MoveBalanceMap[baseMove.Name]; ok {
				mb.Apply(&baseMove)
			}

			baseMove.Name = cleanName(moveName)

			moves = append(moves, &baseMove)
		}

		finalPokemon, err := engine.InitPokemon(basePokemon, moves, parsedPokemon)
		if err != nil {
			return nil, err
		}

		res = append(res, &finalPokemon)
	}

	return res, nil
}

func (cfg *config) loadPokemon(name string) (engine.BasePokemon, error) {
	var p engine.BasePokemon

	data, err := os.ReadFile(fmt.Sprintf("data/pokemon/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Pokemon struct
		err = json.Unmarshal(data, &p)
		if err != nil {
			return engine.BasePokemon{}, fmt.Errorf("failed unmarshaling '%s' Pokemon data: %w", name, err)
		}

		return p, nil
	}

	// Otherwise, fetch the Pokemon data from the API
	pokemonJSON, err := cfg.client.FetchPokemon(name)
	if err != nil {
		return engine.BasePokemon{}, fmt.Errorf("failed fetching Pokemon '%s': %w", name, err)
	}
	fmt.Printf("Fetched '%s' from API\n", name)

	p, err = engine.ToPokemon(pokemonJSON)
	if err != nil {
		return engine.BasePokemon{}, err
	}

	// Save the fetched Pokemon data to a file for future use
	data, err = json.Marshal(p)
	if err != nil {
		return engine.BasePokemon{}, fmt.Errorf("failed marshaling Pokemon JSON data '%s' to file: %w", name, err)
	}
	writeToFile(fmt.Sprintf("data/pokemon/%s.json", name), data)

	return p, nil
}

func (cfg *config) loadMove(name string) (engine.Move, error) {
	var m engine.Move

	if strings.HasPrefix(name, "hidden-power") {
		// If the move is Hidden Power, generate it
		var err error
		m, err = generateHiddenPower(name)
		if err != nil {
			return engine.Move{}, err
		}
		return m, nil
	}

	data, err := os.ReadFile(fmt.Sprintf("data/moves/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Move struct
		err = json.Unmarshal(data, &m)
		if err != nil {
			return engine.Move{}, fmt.Errorf("failed unmarshaling Move '%s' data: %w", name, err)
		}

		return m, nil
	}

	// Otherwise, fetch the Move data from the API
	moveJson, err := cfg.client.FetchMove(name)
	if err != nil {
		return engine.Move{}, fmt.Errorf("failed fetching Move '%s': %w", name, err)
	}
	fmt.Printf("Fetched '%s' from API\n", name)

	m, err = engine.ToMove(moveJson)
	if err != nil {
		return engine.Move{}, err
	}

	// Save the fetched Move data using the internal Move struct to a file for future use
	data, err = json.Marshal(m)
	if err != nil {
		return m, fmt.Errorf("failed marshaling Move JSON data '%s' to file: %w", name, err)
	}
	writeToFile(fmt.Sprintf("data/moves/%s.json", name), data)

	return m, nil
}

func generateHiddenPower(name string) (engine.Move, error) {
	parts := strings.Split(name, "-")
	if len(parts) != 3 {
		return engine.Move{}, fmt.Errorf("type not specified for hidden power")
	}

	moveType := engine.StringToPokemonType(parts[2])
	if moveType == engine.NoType {
		return engine.Move{}, fmt.Errorf("%s is not a valid type for %s", parts[2], name)
	}

	move := engine.Move{
		Name:     "hidden power",
		Type:     moveType,
		Power:    60,
		Accuracy: 100,
		Class:    engine.SpecialClass,
	}

	return move, nil
}

func writeToFile(filename string, data []byte) error {
	dir := filepath.Dir(filename)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}
	return os.WriteFile(filename, data, 0644)
}

func apiName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, ".", "")
	name = strings.ReplaceAll(name, "’", "")
	return name
}

func cleanName(name string) string {
	name = strings.ToLower(name)
	if !hasHyphen(name) && !isRegionalPokemon(name) {
		name = strings.ReplaceAll(name, "-", " ")
	}

	return name
}

func hasHyphen(name string) bool {
	var withHyphen = map[string]struct{}{
		"ho-oh":     {},
		"porygon-z": {},
		"jangmo-o":  {},
		"hakamo-o":  {},
		"kommo-o":   {},
		"ting-lu":   {},
		"chien-pao": {},
		"wo-chien":  {},
		"chi-yu":    {},
	}

	if _, ok := withHyphen[name]; ok {
		return true
	}

	return false
}

func isRegionalPokemon(name string) bool {
	regions := []string{
		"-alola",
		"-galar",
		"-hisui",
		"-paldea",
	}

	name = strings.ToLower(name)

	for _, region := range regions {
		if strings.HasSuffix(name, region) {
			return true
		}
	}

	return false
}
