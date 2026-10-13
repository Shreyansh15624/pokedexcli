package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cliCommand := map[string]cliCommand{
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
		"empty": {
			name:        "empty",
			description: "Singifies an empty input",
			callback:    commandHelp,
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Welcome to PokeDex!")
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInputString := scanner.Text()
		cleanedString := cleanInput(userInputString)
		if len(cleanedString) == 0 {
			continue
		}
		userInput := cleanedString[0]
		command, exists := cliCommand[userInput]
		if exists {
			err := command.callback()
			if err != nil {
				fmt.Println("%v", err)
			}
			// fmt.Println(out)
		} else {
			fmt.Println("unknown command")
		}
	}
}
