package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/englandrecoil/go-blog-aggregator/internal/database"
	"github.com/google/uuid"
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.Arguments) != 0 {
		parsedLimit, err := strconv.Atoi(cmd.Arguments[0])
		if err != nil {
			return fmt.Errorf("can't set this limit. Use %s <limit> or leave blank", cmd.Name)
		}
		if parsedLimit <= 0 {
			return fmt.Errorf("limit must be a positive number")
		}
		limit = parsedLimit
	}

	args := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), args)
	if err != nil {
		return err
	}

	fmt.Printf("Found posts for user %s\n:", user.Name)
	for _, post := range posts {
		fmt.Printf("Published at %s from %s\n", post.PublishedAt, post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("	%v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
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
		// fmt.Printf("Found post: %s, desc: %s\n", value.Title, value.Description)

		pubDate, err := parsePublishDate(value.PubDate)
		if err != nil {
			fmt.Println(err)
			pubDate = time.Now()
		}

		args := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       value.Title,
			Url:         value.Link,
			Description: value.Description,
			PublishedAt: pubDate,
			FeedID:      feed.ID,
		}

		if _, err := s.db.CreatePost(context.Background(), args); err != nil {
			return err
		}
	}
	log.Printf("Feed %s collected, %d posts found", feed.Name, len(feedContent.Channel.Item))
	fmt.Println("=====================================")

	return nil
}

func parsePublishDate(timeDate string) (time.Time, error) {
	pubDate, err := time.Parse(time.RFC1123Z, timeDate)
	if err == nil {
		return pubDate, nil
	}

	pubDate, err = time.Parse(time.RFC1123, timeDate)
	if err == nil {
		return pubDate, nil
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", timeDate)
}
