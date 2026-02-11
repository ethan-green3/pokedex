package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := config{Next: "https://pokeapi.co/api/v2/location-area/", Previous: ""}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		command := cleanInput(userInput)
		if len(command) == 0 {
			continue
		}
		if value, ok := commands[command[0]]; ok {
			err := value.callback(&config)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		fmt.Println("Unknown command")
	}
}
