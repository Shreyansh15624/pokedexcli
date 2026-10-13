package main

import (
	"fmt"
	"os"
	"strings"
)

type config struct {
	commandsMap map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range c.commandsMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}
