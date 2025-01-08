# gator
a rss blog aggregator - this repo is a project i built for the boot.dev course

## requirements

* this was built on macOS
* go v1.23.4 
* postgresql v12+
* pq driver: ensure the github.com/lib/pq package is included in go.mod.

## features

* add rss feeds you want to follow to your database by <name> and <url>
* allows for multiple user profiles (login functionality)
* follow/unfollow feeds to (un-)link them to your logged-in user profile
* view follows of the logged-in user
* automatically aggregate rss items regularly (prioritizing the least recently updated feeds)
* browse an optional number of posts to view in the terminal 
* list user profiles and feeds to get an overview
* reset the entire rss database if needed

## installation

```go install https://github.com/per1Peteia/gator```

## setup

* use a single JSON file (`~/.gatorconfig.json`) to keep track of two things:

   1. Who is currently logged in
   2. The connection credentials for the PostgreSQL database

* The JSON file should have this structure (when prettified):

```
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```


## usage

## what i learned

* learn how to integrate a Go application with a PostgreSQL database 
* practice using SQL skills to query and migrate a database (using sqlc and goose, two lightweight tools for typesafe SQL in Go)
* how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database


## where to go from here

* add sorting and filtering options to the browse command
* add pagination to the browse command
* add concurrency to the agg command so that it can fetch more frequently
* add a search command that allows for fuzzy searching of posts
* add bookmarking or liking posts
* add a TUI that allows you to select a post in the terminal and view it in a more readable format (either in the terminal or open in a browser)
* add an HTTP API (and authentication/authorization) that allows other users to interact with the service remotely
* write a service manager that keeps the agg command running in the background and restarts it if it crashes



