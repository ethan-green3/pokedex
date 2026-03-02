package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ethan-green3/pokedexcli/internal/pokecache"
)

type Response struct {
	Results  []Location `json:"results"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

var cache = pokecache.NewCache(time.Second * 5)

func GetLocationAreas(url string) (Response, error) {
	if val, found := cache.Get(url); found {
		var r Response
		err := json.Unmarshal(val, &r)
		if err != nil {
			return r, fmt.Errorf("Error unmarshaling JSON from cache")
		}
		return r, nil
	}

	client := &http.Client{}
	var r Response
	res, err := client.Get(url)
	if err != nil {
		return r, fmt.Errorf("Error fetching location areas: %w", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return r, fmt.Errorf("Error reading response body: %w", err)
	}
	cache.Add(url, data)
	err = json.Unmarshal(data, &r)
	if err != nil {
		return r, fmt.Errorf("Error unmarshaling JSON into response struct: %w", err)
	}
	return r, nil
}
