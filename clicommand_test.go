package main

import (
	"testing"
)

func TestGetCommands(t *testing.T) {
	cfg := config{}
	commands := getCommands(&cfg)
	_, ok := commands["help"]
	if !ok {
		t.Fatalf("list of commands, doesnt support 'help'")
	}
}
