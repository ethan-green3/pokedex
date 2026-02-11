package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func GetLocationAreas(url string) (Response, error) {
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

	err = json.Unmarshal(data, &r)
	if err != nil {
		return r, fmt.Errorf("Error unmarshaling JSON into response struct: %w", err)
	}
	for _, item := range r.Results {
		fmt.Printf("%s\n", item.Name)
	}
	return r, nil
}
