package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/englandrecoil/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	if _, err := s.db.GetUser(context.Background(), cmd.Arguments[0]); err != nil {
		return fmt.Errorf("user '%s' doesn't exist", cmd.Arguments[0])
	}

	if err := s.cfg.SetUser(cmd.Arguments[0]); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Printf("User %s has been set\n", s.cfg.CurrentUserName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	// check command arguments
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	// check if username already exists
	if user, err := s.db.GetUser(context.Background(), cmd.Arguments[0]); err == nil {
		return fmt.Errorf("user '%s' is already exists", user.Name)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("unexpected error while checking user: %w", err)
	}

	time := time.Now()
	args := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time,
		UpdatedAt: time,
		Name:      cmd.Arguments[0],
	}

	user, err := s.db.CreateUser(context.Background(), args)
	if err != nil {
		return fmt.Errorf("can't create new user: %w", err)
	}

	if err = s.cfg.SetUser(user.Name); err != nil {
		return fmt.Errorf("couldn't set new user: %w", err)
	}

	fmt.Println("User registered successfully")
	printUserInfo(user)

	return nil
}

func printUserInfo(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Created_at:    %v\n", user.CreatedAt)
	fmt.Printf(" * Updated_at:    %v\n", user.UpdatedAt)
	fmt.Printf(" * Name:    %v\n", user.Name)
}
