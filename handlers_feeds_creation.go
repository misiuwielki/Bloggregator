package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/misiuwielki/Bloggregator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
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

func fetchFeed(ctx context.Context, feedUrl string) (*RSSFeed, error) {
	rq, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error on creating request for feed %w", err)
	}
	rq.Header.Set("User-Agent", "gator")
	fmt.Printf("fetching %s\n", feedUrl)
	rsp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error on getting response for feed %w", err)
	}
	defer rsp.Body.Close()
	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error on reading response for feed %w", err)
	}
	feedS := RSSFeed{}
	err = xml.Unmarshal(data, &feedS)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("error on unmarshaling response for feed %w", err)
	}
	feedS.Channel.Title = html.UnescapeString(feedS.Channel.Title)
	feedS.Channel.Description = html.UnescapeString(feedS.Channel.Description)
	for i, item := range feedS.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		feedS.Channel.Item[i] = item
	}
	return &feedS, nil

}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("too few arguments - need name and url")
	}
	name := cmd.arguments[0]
	url := cmd.arguments[1]
	feed, err := s.db.AddFeed(context.Background(), database.AddFeedParams{
		ID:     uuid.New(),
		Name:   name,
		Url:    url,
		UserID: user.ID,
	})
	if err != nil {
		return fmt.Errorf("error on adding feed %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:     uuid.New(),
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error on adding feed %w", err)
	}
	fmt.Printf("%v\n", feed)
	return nil
}
