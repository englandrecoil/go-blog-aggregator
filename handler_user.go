package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Arguments) < 1 {
		return fmt.Errorf("no arguments provided")
	}

	if len(cmd.Arguments) > 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}

	if err := s.cfg.SetUser(cmd.Arguments[0]); err != nil {
		return fmt.Errorf("couldnt set current user: %w", err)
	}

	fmt.Printf("User %s has been set\n", s.cfg.CurrentUserName)

	return nil
}
