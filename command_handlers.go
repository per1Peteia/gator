package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	cfg "github.com/per1Peteia/gator/internal/config"
	"github.com/per1Peteia/gator/internal/database"
	"time"
)

type state struct {
	c  *cfg.Config
	db *database.Queries
}

func handlerList(s *state, cmd command) error {
	// loop over GetUsers slice and print user (current) if user is set to current else print user w/o (current)

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
