package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ethan-green3/pokedexcli/pokeapi"
)

func main() {
	pokedex, err := loadPokedex()
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: could not load pokedex:", err)
		pokedex = make(map[string]pokeapi.PokemonToCatch)
	}

	cfg := config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Pokedex:  pokedex,
	}

	p := tea.NewProgram(newModel(cfg), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error running pokedex:", err)
		os.Exit(1)
	}
}
