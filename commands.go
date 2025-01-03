package main

import "errors"

type command struct {
	name string
	args []string
}

type commands struct {
	callback map[string]func(*state, command) error
}

// registers a new handler function for a command name
func (c *commands) register(name string, f func(*state, command) error) {
	c.callback[name] = f
}

// runs a given command with the provided state if it exists
func (c *commands) run(s *state, cmd command) error {
	callback, exists := c.callback[cmd.name]
	if exists {
		return callback(s, cmd)
	}
	return errors.New("command does not exist")
}
