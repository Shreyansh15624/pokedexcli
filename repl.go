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
	callback    func(c *config, area string) error
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

type locationAreaDetail struct {
	Name              string             `json:"name"`
	PokemonEncounters []pokemonEncounter `json:"pokemon_encounters"`
}

type pokemonEncounter struct {
	Pokemon pokemonInfo `json:"pokemon"`
}

type pokemonInfo struct {
	Name string `json:"name"`
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func commandExit(c *config, area string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, area string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for _, command := range c.commandsMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func commandMapf(c *config, area string) error {
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

func commandMapb(c *config, area string) error {
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

func localPokemonEncounter(c *config, area string) error {
	if area == "" {
		fmt.Println("Please, enter a location-area name to get results...")
		return nil
	}

	baseUrl := "https://pokeapi.co/api/v2/location-area/" + area + "/"

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

	var pokemon_encounters locationAreaDetail

	err := json.Unmarshal(cachedData, &pokemon_encounters)
	if err != nil {
		return err
	}

	fmt.Println("Exploring " + area + "...")
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemon_encounters.PokemonEncounters {
		fmt.Println("-", pokemon.Pokemon.Name)
	}

	return nil
}
