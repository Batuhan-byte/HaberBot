package port

import (
	"context"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// ContentFetcher defines the contract for fetching articles from an external
// content source. Each implementation targets a specific source (HackerNews, RSS, etc.).
type ContentFetcher interface {
	// Fetch retrieves articles from the external source for the given topic.
	Fetch(ctx context.Context, topic *entity.Topic) ([]*entity.Article, error)

	// SourceType returns the type of source this fetcher targets.
	SourceType() valueobject.SourceType
}
