package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Booting up...\nInitialize...\nLoading Database...\nWelcome, Pokedex CLI!")
	fmt.Println()

	//initialize commandline reader
	scanner := bufio.NewScanner(os.Stdin)

	//initialize config and commands
	//config stores temp data for runtime
	cfg := &config{
		next:     "https://pokeapi.co/api/v2/location-area/",
		previous: "",
	}
	commands := getCommands(cfg)

	//loop to read commandline
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		if err := scanner.Err(); err != nil {
			fmt.Fprint(os.Stderr, "reading standard input:", err)
		}
		input := scanner.Text()
		_, ok := commands[input]
		if ok {
			commands[input].callback(cfg)
		} else {
			fmt.Printf("Input not valid: %s\n", input)
		}
	}
}
