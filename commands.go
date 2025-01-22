package main

import "fmt"

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	commandHandler, exists := c.registeredCommands[cmd.Name]
	if !exists {
		return fmt.Errorf("run error: command %s doesn't exist", cmd.Name)
	}

	if err := commandHandler(s, cmd); err != nil {
		return fmt.Errorf("%s error: %w", cmd.Name, err)
	}

	return nil
}
