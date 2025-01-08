package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	cfg "github.com/per1Peteia/gator/internal/config"
	"github.com/per1Peteia/gator/internal/database"
)

// this struct represents application state (config and database)
type state struct {
	c  *cfg.Config
	db *database.Queries
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.args) == 1 {
		parsedLimit, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("couldn't parse limit arg to int: %w", err)
		}
		limit = parsedLimit
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("could not get posts: %w", err)
	}

	fmt.Printf("Found %d posts for user %s:\n", len(posts), user.Name)
	for _, post := range posts {
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %v\n", post.Description)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")
	}

	return nil
}

// takes 1 argument (feed url) and deletes feed_follow record for logged in user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}

	fmt.Println("deleting ...")

	if err := s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{Url: cmd.args[0], UserID: user.ID}); err != nil {
		return fmt.Errorf("couldn't delete feed follow: %w", err)
	}

	fmt.Println("Unfollow successful")

	return nil
}

// takes 1 argument (feed url) and creates a joining table which represents a user following a feed
func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <url>", cmd.name)
	}

	user, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	feedID, err := s.db.GetFeedByID(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't get feed id: %w", err)
	}

	feedFollow, err := s.db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feedID,
		},
	)

	fmt.Println("successful follow")
	fmt.Println()
	printFeedFollow(string(feedFollow.UserName), feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	user, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get user: %w", err)
	}

	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get follows: %w", err)
	}
	if len(follows) == 0 {
		fmt.Printf("%s is not following any feeds at the moment.", user.Name)
	}

	for _, follow := range follows {
		printFeedFollow(follow.UserName, follow.FeedName)
	}

	return nil
}

func printFeedFollow(username, feedname string) {
	fmt.Printf("* User:          %s\n", username)
	fmt.Printf("* Feed:          %s\n", feedname)
}

// takes no arguments and lists all records from feeds table
func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get feeds: %w", err)
	}

	if len(feeds) == 0 {
		return fmt.Errorf("you need to add feeds first.")
	}

	for _, feed := range feeds {
		fmt.Printf(" * Name:					%s\n", feed.Name)
		fmt.Printf(" * URL:						%s\n", feed.Url)
		fmt.Printf(" * Added by User:				%s\n", feed.Feedusername)
		fmt.Println()
		fmt.Println("=====================================")
	}

	return nil
}

// this function takes 2 arguments (name, url), adds a feed to the db, and connects it to the current user
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("Usage: %s <name> <feed_url>\n", cmd.name)
	}

	currentUser, err := s.db.GetUser(context.Background(), s.c.CurrentUserName)
	if err != nil {
		return fmt.Errorf("couldn't get current user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    currentUser.ID,
	},
	)

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    currentUser.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	fmt.Println()
	fmt.Println("Feed followed successfully:")
	printFeedFollow(feedFollow.UserName, feedFollow.FeedName)
	fmt.Println("=====================================")

	return nil
}

// helper function to print feeds
func printFeed(feed database.Feed) {
	fmt.Printf("* ID:            %s\n", feed.ID)
	fmt.Printf("* Created:       %v\n", feed.CreatedAt)
	fmt.Printf("* Updated:       %v\n", feed.UpdatedAt)
	fmt.Printf("* Name:          %s\n", feed.Name)
	fmt.Printf("* URL:           %s\n", feed.Url)
	fmt.Printf("* UserID:        %s\n", feed.UserID)
}

// this function takes 1 argument (duration) and will aggregate feeds
func handlerAgg(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("Usage: %s <duration [1s 1m 1h]>", cmd.name)
	}

	duration, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("%w: valid time units are s, m, h", err)
	}

	fmt.Printf("Collecting feeds every %s\n", cmd.args[0])

	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

// this function takes no arguments and prints all currently registered users to the console
func handlerList(s *state, cmd command) error {
	// loop over GetUsers slice and print '<user> (current)' if user is set to current else print user w/o (current)
	items, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't get users: %w", err)
	}
	if len(items) == 0 {
		fmt.Println("there are no users registered yet")
		return nil
	}
	for _, item := range items {
		if s.c.CurrentUserName == item {
			fmt.Printf("  * %s (current)\n", item)
		} else {
			fmt.Printf("  * %s\n", item)
		}
	}
	return nil
}

// this function takes no arguments and resets the users table to an empty valid table
func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUsers(context.Background()); err != nil {
		return fmt.Errorf("could not reset database: %w", err)
	}
	fmt.Println("users were successfully reset.")

	return nil
}

// this function takes 1 argument (name) and logs in a registered user to be the current user
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	if _, err := s.db.GetUser(context.Background(), cmd.args[0]); err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}
	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}
	fmt.Println()
	fmt.Printf("current user has successfully set to: %s", cmd.args[0])
	fmt.Println()

	return nil
}

// this function takes 1 argument (name) and registers a user to be included in the users table
func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}
	user, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      cmd.args[0],
		},
	)
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	if err := s.c.SetUser(cmd.args[0]); err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Println("User created successfully!")
	printUser(user)
	return nil
}

// helper function to print users
func printUser(user database.User) {
	fmt.Printf(" * ID:			%v\n", user.ID)
	fmt.Printf(" * Name:		%v\n", user.Name)
}
