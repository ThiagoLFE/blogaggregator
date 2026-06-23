package commands

import (
	"blogaggregator/internal/config"
	"errors"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommands map[string]func(*config.State, Command) error
}

func (c *Commands) Register(name string, f func(*config.State, Command) error) {
	c.RegisteredCommands[name] = f
}

func (c *Commands) Run(s *config.State, cmd Command) error {
	f, ok := c.RegisteredCommands[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}
	return f(s, cmd)
}
