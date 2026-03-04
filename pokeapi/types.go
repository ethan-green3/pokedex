package pokeapi

type LocationResponse struct {
	Results  []Location `json:"results"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
