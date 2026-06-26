package pokeapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient http.Client
}

func NewClient() Client {
	return Client{
		httpClient: http.Client{},
	}
}

func (c *Client) FetchPokemon(name string) (PokemonJSON, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	res, err := http.Get(url)
	if err != nil {
		return PokemonJSON{}, fmt.Errorf("error fetching Pokemon data from API: %w", err)
	}
	defer res.Body.Close()

	var pokemonJSON PokemonJSON
	err = json.NewDecoder(res.Body).Decode(&pokemonJSON)
	if err != nil {
		return PokemonJSON{}, fmt.Errorf("error decoding JSON into PokemonJSON: %w", err)
	}

	return pokemonJSON, nil
}

func (c *Client) FetchMove(name string) (MoveJSON, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", name)
	res, err := http.Get(url)
	if err != nil {
		return MoveJSON{}, fmt.Errorf("error fetching Move data from API: %w", err)
	}
	defer res.Body.Close()

	var move MoveJSON
	err = json.NewDecoder(res.Body).Decode(&move)
	if err != nil {
		return MoveJSON{}, fmt.Errorf("error decoding JSON into MoveJSON: %w", err)
	}

	return move, nil
}
