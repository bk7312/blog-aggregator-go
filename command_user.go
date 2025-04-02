package main

import (
	"context"
	"fmt"
	"time"

	"github.com/bk7312/blog-aggregator-go/internal/config"
	"github.com/bk7312/blog-aggregator-go/internal/database"
	"github.com/google/uuid"
)

func handleLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage for %s requires 1 arg", cmd.name)
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not registered: %v", err)
	}

	config.SetUser(user.Name)
	*s.cfg = config.Read()
	return nil
}

func handleRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage for %s requires 1 arg", cmd.name)
	}
	username := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), username)
	if err == nil {
		return fmt.Errorf("user already registered")
	}

	user, err = s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		},
	)
	if err != nil {
		return fmt.Errorf("error registering user in db: %v", err)
	}
	config.SetUser(user.Name)
	*s.cfg = config.Read()
	fmt.Printf("user %s successfully created\n", user.Name)
	return nil
}

func handleReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handleUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		str := user.Name
		if user.Name == s.cfg.CurrentUserName {
			str += " (current)"
		}
		fmt.Printf("* %s\n", str)
	}
	return nil
}
