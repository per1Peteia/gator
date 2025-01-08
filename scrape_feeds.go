package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/per1Peteia/gator/internal/api"
	"github.com/per1Peteia/gator/internal/database"
	"log"
	"strings"
	"time"
)

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("couldn't get fet to fetch")
	}
	fmt.Println("Found feed!")
	scrapeFeed(s.db, nextFeed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	if _, err := db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Printf("couldn't mark feed fetched")
	}

	rssFeed, err := api.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("could not fetch feed")
	}

	for _, rssItem := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, rssItem.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePosts(context.Background(), database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       rssItem.Title,
			Url:         rssItem.Link,
			Description: sql.NullString{String: rssItem.Description, Valid: true},
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed '%s' collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
