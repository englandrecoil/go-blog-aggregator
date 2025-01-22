package main

import (
	"log"
	"os"

	"github.com/englandrecoil/go-blog-aggregator/internal/config"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	programState := state{
		cfg: &cfg,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("warning: not enough arguments provided")
	}

	cmd := command{
		Name:      args[1],
		Arguments: args[2:],
	}

	if err := commands.run(&programState, cmd); err != nil {
		log.Fatal(err)
	}

}
