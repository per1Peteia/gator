package main

import (
	"fmt"
	cfg "github.com/per1Peteia/gator/internal/config"
)

type state struct {
	c *cfg.Config
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}
	fmt.Println()
	fmt.Printf("current user has been set to: %s", cmd.args[0])
	fmt.Println()

	return nil
}
