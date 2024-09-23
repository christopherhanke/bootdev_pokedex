package clicommand

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

func GetCommands(cfg *Config) map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays next 20 location areas in the Pokemon world",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays last 20 location areas in the Pokemon world",
			Callback:    commandMapB,
		},
	}
	return commands
}

func commandHelp(cfg *Config) error {
	fmt.Print("\nUsage of the Pokedex\nList of commands:\n\n")
	commands := GetCommands(cfg)
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.Description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *Config) error {
	os.Exit(0)
	return nil
}

type getLocations struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(cfg *Config) error {
	//exit function if config next is not set
	if cfg.Next == "" {
		fmt.Println("Error cfg has no next")
		return fmt.Errorf("cfg.next undefined")
	}

	//get data from PokeApi, read and work
	resp, err := http.Get(cfg.Next)
	if err != nil {
		fmt.Println("Error get", err)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error read", err)
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

	//update Next and Previous in config
	if locations.Next != "" {
		cfg.Next = locations.Next
	}
	if locations.Previous != "" {
		cfg.Previous = locations.Previous
	}
	return nil
}

func commandMapB(cfg *Config) error {
	//exit function if config previous is not set
	if cfg.Previous == "" {
		fmt.Println("Error cfg has no previous")
		return fmt.Errorf("cfg.previous undefined")
	}

	resp, err := http.Get(cfg.Previous)
	if err != nil {
		fmt.Println("Error get", err)
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error read", err)
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

	//update Next and Previous in config
	if locations.Next != "" {
		cfg.Next = locations.Next
	}
	if locations.Previous != "" {
		cfg.Previous = locations.Previous
	}

	return nil
}
