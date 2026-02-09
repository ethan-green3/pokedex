package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		command := cleanInput(userInput)
		if len(command) == 0 {
			continue
		}
		if value, ok := commands[command[0]]; ok {
			err := value.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else if !ok {
			fmt.Println("Unknown command")
		}
	}
}
