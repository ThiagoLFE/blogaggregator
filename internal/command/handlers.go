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

func HandlerReset(s *config.State, cmd Command) error {
	err := s.DB.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset users db state: %v", err)
	}

	err = s.DB.DeleteFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to reset feeds db state: %v", err)
	}

	fmt.Println("DB restored successfully")
	return nil
}

func HandlerGetUsers(s *config.State, _ Command) error {
	users, err := s.DB.ListUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to read users: %v", err)
	}

	for _, u := range users {
		if u.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", u.Name)
			continue
		}
		fmt.Printf("* %s\n", u.Name)
	}
	return nil
}

func HandlerAgg(s *config.State, _ Command) error {
	rssFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	printFeed(rssFeed)
	return nil
}

func printFeed(feed *RSSFeed) {
	fmt.Println()
	fmt.Println("RSS FEED")
	fmt.Printf("%s\n", feed.Channel.Description)
	fmt.Printf("%s\n", feed.Channel.Link)
	fmt.Printf("%s\n", feed.Channel.Title)

	fmt.Println()
	fmt.Println("Items:")
	for i, _ := range feed.Channel.Item {
		fmt.Printf("%s\n", feed.Channel.Item[i].Description)
		fmt.Printf("%s\n", feed.Channel.Item[i].Title)
		fmt.Printf("%s\n", feed.Channel.Item[i].Link)
		fmt.Printf("%s\n", feed.Channel.Item[i].PubDate)
	}
}

func HandlerAddFeed(s *config.State, c Command) error {
	if len(c.Args) != 2 {
		return fmt.Errorf("Error. Usage go run . addfeed <title> <url>")
	}

	u, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error to found user: %q", err)
	}

	cmd := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      c.Args[0],
		Url:       c.Args[1],
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	feed, err := s.DB.CreateFeed(context.Background(), cmd)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) {
			if pqError.Constraint == "uq_feeds_user_name" {
				return fmt.Errorf("this feed with same name and url already have been posted from you")
			}
		}
		return fmt.Errorf("fail to create feed: %q", err)
	}

	printFeedObjt(s, feed)

	return nil
}

func printFeedObjt(s *config.State, feed database.Feed) {
	author_name := fmt.Sprintf("%v", feed.UserID)
	u, err := s.DB.GetUserByID(context.Background(), feed.UserID)
	if err == nil {
		author_name = u.Name
	}

	fmt.Println("")
	fmt.Println("Feed")
	fmt.Println("")

	fmt.Printf("ID: %v\n", feed.ID)
	fmt.Printf("Name: %v\n", feed.Name)
	fmt.Printf("URL: %v\n", feed.Url)
	fmt.Printf("Author: %v\n", author_name)
	fmt.Printf("Created at: %v\n", feed.CreatedAt.Format("02/01/2006 15:04:05"))
	fmt.Printf("Updated at: %v\n", feed.UpdatedAt.Format("02/01/2006 15:04:05"))

	fmt.Println("")
}
