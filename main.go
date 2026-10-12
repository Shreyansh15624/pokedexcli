package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println(`Welcome to the Pokedex!
Usage

help: Displays a help message
exit: Exit the Pokedex
`)
	return nil
}

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
