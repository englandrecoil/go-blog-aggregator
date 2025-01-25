package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Arguments) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		return errors.New("wrong time between requests. Use <value + ms/s/m/h>")
	}

	fmt.Printf("Feeds will be collected every %v...\n", timeBetweenRequests)
	fmt.Println("=====================================")

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get next feed to fetch: %w", err)
	}

	if err = s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		return fmt.Errorf("couldn't mark fetched feed: %w", err)
	}

	feedContent, err := fetchFeed(context.TODO(), feed.Url)
	if err != nil {
		return err
	}

	for _, value := range feedContent.Channel.Item {
		fmt.Printf("Found post: %s\n", value.Title)
	}
	log.Printf("Feed %s collected, %d posts found", feed.Name, len(feedContent.Channel.Item))
	fmt.Println("=====================================")

	return nil
}
