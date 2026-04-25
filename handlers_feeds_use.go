package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func handlerGetAllFeeds(s *state, cmd command) error {
	feedS, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error on retrieving feeds: %w", err)
	}
	for _, feed := range feedS {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error on identifying user: %w", err)
		}
		fmt.Printf("Feed name: %s\n", feed.Name)
		fmt.Printf("Feed url: %s\n", feed.Url)
		fmt.Printf("Added by: %s\n", user.Name)

	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("time between requests is required, minimum 30 second")
	}
	repTime, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("Invalid time value: %w", err)
	}
	if repTime < 1*time.Minute {
		repTime = 30 * time.Second
		fmt.Println("Timer set to minimum value: 30 seconds")
	}
	ticker := time.NewTicker(repTime)
	for ; ; <-ticker.C {
		err := scrapeFeed(s)
		if err != nil {
			return err
		}
	}
}

func scrapeFeed(s *state) error {
	fmt.Println("starting scraping")
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("to aggregate, add a feed first with 'addfeed' command")
		}
		return fmt.Errorf("error while retrieving feed: %w", err)
	}
	fmt.Println("found feed")
	err = s.db.MarkFeedFetched(context.Background(), feed.Name)
	if err != nil {
		return fmt.Errorf("error while marking fetch: %w", err)
	}
	fmt.Println("feed marked as fetched")
	feedS, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("error while fetching feed: %w", err)
	}
	fmt.Printf("feed %s fetched\n", feedS.Channel.Title)
	for i, item := range feedS.Channel.Item {
		if i < 10 {
			fmt.Printf("Title: %s\n", item.Title)
		}
	}
	return nil
}
