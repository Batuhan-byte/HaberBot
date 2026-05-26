// Package usecase implements application-level business logic.
// Use cases orchestrate domain entities and ports to fulfill specific
// application requirements. They depend only on the domain layer.
package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"github.com/go-shiori/go-readability"
)

// FetchArticlesUseCase orchestrates fetching articles from multiple content
// sources for a given topic, deduplicating by URL before persistence.
type FetchArticlesUseCase struct {
	fetchers    []port.ContentFetcher
	articleRepo port.ArticleRepository
}

// NewFetchArticlesUseCase creates a new FetchArticlesUseCase with the given
// content fetchers and article repository.
func NewFetchArticlesUseCase(fetchers []port.ContentFetcher, articleRepo port.ArticleRepository) *FetchArticlesUseCase {
	return &FetchArticlesUseCase{
		fetchers:    fetchers,
		articleRepo: articleRepo,
	}
}

// Execute fetches articles from all registered fetchers for the given topic,
// deduplicates them by URL, and persists new ones. Returns the count of
// newly fetched articles.
func (uc *FetchArticlesUseCase) Execute(ctx context.Context, topic *entity.Topic) (int, error) {
	var fetched int

	for _, fetcher := range uc.fetchers {
		if !uc.isSourceEnabled(topic, fetcher) {
			continue
		}

		articles, err := fetcher.FetchByKeywords(ctx, topic.Keywords)
		if err != nil {
			slog.Error("failed to fetch articles",
				"source", fetcher.SourceType(),
				"topic", topic.Slug,
				"error", err,
			)
			continue
		}

		saved, err := uc.saveNewArticles(ctx, articles, topic.ID)
		if err != nil {
			return fetched, fmt.Errorf("saving articles for topic %s: %w", topic.Slug, err)
		}
		fetched += saved
	}

	return fetched, nil
}

// isSourceEnabled checks whether the fetcher's source type is allowed for the topic.
func (uc *FetchArticlesUseCase) isSourceEnabled(topic *entity.Topic, fetcher port.ContentFetcher) bool {
	for _, source := range topic.Sources {
		if source == fetcher.SourceType() {
			return true
		}
	}
	return false
}

// saveNewArticles persists articles that don't already exist by URL.
func (uc *FetchArticlesUseCase) saveNewArticles(ctx context.Context, articles []*entity.Article, topicID string) (int, error) {
	var saved int

	for _, article := range articles {
		exists, err := uc.articleRepo.ExistsByURL(ctx, article.OriginalURL)
		if err != nil {
			return saved, fmt.Errorf("checking URL existence: %w", err)
		}
		if exists {
			continue
		}

		// Use go-readability to extract the full text content from the URL
		parsed, err := readability.FromURL(article.OriginalURL, 15*time.Second)
		if err == nil && parsed.TextContent != "" {
			article.OriginalContent = parsed.TextContent
		} else {
			slog.Warn("could not extract full text, falling back to summary", "url", article.OriginalURL)
		}

		article.TopicID = topicID
		if err := uc.articleRepo.Save(ctx, article); err != nil {
			return saved, fmt.Errorf("saving article: %w", err)
		}
		saved++
	}

	return saved, nil
}
