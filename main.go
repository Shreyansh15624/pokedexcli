package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	const interval = 20 * time.Second
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
			"explore": {
				name:        "explore",
				description: "Display the names of all Pokémon found in that area",
				callback:    localPokemonEncounter,
			},
			"catch": {
				name:        "catch",
				description: "Catching a specific pokemon found in the current local area",
				callback:    catchThatPokemon,
			},
		},
		pokemon: make(map[string]Pokemon),
	}

	scanner := bufio.NewScanner(os.Stdin)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr)
	}
	fmt.Println("Welcome to PokeDex!")
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInputString := scanner.Text()
		cleanedString := cleanInput(userInputString)
		if len(cleanedString) == 0 {
			continue
		}
		userInput1 := cleanedString[0]
		var userInput2 string
		if len(cleanedString) > 1 {
			userInput2 = cleanedString[1]
		} else {
			userInput2 = ""
		}
		command, exists := config.commandsMap[userInput1]
		if exists {
			err := command.callback(&config, userInput2)
			if err != nil {
				fmt.Printf("%v", err)
			}
		} else {
			fmt.Println("unknown command")
		}
	}
}
