package clicommand

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/christopherhanke/bootdev_pokedex/internal/pokecache"
	"github.com/christopherhanke/bootdev_pokedex/internal/pokedex"
)

type cliCommand struct {
	Name        string
	Description string
	Callback    func(*Config, ...string) error
}

type Config struct {
	Next     string
	Previous string
	Cache    pokecache.Cache
	Pokedex  pokedex.Pokedex
}

func GetCommands(cfg *Config) map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message.",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex.",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays next 20 location areas in the Pokemon world.",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays last 20 location areas in the Pokemon world.",
			Callback:    commandMapB,
		},
		"explore": {
			Name:        "explore",
			Description: "Explore a given area and list all Pokemon there.",
			Callback:    commandExplore,
		},
		"catch": {
			Name:        "catch",
			Description: "catch a given Pokemon.",
			Callback:    commandCatch,
		},
	}
	return commands
}

func commandHelp(cfg *Config, args ...string) error {
	fmt.Print("\nUsage of the Pokedex\nList of commands:\n\n")
	commands := GetCommands(cfg)
	for key, value := range commands {
		fmt.Printf("%s: %s\n", key, value.Description)
	}
	fmt.Println()
	return nil
}

func commandExit(cfg *Config, args ...string) error {
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

func commandMap(cfg *Config, args ...string) error {
	//exit function if config next is not set
	if cfg.Next == "" {
		fmt.Println("Error cfg has no next")
		return fmt.Errorf("cfg.next undefined")
	}

	var locations getLocations

	//check data in Cache and present if avaible
	if val, ok := cfg.Cache.Get(cfg.Next); ok {
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return err
		}
	} else {
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

		err = json.Unmarshal(data, &locations)
		if err != nil {
			fmt.Println("Error unmarshal", err)
			return err
		}

		//add data to Cache
		if val, err := json.Marshal(locations); err == nil {
			cfg.Cache.Add(cfg.Next, val)
		}
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

func commandMapB(cfg *Config, args ...string) error {
	//exit function if config previous is not set
	if cfg.Previous == "" {
		fmt.Println("Error cfg has no previous")
		return fmt.Errorf("cfg.previous undefined")
	}

	var locations getLocations

	if val, ok := cfg.Cache.Get(cfg.Previous); ok {
		err := json.Unmarshal(val, &locations)
		if err != nil {
			return err
		}
	} else {
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

		err = json.Unmarshal(data, &locations)
		if err != nil {
			fmt.Println("Error unmarshal", err)
			return err
		}

		//add data to Cache
		if val, err := json.Marshal(locations); err == nil {
			cfg.Cache.Add(cfg.Previous, val)
		}
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

type getEncounters struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func commandExplore(cfg *Config, args ...string) error {
	// exit when no area (args) is given
	if len(args) == 0 {
		fmt.Printf("No Area given to explore.\nPlease use: explore <area>\n\n")
		return fmt.Errorf("no area arg given")
	}

	fmt.Printf("Explore was called with arg: %v\n", args[0])

	urlArea := "https://pokeapi.co/api/v2/location-area/" + args[0]
	var encounters getEncounters

	//check if data is in cache already
	if val, ok := cfg.Cache.Get(urlArea); ok {
		err := json.Unmarshal(val, &encounters)
		if err != nil {
			return err
		}
	} else {
		//get data from PokeApi
		resp, err := http.Get(urlArea)
		if err != nil {
			fmt.Println("error get", err)
			return err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return err
		}
		err = json.Unmarshal(data, &encounters)
		if err != nil {
			fmt.Println("error unmarshal", err)
			return err
		}

		//add data to Cache
		if val, err := json.Marshal(encounters); err == nil {
			cfg.Cache.Add(urlArea, val)
		}
	}

	//print list of all Pokemon in area
	for _, pokemon := range encounters.PokemonEncounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}

	return nil
}

func commandCatch(cfg *Config, args ...string) error {
	//exit when no Pokemon is given
	if len(args) == 0 {
		fmt.Printf("No Pokemon given to catch.\nPlease use: catch <pokemon>\n")
		return fmt.Errorf("no Pokemon given")
	}

	urlPokemon := "https://pokeapi.co/api/v2/pokemon/" + args[0]
	var pokemon pokedex.Pokemon

	//check if data is in cache
	if val, ok := cfg.Cache.Get(urlPokemon); ok {
		err := json.Unmarshal(val, &pokemon)
		if err != nil {
			return err
		}
	} else {
		//get data from PokeApi
		resp, err := http.Get(urlPokemon)
		if err != nil {
			fmt.Println("error get", err)
			return err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("error reading body", err)
			return err
		}
		err = json.Unmarshal(data, &pokemon)
		if err != nil {
			fmt.Println("error unmarshal", err)
			return err
		}

		//add to cache
		if val, err := json.Marshal(pokemon); err != nil {
			cfg.Cache.Add(urlPokemon, val)
		}
	}

	//try to catch pokemon
	fmt.Printf("Throwing a pokeball at %v\n", args[0])
	chance := rand.Intn(300)
	if chance >= pokemon.BaseExperience {
		fmt.Printf("%v was caught!\n", args[0])
		cfg.Pokedex.Add(pokemon)
	} else {
		fmt.Printf("%v escaped!\n", args[0])
	}

	return nil
}
