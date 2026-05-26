// Package port defines the interfaces (ports) through which the domain
// communicates with the outside world. All adapters must implement these
// interfaces. This follows the Ports & Adapters (Hexagonal) pattern.
package port

import (
	"context"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// ArticleRepository defines the contract for article persistence operations.
// Implementations may use PostgreSQL, in-memory storage, or any other backend.
type ArticleRepository interface {
	// Save persists a new article or updates an existing one.
	Save(ctx context.Context, article *entity.Article) error

	// FindByID retrieves a single article by its unique identifier.
	// Returns nil and no error if the article is not found.
	FindByID(ctx context.Context, id string) (*entity.Article, error)

	// FindByTopicID retrieves articles belonging to a specific topic
	// with pagination support via limit and offset.
	FindByTopicID(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error)

	// FindRecent retrieves the most recently fetched articles, limited
	// by the given count.
	FindRecent(ctx context.Context, limit int) ([]*entity.Article, error)

	// FindUnprocessed retrieves articles that have not yet been processed
	// by the AI pipeline, limited by the given count.
	FindUnprocessed(ctx context.Context, limit int) ([]*entity.Article, error)

	// ExistsByURL checks whether an article with the given URL already exists.
	ExistsByURL(ctx context.Context, url string) (bool, error)

	// CountByTopicID returns the total number of articles for a given topic.
	CountByTopicID(ctx context.Context, topicID string) (int, error)

	// Search articles by title/content and/or source.
	// Returns articles matching the search query and/or source with pagination.
	Search(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error)
}
