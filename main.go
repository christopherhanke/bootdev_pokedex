package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/christopherhanke/bootdev_pokedex/internal/clicommand"
	"github.com/christopherhanke/bootdev_pokedex/internal/pokecache"
)

const START string = "https://pokeapi.co/api/v2/location-area/"

func main() {
	fmt.Println("Booting up...\nInitialize...\nLoading Database...\nWelcome, Pokedex CLI!")
	fmt.Println()

	//initialize commandline reader
	scanner := bufio.NewScanner(os.Stdin)

	//initialize config and commands
	//config stores temp data for runtime
	cfg := &clicommand.Config{
		Next:     START,
		Previous: "",
		Cache:    pokecache.NewCache(time.Minute * 5),
	}
	commands := clicommand.GetCommands(cfg)

	//loop to read commandline
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, "reading standard input:", err)
		}
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			continue
		}
		commandName := input[0]
		var args []string
		if len(input) > 1 {
			args = input[1:]
		}
		_, ok := commands[commandName]
		if ok {
			commands[commandName].Callback(cfg, args...)
		} else {
			fmt.Printf("Input not valid: %s\n", commandName)
		}
	}
}

func cleanInput(input string) []string {
	loweredInput := strings.ToLower(input)
	output := strings.Fields(loweredInput)
	return output
}
