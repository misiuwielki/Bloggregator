package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/misiuwielki/Bloggregator/internal/database"
)

func handlerFollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("url not passed")
	}
	url := cmd.arguments[0]
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("feed not found, use 'addfeed' first")
		}
		return fmt.Errorf("error on finding feed: %w", err)
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error while following feed: %w", err)
	}
	fmt.Printf("Feed %s is\n", feedFollow.FeedName)
	fmt.Printf("followed by: %s", feedFollow.UserName)
	return nil
}

func handlerFollowsForUser(s *state, cmd command, user database.User) error {
	followsS, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("error on retrieving follows: %w", err)
	}
	for _, follow := range followsS {
		fmt.Printf("%s\n", follow.FeedName)
	}
	return nil
}

func handlerUnfollowFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("url not passed")
	}
	url := cmd.arguments[0]
	err := s.db.RemoveFollowFeed(context.Background(), database.RemoveFollowFeedParams{
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("Error on processing unfollow: %w", err)
	}
	return nil
}
