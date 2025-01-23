package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/englandrecoil/go-blog-aggregator/internal/config"
	"github.com/englandrecoil/go-blog-aggregator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", cfg.Url)
	if err != nil {
		log.Fatalf("couldn't create connection to db: %s", err)
	}
	defer db.Close()
	dbQueries := database.New(db)

	programState := state{
		cfg: &cfg,
		db:  dbQueries,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerListUsers)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Warning: not enough arguments provided. Usage: cli <command> [args...]")
	}

	cmd := command{
		Name:      args[1],
		Arguments: args[2:],
	}

	if err := commands.run(&programState, cmd); err != nil {
		log.Fatal(err)
	}

}
