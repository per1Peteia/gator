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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return err
	}
	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}
	fmt.Println()
	fmt.Printf("current user has been set to: %s", cmd.args[0])
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
		return err
	}

	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("user was successfully created:\nName:%s\nCreatedAt: %v", user.Name, user.CreatedAt)

	return nil
}
