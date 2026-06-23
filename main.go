package main

import (
	commands "blogaggregator/internal/command"
	"blogaggregator/internal/config"
	"blogaggregator/internal/database"
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const DATA_BASE_URL = "postgres://postgres:postgres@localhost:5432/gator"

func main() {
	// Reading the config file
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// Connecting with db
	db, err := sql.Open("postgres", DATA_BASE_URL)
	if err != nil {
		panic(err.Error())
	}

	// Taking our queries
	dbqueries := database.New(db)

	// Setting our State
	programState := &config.State{
		Config: &cfg,
		DB:     dbqueries,
	}

	cmds := commands.Commands{
		RegisteredCommands: make(map[string]func(*config.State, commands.Command) error),
	}

	cmds.Register("login", commands.HandlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli<command> [args...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	err = cmds.Run(programState, commands.Command{
		Name: cmdName,
		Args: cmdArgs,
	})
	if err != nil {
		log.Fatal(err)
	}
}
