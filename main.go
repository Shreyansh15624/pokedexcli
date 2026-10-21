package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	const interval = 10 * time.Second
	config := config{
		cache: pokecache.NewCache(interval),
		commandsMap: map[string]cliCommand{
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
				description: "Displays each 20 next location areas in the pokemon world for each subsequent call",
				callback:    commandMapf,
			},
			"mapb": {
				name:        "mapb",
				description: "Dispalys each 20 previous location areas in the pokemon world for each subsequent call",
				callback:    commandMapb,
			},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to PokeDex!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInputString := scanner.Text()
		cleanedString := cleanInput(userInputString)
		if len(cleanedString) == 0 {
			continue
		}
		userInput := cleanedString[0]
		command, exists := config.commandsMap[userInput]
		if exists {
			err := command.callback(&config)
			if err != nil {
				fmt.Printf("%v", err)
			}
			// fmt.Println(out)
		} else {
			fmt.Println("unknown command")
		}
	}
}
