package main

import (
	"database/sql"
	"github.com/per1Peteia/gator/internal/config"
	"github.com/per1Peteia/gator/internal/database"
	"log"
	"os"
)

// the underscore tells go that this is imported for side effects not usage
import _ "github.com/lib/pq"

func main() {
	config, err := cfg.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	db, err := sql.Open("postgres", "postgres://peripeteia:@localhost:5432/gator?sslmode=disable")
	dbQueries := database.New(db)

	appState := state{
		c:  &config,
		db: dbQueries,
	}

	commands := commands{
		callback: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerList)
	commands.register("agg", handlerAgg)
	commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	commands.register("feeds", handlerListFeeds)
	commands.register("follow", middlewareLoggedIn(handlerFollow))
	commands.register("following", middlewareLoggedIn(handlerFollowing))
	commands.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	commands.register("browse", middlewareLoggedIn(handlerBrowse))

	if len(os.Args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}
	command := command{
		name: os.Args[1],
		args: os.Args[2:],
	}
	if err := commands.run(&appState, command); err != nil {
		log.Fatal(err)
	}
}
