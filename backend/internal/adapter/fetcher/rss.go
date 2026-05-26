package fetcher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/mmcdole/gofeed"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// RSSFetcher implements port.ContentFetcher for RSS feeds.
type RSSFetcher struct {
	parser    *gofeed.Parser
	feeds     []string
	sanitizer *HTMLSanitizer
}

// NewRSSFetcher creates a new RSSFetcher that parses the given RSS feed URLs.
func NewRSSFetcher(feeds []string) *RSSFetcher {
	return &RSSFetcher{
		parser:    gofeed.NewParser(),
		feeds:     feeds,
		sanitizer: NewHTMLSanitizer(),
	}
}

// FetchByKeywords retrieves items from all configured RSS feeds that match
// any of the given keywords in their title or description.
func (f *RSSFetcher) FetchByKeywords(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error) {
	var articles []*entity.Article

	for _, feedURL := range f.feeds {
		feedArticles, err := f.parseFeed(ctx, feedURL, keywords)
		if err != nil {
			slog.Error("failed to parse RSS feed",
				"url", feedURL, "error", err,
			)
			continue
		}
		articles = append(articles, feedArticles...)
	}

	return articles, nil
}

// SourceType returns SourceRSS.
func (f *RSSFetcher) SourceType() valueobject.SourceType {
	return valueobject.SourceRSS
}

// parseFeed parses a single RSS feed and returns articles matching keywords.
func (f *RSSFetcher) parseFeed(ctx context.Context, feedURL string, keywords []valueobject.TopicKeyword) ([]*entity.Article, error) {
	feed, err := f.parser.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("parsing feed %s: %w", feedURL, err)
	}

	var articles []*entity.Article
	for _, item := range feed.Items {
		if !f.itemMatchesKeywords(item, keywords) {
			continue
		}
		articles = append(articles, f.toArticle(item))
	}

	return articles, nil
}

// itemMatchesKeywords checks if an RSS item's title or description matches
// any of the given keywords.
func (f *RSSFetcher) itemMatchesKeywords(item *gofeed.Item, keywords []valueobject.TopicKeyword) bool {
	titleMatch := matchesKeywords(item.Title, keywords)
	descriptionMatch := matchesKeywords(item.Description, keywords)
	return titleMatch || descriptionMatch
}

// toArticle converts a gofeed.Item into a domain Article entity.
func (f *RSSFetcher) toArticle(item *gofeed.Item) *entity.Article {
	content := item.Description
	if item.Content != "" {
		content = item.Content
	}

	if f.sanitizer != nil {
		if sanitized, err := f.sanitizer.Sanitize(content); err == nil && sanitized != "" {
			content = sanitized
		}
	}

	imageURL := extractImageURL(item)

	return &entity.Article{
		Title:           item.Title,
		OriginalURL:     item.Link,
		SourceType:      valueobject.SourceRSS,
		OriginalContent: content,
		ImageURL:        imageURL,
		FetchedAt:       time.Now(),
	}
}

// extractImageURL attempts to extract an image URL from the feed item.
func extractImageURL(item *gofeed.Item) string {
	if item.Image != nil && item.Image.URL != "" {
		return item.Image.URL
	}

	if len(item.Enclosures) > 0 {
		for _, enc := range item.Enclosures {
			if enc.Type == "image/jpeg" || enc.Type == "image/png" {
				return enc.URL
			}
		}
	}

	return ""
}
