package port

import (
	"context"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// ContentFetcher defines the contract for fetching articles from an external
// content source. Each implementation targets a specific source (HackerNews, RSS, etc.).
type ContentFetcher interface {
	// FetchByKeywords retrieves articles from the external source that match
	// any of the given keywords. Returns a slice of newly constructed articles.
	FetchByKeywords(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error)

	// SourceType returns the type of source this fetcher targets.
	SourceType() valueobject.SourceType
}
