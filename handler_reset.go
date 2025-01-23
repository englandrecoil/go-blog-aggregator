package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Arguments) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("couldn't delete data")
	}

	fmt.Println("Data successfully deleted")
	return nil
}
