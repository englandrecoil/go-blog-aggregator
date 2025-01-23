package main

import (
	"context"
	"fmt"
	"time"

	"github.com/englandrecoil/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
	}

	currentUserDB, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	time := time.Now()
	args := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time,
		UpdatedAt: time,
		Name:      cmd.Arguments[0],
		Url:       cmd.Arguments[1],
		UserID:    currentUserDB.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), args)
	if err != nil {
		return fmt.Errorf("couldn't create new feed: %w", err)
	}

	fmt.Println("New feed created successfully:")
	fmt.Println("=====================================")
	printFeed(feed, currentUserDB)
	fmt.Println("=====================================")

	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found")
		return nil
	}

	fmt.Printf("Found %d feeds:\n", len(feeds))
	fmt.Println("=====================================")

	for _, value := range feeds {
		user, err := s.db.GetUserByID(context.Background(), value.UserID)
		if err != nil {
			return err
		}

		printFeed(value, user)
		fmt.Println("=====================================")
	}

	return nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* User:          %s\n", user.Name)
}
