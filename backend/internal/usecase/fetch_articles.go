// Package usecase implements application-level business logic.
// Use cases orchestrate domain entities and ports to fulfill specific
// application requirements. They depend only on the domain layer.
package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
)

// FetchArticlesUseCase orchestrates fetching articles from multiple content
// sources for a given topic, deduplicating by URL before persistence.
type FetchArticlesUseCase struct {
	fetchers     []port.ContentFetcher
	articleRepo  port.ArticleRepository
	pendingLimit int
}

// NewFetchArticlesUseCase creates a new FetchArticlesUseCase with the given
// content fetchers, article repository, and pending limit.
func NewFetchArticlesUseCase(fetchers []port.ContentFetcher, articleRepo port.ArticleRepository, pendingLimit int) *FetchArticlesUseCase {
	return &FetchArticlesUseCase{
		fetchers:     fetchers,
		articleRepo:  articleRepo,
		pendingLimit: pendingLimit,
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

		articles, err := fetcher.Fetch(ctx, topic)
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
		// Skip non-article repository links (e.g. GitHub, GitLab, Bitbucket)
		urlLower := strings.ToLower(article.OriginalURL)
		if strings.Contains(urlLower, "github.com") || strings.Contains(urlLower, "gitlab.com") || strings.Contains(urlLower, "bitbucket.org") {
			slog.Info("skipping non-article code repository URL", "url", article.OriginalURL)
			continue
		}

		exists, err := uc.articleRepo.ExistsByURL(ctx, article.OriginalURL)
		if err != nil {
			return saved, fmt.Errorf("checking URL existence: %w", err)
		}
		if exists {
			continue
		}

		// Use the new helper to get clean HTML content and image
		htmlContent, scrapedImage := scrapeFullArticle(ctx, article.OriginalURL)
		if htmlContent != "" {
			article.OriginalContent = htmlContent
		} else {
			slog.Warn("could not extract full HTML, falling back to summary/title", "url", article.OriginalURL)
			if article.OriginalContent == "" {
				article.OriginalContent = fmt.Sprintf("<p>%s</p>", article.Title)
			}
		}

		if article.ImageURL == "" && scrapedImage != "" {
			article.ImageURL = scrapedImage
		}

		article.TopicID = topicID
		if err := uc.articleRepo.Save(ctx, article); err != nil {
			return saved, fmt.Errorf("saving article: %w", err)
		}
		
		// Trim pending articles for this topic to enforce the limit
		if err := uc.articleRepo.TrimPendingByTopic(ctx, topicID, uc.pendingLimit); err != nil {
			slog.Error("failed to trim pending articles after insert", "topic_id", topicID, "error", err)
		}

		saved++
	}

	return saved, nil
}

// TrimPending triggers a trim operation on the repository for the given topicID and configured limit.
func (uc *FetchArticlesUseCase) TrimPending(ctx context.Context, topicID string) error {
	return uc.articleRepo.TrimPendingByTopic(ctx, topicID, uc.pendingLimit)
}

// scrapeFullArticle fetches the article page HTML and extracts clean readability content and og:image meta tags.
func scrapeFullArticle(ctx context.Context, targetURL string) (string, string) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return "", ""
	}
	// Standard desktop browser User-Agent to avoid scraping blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "tr-TR,tr;q=0.9,en-US;q=0.8,en;q=0.7")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn("HTTP request failed during scraping", "url", targetURL, "error", err)
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("unexpected status code during scraping", "url", targetURL, "status", resp.StatusCode)
		return "", ""
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ""
	}

	// 1. Extract og:image using goquery
	var scrapedImage string
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyBytes))
	if err == nil {
		doc.Find("meta").Each(func(i int, s *goquery.Selection) {
			if prop, exists := s.Attr("property"); exists && prop == "og:image" {
				scrapedImage, _ = s.Attr("content")
			}
			if name, exists := s.Attr("name"); exists && (name == "twitter:image" || name == "og:image") && scrapedImage == "" {
				scrapedImage, _ = s.Attr("content")
			}
		})
	}

	// 2. Extract clean HTML content using go-readability
	var htmlContent string
	parsedURL, _ := url.Parse(targetURL)
	parsed, err := readability.FromReader(bytes.NewReader(bodyBytes), parsedURL)
	if err == nil && parsed.Content != "" {
		htmlContent = parsed.Content
	} else {
		// Fallback to text content if HTML content parsing fails
		if err == nil && parsed.TextContent != "" {
			htmlContent = fmt.Sprintf("<p>%s</p>", parsed.TextContent)
		}
	}

	return htmlContent, scrapedImage
}

