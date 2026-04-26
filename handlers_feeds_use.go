package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/misiuwielki/Bloggregator/internal/database"
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
			published, err := time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				return fmt.Errorf("Error on parsing date: %w", err)
			}
			post, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
				ID:    uuid.New(),
				Title: item.Title,
				Url:   item.Link,
				Description: sql.NullString{
					String: item.Description,
					Valid:  item.Description != ""},
				PublishedAt: published,
				FeedID:      feed.ID,
			})
			if err != nil {
				if pqErr, ok := err.(*pq.Error); ok {
					if pqErr.Code == "23505" {
						continue
					}
				}
				return fmt.Errorf("Error on saving post %w", err)
			}
			fmt.Printf("Saved post: %s\n", post.Title)
		}
	}
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.arguments) > 0 {
		if cmd.arguments[0] == "nolimit" {
			limit = 99999999
		} else {
			var err error
			limit, err = strconv.Atoi(cmd.arguments[0])
			if err != nil {
				return fmt.Errorf("limist must be a number or 'nolimit")
			}
		}
	}
	postS, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("Error on retrieving posts %w", err)
	}
	for _, post := range postS {
		description := ""
		if post.Description.Valid {
			description = post.Description.String
		}
		fmt.Printf("Title: %v\n", post.Title)
		fmt.Printf("Text: %v\n", description)
		fmt.Printf("Published at: %v\n", post.PublishedAt)
	}
	return nil
}
