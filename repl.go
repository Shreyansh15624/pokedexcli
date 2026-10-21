package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/internal/pokecache"
	"strings"
)

type config struct {
	cache       *pokecache.Cache
	commandsMap map[string]cliCommand
	next        *string
	prev        *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type returnedJSON struct {
	Count   int
	Next    *string
	Prev    *string `json:"previous"`
	Results []resultJSON
}

type resultJSON struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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

func commandMapf(c *config) error {
	var baseUrl string
	if c.next == nil {
		baseUrl = "https://pokeapi.co/api/v2/location-area/"
	} else {
		baseUrl = *c.next
	}

	cachedData, ok := c.cache.Get(baseUrl)
	if !ok {
		response, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		if response.StatusCode > 299 {
			err := fmt.Sprintf("Response failed with status code: %d and\nbody: %v\n", response.StatusCode, response.Body)
			return errors.New(err)
		}
		cachedData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		c.cache.Add(baseUrl, cachedData)
	}

	var results returnedJSON
	err := json.Unmarshal(cachedData, &results)
	if err != nil {
		return err
	}
	c.next = results.Next
	c.prev = results.Prev

	for _, location := range results.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(c *config) error {
	var baseUrl string
	if c.prev == nil {
		fmt.Println("you're on the first page")
		return nil
		//baseUrl = "https://pokeapi.co/api/v2/location-area/"
	} else {
		baseUrl = *c.prev
	}

	cachedData, ok := c.cache.Get(baseUrl)
	if !ok {
		response, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		cachedData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode > 299 {
			err := fmt.Sprintf("Response failed with status code: %d and\nbody: %v\n", response.StatusCode, cachedData)
			return errors.New(err)
		}

		c.cache.Add(baseUrl, cachedData)
	}

	var results returnedJSON

	err := json.Unmarshal(cachedData, &results)
	if err != nil {
		return err
	}
	c.next = results.Next
	c.prev = results.Prev

	for _, location := range results.Results {
		fmt.Println(location.Name)
	}
	return nil
}
