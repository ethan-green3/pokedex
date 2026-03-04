package main

import (
	"math/rand"

	"github.com/ethan-green3/pokedexcli/pokeapi"
)

func isPreviousNil(prev *string) bool {
	if prev == nil {
		return true
	}
	return false
}

func TryCatch(pokemon pokeapi.PokemonToCatch) (bool, error) {
	var chanceToFail float64
	exp := pokemon.BaseExperience
	if exp < 100 {
		chanceToFail = 0.5
	} else if exp > 100 && exp < 200 {
		chanceToFail = 0.6
	} else if exp > 200 && exp < 300 {
		chanceToFail = 0.7
	} else if exp > 300 && exp < 350 {
		chanceToFail = 0.8
	} else {
		chanceToFail = 0.9
	}

	roll := rand.Float64()
	if roll > chanceToFail {
		return true, nil
	} else {
		return false, nil
	}

}
