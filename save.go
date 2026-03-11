package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/ethan-green3/pokedexcli/pokeapi"
)

const pokedexFile = "pokedex.json"

func savePokedex(pokedex map[string]pokeapi.PokemonToCatch) error {
	data, err := json.MarshalIndent(pokedex, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pokedex: %w", err)
	}

	if err := os.WriteFile(pokedexFile, data, 0644); err != nil {
		return fmt.Errorf("write pokedex file: %w", err)
	}

	return nil
}

func loadPokedex() (map[string]pokeapi.PokemonToCatch, error) {
	data, err := os.ReadFile(pokedexFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return make(map[string]pokeapi.PokemonToCatch), nil
		}
		return nil, fmt.Errorf("read pokedex file: %w", err)
	}

	var pokedex map[string]pokeapi.PokemonToCatch
	if err := json.Unmarshal(data, &pokedex); err != nil {
		return nil, fmt.Errorf("unmarshal pokedex: %w", err)
	}

	if pokedex == nil {
		pokedex = make(map[string]pokeapi.PokemonToCatch)
	}

	return pokedex, nil
}
