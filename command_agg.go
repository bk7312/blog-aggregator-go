package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/bk7312/blog-aggregator-go/internal/database"
	"github.com/google/uuid"
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

func handleAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage for %s requires 1 arg", cmd.name)
	}
	interval, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("invalid time interval")
	}
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			return fmt.Errorf("error scraping feed: %v", err)
		}
	}
}

func handleBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err == nil {
			limit = parsedLimit
		}
	}
	posts, err := s.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting post: %v", err)
	}
	fmt.Printf("Showing %v articles for %s\n", len(posts), user.Name)
	fmt.Println("----------")
	for _, post := range posts {
		fmt.Println("Title:", html.UnescapeString(post.Title))
		fmt.Println("Link:", html.UnescapeString(post.Url))
		fmt.Println("Description:", html.UnescapeString(post.Description.String))
		fmt.Println("PubDate:", html.UnescapeString(post.PublishedAt.String()))
		fmt.Println("-----")
	}
	return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "gator")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var data RSSFeed
	err = xml.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}
	newTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: newTime,
		UpdatedAt:     time.Now(),
	})
	fmt.Printf("Fetched %v articles for %v\n", len(feed.Channel.Item), feed.Channel.Title)
	fmt.Println("----------")

	for _, item := range feed.Channel.Item {
		pubTime, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			return fmt.Errorf("unable to parse time: %v", err)
		}
		var desc sql.NullString
		if item.Description == "" {
			desc = sql.NullString{
				Valid: false,
			}
		} else {
			desc = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		_, err = s.db.GetPostByUrl(context.Background(), item.Link)
		if err == nil {
			continue
		}

		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: pubTime,
			FeedID:      nextFeed.ID,
		})
		if err != nil {
			return fmt.Errorf("error saving post: %v", err)
		}
		fmt.Println("Title:", html.UnescapeString(item.Title))
	}
	return nil
}
