package main

import (
	"context"
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
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	commands.register("following", middlewareLoggedIn(handlerListFollowedFeeds))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollowFeed))

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

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUserDB, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, currentUserDB)
	}
}
