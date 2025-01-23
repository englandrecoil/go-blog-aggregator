package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link>chardata"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("couldn't create new request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("can't send request to server: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return &RSSFeed{}, fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}

	feed := RSSFeed{}
	// decode xml data from server to struct
	bodyData, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error reading response body: %s", err)
	}

	if err := xml.Unmarshal(bodyData, &feed); err != nil {
		return &RSSFeed{}, fmt.Errorf("error decoding response body: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for index := range feed.Channel.Item {
		feed.Channel.Item[index].Title = html.UnescapeString(feed.Channel.Item[index].Title)
		feed.Channel.Item[index].Description = html.UnescapeString(feed.Channel.Item[index].Description)
	}

	return &feed, nil
}
