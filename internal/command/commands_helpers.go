package commands

import (
	"blogaggregator/internal/config"
	"fmt"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("username is required to login")
	}
	name := cmd.Args[0]

	if err := s.Config.SetUser(name); err != nil {
		return err
	}

	fmt.Printf("Welcome %s, the loggin was succefully!\n", cmd.Args[0])
	return nil
}
