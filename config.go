package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fred-bonn/nuzlocke-verifier/internal/parser"
	"github.com/fred-bonn/nuzlocke-verifier/internal/pokeapi"
)

type config struct {
	client pokeapi.Client
}

func (cfg *config) validateInput(trainerPath string) ([]*pokemon, error) {
	trainerFullPath, err := filepath.Abs(trainerPath)
	if err != nil {
		return nil, fmt.Errorf("failed getting absolute path: %w", err)
	}

	trainerPokemon, err := parser.ReadShowdownFile(trainerFullPath)
	if err != nil {
		return nil, fmt.Errorf("failed reading showdown file: %w", err)
	}
	if len(trainerPokemon) == 0 {
		return nil, err
	}

	trainerParty, err := cfg.loadShowdown(trainerPokemon)
	if err != nil {
		return nil, fmt.Errorf("failed loading showdown file: %w", err)
	}

	return trainerParty, nil
}

func (cfg *config) loadShowdown(mons []parser.ParsedPokemon) ([]*pokemon, error) {
	var res []*pokemon

	for _, mon := range mons {
		var moves []*Move

		basePokemon, err := cfg.loadPokemon(apiName(mon.Name))
		if err != nil {
			return nil, err
		}
		basePokemon.Name = cleanName(mon.Name, false)

		for _, moveName := range mon.Moves {
			baseMove, err := cfg.loadMove(apiName(moveName))
			if err != nil {
				return nil, err
			}
			baseMove.Name = cleanName(moveName, true)

			moves = append(moves, &baseMove)
		}

		finalPokemon, err := initPokemon(basePokemon, mon.Level, mon.IVs, mon.Nature, moves, mon.HP, stringToAilmentState(mon.Status))
		if err != nil {
			return nil, err
		}

		item, err := registerItem(stringToItemState(strings.ToLower(mon.Item)), &finalPokemon)
		if err != nil {
			return nil, err
		}
		finalPokemon.item = item

		finalPokemon.ability = stringToAbility(strings.ToLower(mon.Ability))

		res = append(res, &finalPokemon)
	}

	return res, nil
}

func (cfg *config) loadPokemon(name string) (BasePokemon, error) {
	var p BasePokemon

	data, err := os.ReadFile(fmt.Sprintf("data/pokemon/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Pokemon struct
		err = json.Unmarshal(data, &p)
		if err != nil {
			return BasePokemon{}, fmt.Errorf("failed unmarshaling '%s' Pokemon data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
	} else {
		// Otherwise, fetch the Pokemon data from the API
		pokemonJSON, err := cfg.client.FetchPokemon(name)
		if err != nil {
			return BasePokemon{}, fmt.Errorf("failed fetching Pokemon '%s': %w", name, err)
		}
		fmt.Printf("Fetched '%s' from API\n", name)

		p = toPokemon(pokemonJSON)

		// Save the fetched Pokemon data to a file for future use
		vlogln(p)
		data, err = json.Marshal(p)
		vlogln(data)
		if err != nil {
			return BasePokemon{}, fmt.Errorf("failed marshaling Pokemon JSON data '%s' to file: %w", name, err)
		}
		writeToFile(fmt.Sprintf("data/pokemon/%s.json", name), data)
	}

	return p, nil
}

func (cfg *config) loadMove(name string) (Move, error) {
	var m Move

	if strings.HasPrefix(name, "hidden-power") {
		// If the move is Hidden Power, generate it
		var err error
		m, err = generateHiddenPower(name)
		if err != nil {
			return Move{}, err
		}
		return m, nil
	}

	data, err := os.ReadFile(fmt.Sprintf("data/moves/%s.json", name))
	if err == nil {
		// If the file exists and is read successfully, unmarshal it into a Move struct
		err = json.Unmarshal(data, &m)
		if err != nil {
			return Move{}, fmt.Errorf("failed unmarshaling Move '%s' data: %w", name, err)
		}
		fmt.Printf("Loaded '%s' from file\n", name)
	} else {
		// Otherwise, fetch the Move data from the API
		moveJson, err := cfg.client.FetchMove(name)
		if err != nil {
			return Move{}, fmt.Errorf("failed fetching Move '%s': %w", name, err)
		}
		fmt.Printf("Fetched '%s' from API\n", name)

		m = toMove(moveJson)

		// Save the fetched Move data using the internal Move struct to a file for future use
		data, err = json.Marshal(m)
		if err != nil {
			return m, fmt.Errorf("failed marshaling Move JSON data '%s' to file: %w", name, err)
		}
		writeToFile(fmt.Sprintf("data/moves/%s.json", name), data)
	}

	if mb, ok := moveBalanceMap[m.Name]; ok {
		mb.apply(&m)
	}

	return m, nil
}

func generateHiddenPower(name string) (Move, error) {
	parts := strings.Split(name, "-")
	if len(parts) != 3 {
		return Move{}, fmt.Errorf("type not specified for hidden power")
	}

	move := Move{
		Name:     "hidden power",
		Type:     stringToPokemonType(parts[2]),
		Power:    60,
		Accuracy: 100,
		Class:    specialClass,
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

func cleanName(name string, move bool) string {
	name = strings.ToLower(name)
	if !hasHyphen(name) && !isRegionalPokemon(name) {
		name = strings.ReplaceAll(name, "-", " ")
	}

	return name
}

func hasHyphen(name string) bool {
	var pokemonWithHyphen = map[string]struct{}{
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

	if _, ok := pokemonWithHyphen[name]; ok {
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
