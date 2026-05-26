package entity

import (
	"time"

	"haberbot/internal/domain/valueobject"
)

// Topic represents a content category that defines which articles to fetch.
// Each topic has a set of keywords for filtering, allowed source types,
// and optional RSS feed URLs.
type Topic struct {
	// ID is the unique identifier for the topic.
	ID string `json:"id"`
	// Name is the human-readable display name (e.g. "Yapay Zeka").
	Name string `json:"name"`
	// Slug is the URL-friendly identifier (e.g. "yapay-zeka").
	Slug string `json:"slug"`
	// Keywords is the list of keywords used to filter articles for this topic.
	Keywords []valueobject.TopicKeyword `json:"keywords"`
	// Sources lists the source types enabled for this topic.
	Sources []valueobject.SourceType `json:"sources"`
	// RSSFeeds holds the RSS feed URLs to be crawled for this topic.
	RSSFeeds []string `json:"rss_feeds"`
	// IsActive indicates whether this topic is currently being fetched.
	IsActive bool `json:"is_active"`
	// CreatedAt records when the topic was created.
	CreatedAt time.Time `json:"created_at"`
}
