// Package entity defines the core domain entities for the HaberBot application.
// Entities encapsulate enterprise-wide business rules and are the innermost layer
// of the Clean Architecture.
package entity

import (
	"time"

	"haberbot/internal/domain/valueobject"
)

// Article represents a fetched and optionally processed news article.
// Articles are sourced from external providers (HackerNews, RSS) and may be
// translated and summarized into Turkish by an AI processor.
type Article struct {
	// ID is the unique identifier for the article.
	ID string `json:"id"`
	// Title is the original English title of the article.
	Title string `json:"title"`
	// TurkishTitle is the AI-translated Turkish title.
	TurkishTitle string `json:"title_tr"`
	// OriginalURL is the canonical URL where the article was published.
	OriginalURL string `json:"original_url"`
	// SourceType indicates the origin source (hackernews, rss).
	SourceType valueobject.SourceType `json:"source"`
	// OriginalContent holds the raw English content or description.
	OriginalContent string `json:"original_content"`
	// TurkishContent is the AI-translated Turkish content.
	TurkishContent string `json:"content_tr"`
	// TurkishSummary is the AI-generated Turkish "hap bilgi" summary.
	TurkishSummary string `json:"summary_tr"`
	// Score is the relevance or popularity score from the source.
	Score int `json:"score"`
	// TopicID references the topic this article belongs to.
	TopicID string `json:"topic_id"`
	// ImageURL is an optional image associated with the article.
	ImageURL string `json:"image_url"`
	// ProcessedAt records when the article was processed by AI. Nil means unprocessed.
	ProcessedAt *time.Time `json:"processed_at"`
	// FetchedAt records when the article was fetched from its source.
	FetchedAt time.Time `json:"fetched_at"`
	// CreatedAt records when the article was persisted in the database.
	CreatedAt time.Time `json:"created_at"`
	// IsApproved indicates if the article has been approved by an administrator.
	IsApproved bool `json:"is_approved"`
	// IsHidden indicates if the article is temporarily hidden/soft-deleted.
	IsHidden bool `json:"is_hidden"`
	// ApprovedAt records when the article was approved by an administrator. Nil means not approved.
	ApprovedAt *time.Time `json:"approved_at"`
}

// IsProcessed returns true if the article has been processed by the AI pipeline,
// indicated by a non-nil ProcessedAt timestamp.
func (a *Article) IsProcessed() bool {
	return a.ProcessedAt != nil
}
