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

func (c *Client) FetchPokemon(name string) (BasePokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	res, err := http.Get(url)
	if err != nil {
		return BasePokemon{}, fmt.Errorf("error fetching Pokemon data from API: %w", err)
	}
	defer res.Body.Close()

	var pokemonJSON pokemonJSON
	err = json.NewDecoder(res.Body).Decode(&pokemonJSON)
	if err != nil {
		return BasePokemon{}, fmt.Errorf("error decoding JSON into PokemonJSON: %w", err)
	}

	pokemon := pokemonJSON.ToPokemon()

	return pokemon, nil
}

func (c *Client) FetchMove(name string) (BaseMove, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", name)
	res, err := http.Get(url)
	if err != nil {
		return BaseMove{}, fmt.Errorf("error fetching Move data from API: %w", err)
	}
	defer res.Body.Close()

	var moveJSON moveJSON
	err = json.NewDecoder(res.Body).Decode(&moveJSON)
	if err != nil {
		return BaseMove{}, fmt.Errorf("error decoding JSON into MoveJSON: %w", err)
	}

	move := moveJSON.ToMove()

	return move, nil
}
