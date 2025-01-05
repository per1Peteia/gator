package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/per1Peteia/gator/internal/api"
	cfg "github.com/per1Peteia/gator/internal/config"
	"github.com/per1Peteia/gator/internal/database"
	"time"
)

type state struct {
	c  *cfg.Config
	db *database.Queries
}

func handlerAddfeed(s *state, cmd command) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <name> <feed_url>\n", cmd.name)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUser.ID,
	},
	)
	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("=====================================")

	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

func handlerAgg(s *state, cmd command) error {
	url := "https://www.wagslane.dev/index.xml"
	rssFeed, err := api.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("could not fetch RSS Feed: %w", err)
	}
	rssFeed = api.UnescapeStrings(rssFeed)
	fmt.Printf("%v", *rssFeed)
	return nil
}

func handlerList(s *state, cmd command) error {
	// loop over GetUsers slice and print '<user> (current)' if user is set to current else print user w/o (current)
	items, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}
	if len(items) == 0 {
		fmt.Println("there are no users registered yet")
		return nil
	}
	for _, item := range items {
		if s.c.CurrentUserName == item {
			fmt.Printf("  * %s (current)\n", item)
		} else {
			fmt.Printf("  * %s\n", item)
		}
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("could not reset database: %w", err)
	}
	fmt.Println("users were successfully reset.")

	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}
	fmt.Println()
	fmt.Printf("current user has successfully set to: %s", cmd.args[0])
	fmt.Println()

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Println("User created successfully!")
	printUser(user)
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:			%v\n", user.ID)
	fmt.Printf(" * Name:		%v\n", user.Name)
}
