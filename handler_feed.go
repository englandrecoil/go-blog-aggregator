package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/englandrecoil/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAddFeed(s *state, cmd command, currentUserDB database.User) error {
	if len(cmd.Arguments) != 2 {
		return fmt.Errorf("usage: %s <name> <url>", cmd.Name)
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

	argsToFollowFeed := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time,
		UpdatedAt: time,
		UserID:    currentUserDB.ID,
		FeedID:    feed.ID,
	}

	_, err = s.db.CreateFeedFollow(context.Background(), argsToFollowFeed)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			fmt.Println("You are already following this feed!")
			return nil
		}
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

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

func handlerFollowFeed(s *state, cmd command, currentUser database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	time := time.Now()

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("feed doesn't exist in db: %w", err)
	}

	args := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time,
		UpdatedAt: time,
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	}

	followedFeed, err := s.db.CreateFeedFollow(context.Background(), args)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			fmt.Println("You are already following this feed!")
			return nil
		}
		return fmt.Errorf("couldn't follow feed: %w", err)
	}

	fmt.Println("New feed successfully followed:")
	fmt.Printf("Feed name: %s\n", followedFeed.FeedName)
	fmt.Printf("Current user: %s\n", followedFeed.UserName)

	return nil
}

func handlerListFollowedFeeds(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	followedFeed, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get followed feeds: %w", err)
	}

	if len(followedFeed) == 0 {
		fmt.Printf("No followed feeds found for user %s\n", user.Name)
		return nil
	}

	fmt.Printf("Subscriptions for user %s:\n", user.Name)
	for index, value := range followedFeed {
		fmt.Printf("%d. %s\n", index+1, value.FeedName)
	}

	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByURL(context.Background(), cmd.Arguments[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed with provided url: %w", err)
	}

	// check if user is following provided feed to prevent "deleting" feeds that user doesnt actually follow
	isFollow, err := checkFollow(s, feed)
	if err != nil {
		return err
	}

	if !isFollow {
		fmt.Println("You're not following this feed")
		return nil
	}

	args := database.DeleteFollowedFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFollowedFeed(context.Background(), args)
	if err != nil {
		return fmt.Errorf("couldn't unfollow feed: %w", err)
	}

	fmt.Printf("Successfully unfollowed %s feed\n", feed.Name)

	return nil
}

func checkFollow(s *state, feed database.Feed) (bool, error) {
	isFollow := false
	followedFeed, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return isFollow, fmt.Errorf("couldn't get followed feeds: %w", err)
	}

	for _, value := range followedFeed {
		if value.FeedID == feed.ID {
			isFollow = true
			break
		}
	}

	return isFollow, nil
}

func printFeed(feed database.Feed, user database.User) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* Added by user:          %s\n", user.Name)
	fmt.Printf("* Last fetched:           %v\n", feed.LastFetchedAt)
}
