package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ethan-green3/pokedexcli/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: "", Pokedex: make(map[string]pokeapi.PokemonToCatch)}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		command := cleanInput(userInput)
		if len(command) == 0 {
			continue
		}
		if value, ok := commands[command[0]]; ok {
			err := value.callback(&config, command...)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println("Unknown command")
	}
}
