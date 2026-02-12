package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/ethan-green3/pokedexcli/pokeapi"
)

var commands map[string]cliCommand

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	Next     string
	Previous string
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, name := range names {
		cmd := commands[name]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandMap(c *config) error {
	res, err := pokeapi.GetLocationAreas(c.Next)
	if err != nil {
		return err
	}
	printResponse(res)

	c.Next = *res.Next
	if isPreviousNil(res.Previous) {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
		return nil
	}

	c.Previous = *res.Previous
	return nil
}

func commandMapb(c *config) error {
	res, err := pokeapi.GetLocationAreas(c.Previous)
	if err != nil {
		return err
	}
	printResponse(res)
	c.Next = *res.Next
	if isPreviousNil(res.Previous) {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
		return nil
	}
	c.Previous = *res.Previous
	return nil
}

func isPreviousNil(prev *string) bool {
	if prev == nil {
		return true
	}
	return false
}

func printResponse(res pokeapi.Response) {
	for _, area := range res.Results {
		fmt.Println(area.Name)
	}
}

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
	}
}
