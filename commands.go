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
	callback    func(c *config, args ...string) error
}

type config struct {
	Next     string
	Previous string
	Pokedex  map[string]pokeapi.PokemonToCatch
}

func init() {
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex\n",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message\n",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas in the Pokemon world\n",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of the previous 20 location areas in the Pokemon world\n",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location to learn more about the Pokemon that are there\n\tExample usage: explore mt-coronet-2f\n",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon\n\tExample usage: catch squirtle\n",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon that you have caught to learn moreabout it\n\tExample usage: inspect pikachu\n",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List out the contents of your pokedex\n",
			callback:    commandPokedex,
		},
	}
}

func commandExit(c *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args ...string) error {
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

func commandMap(c *config, args ...string) error {
	res, err := pokeapi.GetLocationAreas(c.Next)
	if err != nil {
		return err
	}
	for _, item := range res.Results {
		fmt.Printf("%s\n", item.Name)
	}
	c.Next = *res.Next
	if isPreviousNil(res.Previous) {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
		return nil
	}

	c.Previous = *res.Previous
	return nil
}

func commandMapb(c *config, args ...string) error {
	res, err := pokeapi.GetLocationAreas(c.Previous)
	if err != nil {
		return err
	}
	for _, item := range res.Results {
		fmt.Printf("%s\n", item.Name)
	}
	c.Next = *res.Next
	if isPreviousNil(res.Previous) {
		c.Previous = "https://pokeapi.co/api/v2/location-area/"
		return nil
	}
	c.Previous = *res.Previous
	return nil
}

func commandExplore(c *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + args[1]
	exploreResponse, err := pokeapi.ExploreLocation(url)
	fmt.Println("Exploring:", args[1])
	if err != nil {
		return fmt.Errorf("Error within exploring a location: %w", err)
	}
	fmt.Println("Found Pokemon:")
	for _, p := range exploreResponse.Encounters {
		fmt.Printf("- %s\n", p.Pokemon.Name)
	}
	return nil
}

func commandCatch(c *config, args ...string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + args[1]
	res, err := pokeapi.CatchPokemon(url)
	if err != nil {
		return fmt.Errorf("Error within CatchPokemon call: %w", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", args[1])
	catch, err := TryCatch(res)
	if catch {
		fmt.Println(args[1], "was caught!")
		c.Pokedex[args[1]] = res
		if err := savePokedex(c.Pokedex); err != nil {
			return fmt.Errorf("Error saving to Pokedex: %w", err)
		}
		fmt.Println("You may now inspect it with the inspect command")
	} else {
		fmt.Println(args[1], "escaped!")
	}
	return nil
}

func commandInspect(c *config, args ...string) error {
	pokemon, ok := c.Pokedex[args[1]]
	if ok {
		fmt.Println("Name:", pokemon.Name)
		fmt.Println("Weight", pokemon.Weight)
		fmt.Println("Stats:")
		for _, val := range pokemon.Stats {
			fmt.Printf("-%s: %d\n", val.Stat.Name, val.BaseStat)
		}
		fmt.Println("Types:")
		for _, val := range pokemon.Types {
			fmt.Printf("-%s\n", val.Type.Name)
		}
		return nil
	}
	return fmt.Errorf("That Pokemon is not in your Pokedex, you need to catch it first!")
}

func commandPokedex(c *config, args ...string) error {
	fmt.Println("Your pokedex")
	for _, val := range c.Pokedex {
		fmt.Printf("- %s\n", val.Name)
	}
	return nil
}
