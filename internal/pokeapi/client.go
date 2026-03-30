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

func (c *Client) FetchPokemon(name string) (Pokemon, error) {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error fetching Pokemon data: %w", err)
	}
	defer res.Body.Close()

	var pokemonJSON PokemonJSON
	err = json.NewDecoder(res.Body).Decode(&pokemonJSON)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error decoding JSON into PokemonJSON: %w", err)
	}

	pokemon := pokemonJSON.ToPokemon()

	return pokemon, nil
}
