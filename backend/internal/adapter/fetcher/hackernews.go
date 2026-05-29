// Package fetcher provides concrete implementations of the ContentFetcher
// port for various external content sources.
package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

// hackerNewsBaseURL is the HackerNews Firebase API base URL.
const hackerNewsBaseURL = "https://hacker-news.firebaseio.com/v0"

// maxTopStories is the number of top stories to fetch from HackerNews.
const maxTopStories = 100

// httpClientTimeout is the timeout for individual HTTP requests.
const httpClientTimeout = 10 * time.Second

// HackerNewsFetcher implements port.ContentFetcher for the HackerNews API.
type HackerNewsFetcher struct {
	client *http.Client
}

// NewHackerNewsFetcher creates a new HackerNewsFetcher with a configured
// HTTP client.
func NewHackerNewsFetcher() *HackerNewsFetcher {
	return &HackerNewsFetcher{
		client: &http.Client{Timeout: httpClientTimeout},
	}
}

// hackerNewsItem represents a single HackerNews story item from the API.
type hackerNewsItem struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
	Type  string `json:"type"`
}

// Fetch retrieves top stories from HackerNews that match the topic's keywords.
func (f *HackerNewsFetcher) Fetch(ctx context.Context, topic *entity.Topic) ([]*entity.Article, error) {
	return f.FetchByKeywords(ctx, topic.Keywords)
}

// FetchByKeywords retrieves top stories from HackerNews that match any of
// the given keywords in their title.
func (f *HackerNewsFetcher) FetchByKeywords(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error) {
	storyIDs, err := f.fetchTopStoryIDs(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetching top stories: %w", err)
	}

	var articles []*entity.Article
	for _, storyID := range storyIDs {
		item, err := f.fetchItem(ctx, storyID)
		if err != nil {
			continue
		}

		if item.URL == "" || item.Type != "story" {
			continue
		}

		if !matchesKeywords(item.Title, keywords) {
			continue
		}

		articles = append(articles, f.toArticle(item))
	}

	return articles, nil
}

// SourceType returns SourceHackerNews.
func (f *HackerNewsFetcher) SourceType() valueobject.SourceType {
	return valueobject.SourceHackerNews
}

// fetchTopStoryIDs retrieves the list of top story IDs, capped at maxTopStories.
func (f *HackerNewsFetcher) fetchTopStoryIDs(ctx context.Context) ([]int, error) {
	url := fmt.Sprintf("%s/topstories.json", hackerNewsBaseURL)

	body, err := f.httpGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var ids []int
	if err := json.Unmarshal(body, &ids); err != nil {
		return nil, fmt.Errorf("decoding story IDs: %w", err)
	}

	if len(ids) > maxTopStories {
		ids = ids[:maxTopStories]
	}
	return ids, nil
}

// fetchItem retrieves a single HackerNews item by its ID.
func (f *HackerNewsFetcher) fetchItem(ctx context.Context, id int) (*hackerNewsItem, error) {
	url := fmt.Sprintf("%s/item/%d.json", hackerNewsBaseURL, id)

	body, err := f.httpGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var item hackerNewsItem
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, fmt.Errorf("decoding item %d: %w", id, err)
	}
	return &item, nil
}

// httpGet performs a GET request and returns the response body bytes.
func (f *HackerNewsFetcher) httpGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	var buf []byte
	buf, err = readBody(resp)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}
	return buf, nil
}

// toArticle converts a HackerNews item into a domain Article entity.
func (f *HackerNewsFetcher) toArticle(item *hackerNewsItem) *entity.Article {
	return &entity.Article{
		Title:           item.Title,
		OriginalURL:     item.URL,
		SourceType:      valueobject.SourceHackerNews,
		OriginalContent: item.Title,
		Score:           item.Score,
		FetchedAt:       time.Now(),
	}
}

// matchesKeywords checks if the text contains any of the given keywords.
func matchesKeywords(text string, keywords []valueobject.TopicKeyword) bool {
	lowerText := strings.ToLower(text)
	for _, keyword := range keywords {
		if strings.Contains(lowerText, strings.ToLower(keyword.String())) {
			return true
		}
	}
	return false
}

// readBody reads the full response body into a byte slice.
func readBody(resp *http.Response) ([]byte, error) {
	const maxBodySize = 1 << 20 // 1 MB
	buf := make([]byte, 0, 4096)
	for {
		if len(buf) >= maxBodySize {
			return buf, nil
		}
		readBuf := make([]byte, 4096)
		n, err := resp.Body.Read(readBuf)
		buf = append(buf, readBuf[:n]...)
		if err != nil {
			if err.Error() == "EOF" {
				return buf, nil
			}
			return buf, err
		}
	}
}
