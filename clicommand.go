package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

func getCommands(cfg *config) map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays last 20 location areas in the Pokemon world",
			callback:    commandMapB,
		},
	}
	return commands
}

func commandHelp(cfg *config) error {
	fmt.Print("\nUsage of the Pokedex\nList of commands:\n\n")
	commands := getCommands(cfg)
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *config) error {
	os.Exit(0)
	return nil
}

type getLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(cfg *config) error {
	fmt.Println("command map")
	resp, err := http.Get(cfg.next)
	if err != nil {
		fmt.Println("Error get")
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error read")
		return err
	}

	var locations getLocations
	err = json.Unmarshal(data, &locations)
	if err != nil {
		fmt.Println("Error unmarshal", err)
		return err
	}
	for _, val := range locations.Results {
		fmt.Println(val.Name)
	}
	return nil
}

func commandMapB(cfg *config) error {
	fmt.Println("command mapb")
	return nil
}
