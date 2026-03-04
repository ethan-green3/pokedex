package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethan-green3/pokedexcli/internal/pokecache"
)

var cache = pokecache.NewCache(time.Second * 30)

type LocationResponse struct {
	Results  []Location `json:"results"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Returns a list of 20 locations based on the given URL by the map and mapb commands
func GetLocationAreas(url string) (LocationResponse, error) {
	if val, found := cache.Get(url); found {
		var r LocationResponse
		err := json.Unmarshal(val, &r)
		if err != nil {
			return r, fmt.Errorf("Error unmarshaling JSON from cache")
		}
		return r, nil
	}

	client := &http.Client{}
	var r LocationResponse
	res, err := client.Get(url)
	if err != nil {
		return r, fmt.Errorf("Error fetching location areas: %w", err)
	}

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return r, fmt.Errorf("Received 400 Error code from PokeAPI: %w", err)
	}

	if res.StatusCode >= 500 {
		return r, fmt.Errorf("Received Internal Server Error from PokeAPI: %w", err)
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return r, fmt.Errorf("Error reading response body: %w", err)
	}

	cache.Add(url, data)
	if err = json.Unmarshal(data, &r); err != nil {
		return r, fmt.Errorf("Error unmarshaling JSON into response struct: %w", err)
	}
	return r, nil
}

type ExploreResponse struct {
	Encounters []PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonEncounters struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// Explores a given location by the explore command to return back a response that has a slice of pokemon encounters which include
// a Name and URL for the pokemon in that encounter
func ExploreLocation(url string) (ExploreResponse, error) {
	if val, found := cache.Get(url); found {
		var expRes ExploreResponse
		err := json.Unmarshal(val, &expRes)
		if err != nil {
			return expRes, fmt.Errorf("Error unmarshaling explore response data from cache: %w", err)
		}
		return expRes, nil
	}

	client := &http.Client{}
	var expRes ExploreResponse
	res, err := client.Get(url)
	if err != nil {
		return expRes, fmt.Errorf("Error exploring area: %w", err)
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return expRes, fmt.Errorf("Error reading data from exlore response body: %w", err)
	}

	cache.Add(url, data)
	if err = json.Unmarshal(data, &expRes); err != nil {
		return expRes, fmt.Errorf("Error unmarshaling JSON into Explore resposne struct: %w", err)
	}
	return expRes, nil
}

type PokemonToCatch struct {
	BaseExperience int     `json:"base_experience"`
	Name           string  `json:"name"`
	ID             int     `json:"id"`
	Weight         int     `json:"weight"`
	Height         int     `json:"height"`
	Stats          []Stats `json:"stats"`
	Types          []Types `json:"types"`
}

type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Types struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"Stat"`
}

func CatchPokemon(url string) (PokemonToCatch, error) {
	if val, found := cache.Get(url); found {
		var p PokemonToCatch
		if err := json.Unmarshal(val, &p); err != nil {
			return p, fmt.Errorf("Error unmarshaling pokemon to catch from cache")
		}
		return p, nil
	}
	client := http.Client{}
	var p PokemonToCatch
	res, err := client.Get(url)
	if err != nil {
		return p, fmt.Errorf("Error getting data for pokemon to catch: %w", err)
	}
	data, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	cache.Add(url, data)
	if err = json.Unmarshal(data, &p); err != nil {
		return p, fmt.Errorf("Error unmarshaling API response into Pokemon to catch struct: %w", err)
	}

	return p, nil
}
