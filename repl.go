package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"pokedexcli/internal/pokecache"
	"strings"
)

type config struct {
	cache       *pokecache.Cache
	commandsMap map[string]cliCommand
	pokemon     map[string]Pokemon
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

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
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

	cacheData, ok := c.cache.Get(baseUrl)
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
		cacheData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		c.cache.Add(baseUrl, cacheData)
	}

	var results returnedJSON
	err := json.Unmarshal(cacheData, &results)
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

	cacheData, ok := c.cache.Get(baseUrl)
	if !ok {
		response, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		cacheData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode > 299 {
			err := fmt.Sprintf("Response failed with status code: %d and\nbody: %v\n", response.StatusCode, cacheData)
			return errors.New(err)
		}

		c.cache.Add(baseUrl, cacheData)
	}

	var results returnedJSON

	err := json.Unmarshal(cacheData, &results)
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

	cacheData, ok := c.cache.Get(baseUrl)
	if !ok {
		response, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		cacheData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode > 299 {
			err := fmt.Sprintf("Response failed with status code: %d and\nbody: %v\n", response.StatusCode, cacheData)
			return errors.New(err)
		}

		c.cache.Add(baseUrl, cacheData)
	}

	var pokemon_encounters locationAreaDetail

	err := json.Unmarshal(cacheData, &pokemon_encounters)
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

func catchThatPokemon(c *config, pokemon string) error {
	if pokemon == "" {
		fmt.Println("Please, specify the pokemon you want to catch...")
		return nil
	}

	baseUrl := "https://pokeapi.co/api/v2/pokemon/" + pokemon + "/"
	cacheData, ok := c.cache.Get(pokemon)
	if !ok {
		response, err := http.Get(baseUrl)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		cacheData, err = io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		if response.StatusCode > 299 {
			err := fmt.Sprintf("Response failed with status code: %d and\nbody: %v\n", response.StatusCode, cacheData)
			return errors.New(err)
		}
		c.cache.Add(pokemon, cacheData)
	}

	var targetPokemon Pokemon

	err := json.Unmarshal(cacheData, &targetPokemon)
	if err != nil {
		return err
	}

	fmt.Println("Throwing a Pokeball at " + pokemon + "...")

	catchPossibility := rand.IntN(200)
	if targetPokemon.BaseExperience > catchPossibility {
		c.pokemon[pokemon] = targetPokemon
		fmt.Println(pokemon + " was caught!")
	} else {
		fmt.Println(pokemon + " escaped!")
	}

	return nil
}

func inspectPokemon(c *config, pokemon string) error {
	if pokemon == "" {
		fmt.Println("specifiy the pokemon with name you want to inspect")
		return nil
	}

	pokemonDetails, ok := c.pokemon[pokemon]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}

	fmt.Printf("ID: %d\n", pokemonDetails.ID)
	fmt.Printf("Name: %s\n", pokemonDetails.Name)
	fmt.Printf("Base Experience: %d\n", pokemonDetails.BaseExperience)
	fmt.Printf("Height: %d\n", pokemonDetails.Height)
	fmt.Printf("Weight: %d\n", pokemonDetails.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonDetails.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, typ := range pokemonDetails.Types {
		fmt.Printf("  - %s\n", typ.Type.Name)
	}

	return nil
}

func pokedex(c *config, pokemon string) error {
	fmt.Println("Your Pokedex:")
	for name := range c.pokemon {
		fmt.Println(" - " + name)
	}
	return nil
}
