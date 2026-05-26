package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/valueobject"
)

func TestArticle_IsProcessed(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		article  Article
		expected bool
	}{
		{
			name:     "nil ProcessedAt returns false",
			article:  Article{ProcessedAt: nil},
			expected: false,
		},
		{
			name: "non-nil ProcessedAt returns true",
			article: Article{
				ProcessedAt: &now,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.article.IsProcessed()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestArticle_Fields(t *testing.T) {
	now := time.Now()

	article := Article{
		ID:               "test-id",
		Title:            "AI Revolution",
		TurkishTitle:     "Yapay Zeka Devrimi",
		OriginalURL:      "https://example.com/ai",
		SourceType:       valueobject.SourceHackerNews,
		OriginalContent:  "English content here",
		TurkishContent:   "Türkçe içerik burada",
		TurkishSummary:   "Hap bilgi özeti",
		Score:            42,
		TopicID:          "topic-1",
		ImageURL:         "https://example.com/image.jpg",
		ProcessedAt:      &now,
		FetchedAt:        now,
		CreatedAt:        now,
	}

	require.Equal(t, "test-id", article.ID)
	assert.Equal(t, "AI Revolution", article.Title)
	assert.Equal(t, "Yapay Zeka Devrimi", article.TurkishTitle)
	assert.Equal(t, "https://example.com/ai", article.OriginalURL)
	assert.Equal(t, valueobject.SourceHackerNews, article.SourceType)
	assert.Equal(t, 42, article.Score)
	assert.Equal(t, "topic-1", article.TopicID)
	assert.True(t, article.IsProcessed())
}

func TestArticle_DefaultValues(t *testing.T) {
	article := Article{}

	assert.False(t, article.IsProcessed())
	assert.Empty(t, article.ID)
	assert.Empty(t, article.Title)
	assert.Equal(t, 0, article.Score)
	assert.True(t, article.FetchedAt.IsZero())
	assert.True(t, article.CreatedAt.IsZero())
}
