package clicommand

import (
	"testing"
	"time"

	"github.com/christopherhanke/bootdev_pokedex/internal/pokecache"
)

func TestGetCommands(t *testing.T) {
	cfg := &Config{}
	commands := GetCommands(cfg)
	_, ok := commands["help"]
	if !ok {
		t.Fatalf("list of commands, doesnt support 'help'")
	}
}

func TestCommandMapEmpty(t *testing.T) {
	cfg := &Config{}
	err := commandMap(cfg)
	if err == nil {
		t.Fatalf("Error in commandMap: cfg.next defined instead of nil")
	}
}

func TestCommandMapNext(t *testing.T) {
	cfg := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Cache:    pokecache.NewCache(time.Minute * 5),
	}
	err := commandMap(cfg)
	if err != nil {
		t.Fatalf("Error in commandMap: cfg.next undefiend instead of value")
	}
}
