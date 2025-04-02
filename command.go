package main

import (
	"fmt"
	"sort"

	"github.com/bk7312/blog-aggregator-go/internal/config"
	"github.com/bk7312/blog-aggregator-go/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("command not found: %s", cmd.name)
	}
	return handler(s, cmd)
}

func handleHelp(cmds commands) error {
	fmt.Println("List of available commands:")
	c := make([]string, 0, len(cmds.handlers))
	for h := range cmds.handlers {
		c = append(c, h)
	}
	sort.Strings(c)
	for _, key := range c {
		fmt.Println(key)
	}
	return nil
}

func passCmds(handler func(commands) error, c *commands) func(*state, command) error {
	return func(s *state, cmd command) error {
		return handler(*c)
	}
}
