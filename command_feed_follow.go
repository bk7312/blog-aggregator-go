package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bk7312/blog-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func handleFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage for %s requires 1 arg", cmd.name)
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s followed %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handleFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}
	fmt.Printf("%s following %v feeds\n", user.Name, len(feeds))
	for _, feed := range feeds {
		fmt.Println(feed.Name)
	}
	return nil
}

func handleUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage for %s requires 1 arg", cmd.name)
	}
	url := cmd.args[0]
	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return err
	}
	fmt.Printf("%s unfollowed %s\n", user.Name, feed.Name)
	return nil
}
