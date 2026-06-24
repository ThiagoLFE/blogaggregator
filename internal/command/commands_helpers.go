package commands

import (
	"blogaggregator/internal/config"
	"blogaggregator/internal/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func HandlerLogin(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("username is required to login")
	}
	name := cmd.Args[0]

	user, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Printf("The current user not exists, please register with cmd register <name> first\n")
			os.Exit(1)
		}
		return err
	}

	if err := s.Config.SetUser(user.Name); err != nil {
		return err
	}

	fmt.Printf("Welcome %s, the loggin was succefully!\n", cmd.Args[0])
	return nil
}

func HandlerRegistration(s *config.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return errors.New("please insert a username to be registred")
	}
	name := cmd.Args[0]

	user, err := s.DB.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      name,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && pqErr.Constraint == "users_name_key" {
				fmt.Printf("username %q already exists\n", name)
				os.Exit(1)
			}
		}
		return err
	}

	if err := s.Config.SetUser(user.Name); err != nil {
		return err
	}
	fmt.Println("===================================")
	fmt.Println("==== Registered successfully!! ====")
	fmt.Println("===================================")
	fmt.Printf("ID: %s\n", user.ID.String())
	fmt.Printf("Name: %s\n", user.Name)
	fmt.Printf("Created at: %v\n", user.CreatedAt.Format("02/01/2006 15:04:05"))
	fmt.Printf("Updated At: %v\n", user.UpdatedAt.Format("02/01/2006 15:04:05"))
	fmt.Println("===================================")
	fmt.Printf("You are logged as %s now\n", user.Name)
	fmt.Println("===================================")

	return nil
}
