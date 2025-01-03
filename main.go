package main

import (
	"log"
	"os"
	"github.com/per1Peteia/gator/internal/config"
)

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	appState := state{
		c: &config,
	}

	commands := commands{
		callback: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	command := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err := commands.run(&appState, command); err != nil {
		log.Fatal(err)
	}
	
}
