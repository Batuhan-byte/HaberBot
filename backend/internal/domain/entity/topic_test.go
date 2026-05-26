package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"haberbot/internal/domain/valueobject"
)

func TestTopic_Fields(t *testing.T) {
	now := time.Now()

	topic := Topic{
		ID:        "topic-1",
		Name:      "Yapay Zeka",
		Slug:      "yapay-zeka",
		Keywords:  []valueobject.TopicKeyword{"ai", "machine learning"},
		Sources:   []valueobject.SourceType{valueobject.SourceHackerNews, valueobject.SourceRSS},
		RSSFeeds:  []string{"https://example.com/rss"},
		IsActive:  true,
		CreatedAt: now,
	}

	assert.Equal(t, "topic-1", topic.ID)
	assert.Equal(t, "Yapay Zeka", topic.Name)
	assert.Equal(t, "yapay-zeka", topic.Slug)
	assert.Len(t, topic.Keywords, 2)
	assert.Len(t, topic.Sources, 2)
	assert.Len(t, topic.RSSFeeds, 1)
	assert.True(t, topic.IsActive)
}

func TestTopic_Inactive(t *testing.T) {
	topic := Topic{
		ID:       "topic-2",
		Name:     "Inactive Topic",
		Slug:     "inactive",
		IsActive: false,
	}

	assert.False(t, topic.IsActive)
	assert.Empty(t, topic.Keywords)
	assert.Empty(t, topic.Sources)
}

func TestTopic_EmptyRSSFeeds(t *testing.T) {
	topic := Topic{
		ID:       "topic-3",
		Name:     "No RSS",
		Slug:     "no-rss",
		IsActive: true,
	}

	assert.Empty(t, topic.RSSFeeds)
}
