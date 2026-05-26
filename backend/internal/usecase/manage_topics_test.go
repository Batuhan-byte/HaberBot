package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"haberbot/internal/domain/valueobject"
)

type mockTopicRepo struct {
	port.TopicRepository
	saveFunc    func(ctx context.Context, topic *entity.Topic) error
	findAllFunc func(ctx context.Context) ([]*entity.Topic, error)
	findBySlugFunc func(ctx context.Context, slug string) (*entity.Topic, error)
	deleteFunc  func(ctx context.Context, id string) error
	findByIDFunc func(ctx context.Context, id string) (*entity.Topic, error)
}

func (m *mockTopicRepo) Save(ctx context.Context, topic *entity.Topic) error   { return m.saveFunc(ctx, topic) }
func (m *mockTopicRepo) FindAll(ctx context.Context) ([]*entity.Topic, error)    { return m.findAllFunc(ctx) }
func (m *mockTopicRepo) FindBySlug(ctx context.Context, slug string) (*entity.Topic, error) { return m.findBySlugFunc(ctx, slug) }
func (m *mockTopicRepo) Delete(ctx context.Context, id string) error            { return m.deleteFunc(ctx, id) }
func (m *mockTopicRepo) FindByID(ctx context.Context, id string) (*entity.Topic, error) { return m.findByIDFunc(ctx, id) }

func TestManageTopicsUseCase_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		topic := &entity.Topic{Name: "AI", Slug: "ai"}
		var savedTopic *entity.Topic
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, st *entity.Topic) error {
				savedTopic = st
				return nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Create(ctx, topic)
		require.NoError(t, err)
		assert.Equal(t, topic, savedTopic)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, _ *entity.Topic) error {
				return errors.New("save error")
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Create(ctx, &entity.Topic{})
		assert.ErrorContains(t, err, "creating topic")
	})
}

func TestManageTopicsUseCase_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []*entity.Topic{{Name: "AI", Slug: "ai"}}
		repo := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return expected, nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		got, err := uc.GetAll(ctx)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("empty list", func(t *testing.T) {
		repo := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return []*entity.Topic{}, nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		got, err := uc.GetAll(ctx)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return nil, errors.New("db error")
			},
		}
		uc := NewManageTopicsUseCase(repo)
		_, err := uc.GetAll(ctx)
		assert.ErrorContains(t, err, "listing all topics")
	})
}

func TestManageTopicsUseCase_GetBySlug(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		expected := &entity.Topic{Name: "AI", Slug: "ai"}
		repo := &mockTopicRepo{
			findBySlugFunc: func(_ context.Context, slug string) (*entity.Topic, error) {
				assert.Equal(t, "ai", slug)
				return expected, nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		got, err := uc.GetBySlug(ctx, "ai")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found", func(t *testing.T) {
		repo := &mockTopicRepo{
			findBySlugFunc: func(_ context.Context, _ string) (*entity.Topic, error) {
				return nil, nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		got, err := uc.GetBySlug(ctx, "unknown")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockTopicRepo{
			findBySlugFunc: func(_ context.Context, _ string) (*entity.Topic, error) {
				return nil, errors.New("db error")
			},
		}
		uc := NewManageTopicsUseCase(repo)
		_, err := uc.GetBySlug(ctx, "ai")
		assert.ErrorContains(t, err, "finding topic by slug")
	})
}

func TestManageTopicsUseCase_Update(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		topic := &entity.Topic{ID: "1", Name: "Updated"}
		var savedTopic *entity.Topic
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, st *entity.Topic) error {
				savedTopic = st
				return nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Update(ctx, topic)
		require.NoError(t, err)
		assert.Equal(t, topic, savedTopic)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, _ *entity.Topic) error {
				return errors.New("update error")
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Update(ctx, &entity.Topic{})
		assert.ErrorContains(t, err, "updating topic")
	})
}

func TestManageTopicsUseCase_Delete(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := &mockTopicRepo{
			deleteFunc: func(_ context.Context, id string) error {
				assert.Equal(t, "1", id)
				return nil
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Delete(ctx, "1")
		require.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockTopicRepo{
			deleteFunc: func(_ context.Context, _ string) error {
				return errors.New("delete error")
			},
		}
		uc := NewManageTopicsUseCase(repo)
		err := uc.Delete(ctx, "1")
		assert.ErrorContains(t, err, "deleting topic")
	})
}

type mockContentFetcher struct {
	port.ContentFetcher
	fetchByKeywordsFunc func(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error)
	sourceTypeFunc      func() valueobject.SourceType
}

func (m *mockContentFetcher) FetchByKeywords(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error) {
	return m.fetchByKeywordsFunc(ctx, keywords)
}
func (m *mockContentFetcher) SourceType() valueobject.SourceType { return m.sourceTypeFunc() }

func TestFetchArticlesUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("fetches and saves new articles", func(t *testing.T) {
		articles := []*entity.Article{
			{Title: "AI News", OriginalURL: "https://example.com/1"},
			{Title: "ML Update", OriginalURL: "https://example.com/2"},
		}

		fetcher := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceHackerNews },
			fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
				return articles, nil
			},
		}

		var savedURLs []string
		repo := &mockArticleRepo{
			existsByURLFunc: func(_ context.Context, url string) (bool, error) {
				return false, nil
			},
			saveFunc: func(_ context.Context, article *entity.Article) error {
				savedURLs = append(savedURLs, article.OriginalURL)
				return nil
			},
		}

		topic := &entity.Topic{
			ID:      "topic-1",
			Slug:    "ai",
			Keywords: []valueobject.TopicKeyword{"ai"},
			Sources:  []valueobject.SourceType{valueobject.SourceHackerNews},
		}

		uc := NewFetchArticlesUseCase([]port.ContentFetcher{fetcher}, repo)
		count, err := uc.Execute(ctx, topic)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, savedURLs, 2)
	})

	t.Run("skips existing URLs", func(t *testing.T) {
		fetcher := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceHackerNews },
			fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
				return []*entity.Article{
					{Title: "Duplicate", OriginalURL: "https://example.com/dup"},
				}, nil
			},
		}

		repo := &mockArticleRepo{
			existsByURLFunc: func(_ context.Context, url string) (bool, error) {
				assert.Equal(t, "https://example.com/dup", url)
				return true, nil
			},
		}

		uc := NewFetchArticlesUseCase([]port.ContentFetcher{fetcher}, repo)
		count, err := uc.Execute(ctx, &entity.Topic{
			ID: "topic-1", Slug: "ai",
			Keywords: []valueobject.TopicKeyword{"ai"},
			Sources: []valueobject.SourceType{valueobject.SourceHackerNews},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("skips disabled sources", func(t *testing.T) {
		fetcher := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceRSS },
		}

		uc := NewFetchArticlesUseCase([]port.ContentFetcher{fetcher}, &mockArticleRepo{})
		count, err := uc.Execute(ctx, &entity.Topic{
			ID: "topic-1", Slug: "ai",
			Keywords: []valueobject.TopicKeyword{"ai"},
			Sources:  []valueobject.SourceType{valueobject.SourceHackerNews},
		})
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("continues on fetcher error", func(t *testing.T) {
		fetcher1 := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceHackerNews },
			fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
				return nil, errors.New("network error")
			},
		}
		fetcher2 := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceRSS },
			fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
				return []*entity.Article{{Title: "RSS Article", OriginalURL: "https://rss.com/1"}}, nil
			},
		}

		repo := &mockArticleRepo{
			existsByURLFunc: func(_ context.Context, _ string) (bool, error) { return false, nil },
			saveFunc:        func(_ context.Context, _ *entity.Article) error { return nil },
		}

		uc := NewFetchArticlesUseCase([]port.ContentFetcher{fetcher1, fetcher2}, repo)
		count, err := uc.Execute(ctx, &entity.Topic{
			ID: "topic-1", Slug: "ai",
			Keywords: []valueobject.TopicKeyword{"ai"},
			Sources:  []valueobject.SourceType{valueobject.SourceHackerNews, valueobject.SourceRSS},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}

func TestProcessArticlesUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("processes unprocessed articles", func(t *testing.T) {
		articles := []*entity.Article{
			{ID: "1", Title: "AI News", OriginalContent: "Content 1"},
			{ID: "2", Title: "ML Update", OriginalContent: "Content 2"},
		}

		repo := &mockArticleRepo{
			findUnprocessedFunc: func(_ context.Context, batchSize int) ([]*entity.Article, error) {
				assert.Equal(t, 5, batchSize)
				return articles, nil
			},
			saveFunc: func(_ context.Context, article *entity.Article) error {
				assert.NotEmpty(t, article.TurkishTitle)
				assert.NotEmpty(t, article.TurkishContent)
				assert.NotNil(t, article.ProcessedAt)
				return nil
			},
		}

		ai := &mockAIProcessor{
			translateFunc: func(_ context.Context, title, content string) (string, string, error) {
				return "Türkçe: " + title, "Türkçe: " + content, nil
			},
		}

		uc := NewProcessArticlesUseCase(ai, repo)
		count, err := uc.Execute(ctx, 5)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("continues on AI error", func(t *testing.T) {
		articles := []*entity.Article{
			{ID: "1", Title: "OK", OriginalContent: "Content"},
			{ID: "2", Title: "Fail", OriginalContent: "Bad"},
		}

		callCount := 0
		repo := &mockArticleRepo{
			findUnprocessedFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return articles, nil
			},
			saveFunc: func(_ context.Context, _ *entity.Article) error { return nil },
		}

		ai := &mockAIProcessor{
			translateFunc: func(_ context.Context, title, content string) (string, string, error) {
				callCount++
				if callCount == 1 {
					return "Title", "Content", nil
				}
				return "", "", errors.New("ai error")
			},
		}

		uc := NewProcessArticlesUseCase(ai, repo)
		count, err := uc.Execute(ctx, 5)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})

	t.Run("no unprocessed articles", func(t *testing.T) {
		repo := &mockArticleRepo{
			findUnprocessedFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
		}
		uc := NewProcessArticlesUseCase(&mockAIProcessor{}, repo)
		count, err := uc.Execute(ctx, 5)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("context cancellation", func(t *testing.T) {
		articles := []*entity.Article{
			{ID: "1", Title: "Test", OriginalContent: "Content"},
		}

		repo := &mockArticleRepo{
			findUnprocessedFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return articles, nil
			},
		}
		ai := &mockAIProcessor{
			translateFunc: func(_ context.Context, _, _ string) (string, string, error) {
				return "T", "C", nil
			},
		}

		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		uc := NewProcessArticlesUseCase(ai, repo)
		_, err := uc.Execute(cancelledCtx, 5)
		assert.ErrorContains(t, err, "context cancelled")
	})
}

func TestDailyPipelineUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("ExecuteFetch fetches for all active topics", func(t *testing.T) {
		topics := []*entity.Topic{
			{ID: "1", Slug: "ai", IsActive: true},
			{ID: "2", Slug: "ml", IsActive: true},
			{ID: "3", Slug: "inactive", IsActive: false},
		}

		topicUC := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return topics, nil
			},
		}

		fetchCalled := 0
		fetchUC := &FetchArticlesUseCase{}
		*fetchUC = *NewFetchArticlesUseCase([]port.ContentFetcher{}, &mockArticleRepo{})

		processUC := NewProcessArticlesUseCase(&mockAIProcessor{}, &mockArticleRepo{})

		manageUC := NewManageTopicsUseCase(topicUC)
		// Override fetchUC.Execute to count calls using mock fetcher
		mockFetch := &mockContentFetcher{
			sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceHackerNews },
			fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
				fetchCalled++
				return []*entity.Article{}, nil
			},
		}
		*fetchUC = *NewFetchArticlesUseCase([]port.ContentFetcher{mockFetch}, &mockArticleRepo{
			existsByURLFunc: func(_ context.Context, _ string) (bool, error) { return false, nil },
			saveFunc:        func(_ context.Context, _ *entity.Article) error { return nil },
		})

		// Create topics with matching sources so fetcher is enabled
		for _, t := range topics {
			if t.IsActive {
				t.Sources = []valueobject.SourceType{valueobject.SourceHackerNews}
				t.Keywords = []valueobject.TopicKeyword{"ai"}
			}
		}

		uc := NewDailyPipelineUseCase(fetchUC, processUC, manageUC)
		total, err := uc.ExecuteFetch(ctx)
		require.NoError(t, err)
		assert.Equal(t, 0, total)
		assert.Equal(t, 2, fetchCalled)
	})

	t.Run("ExecuteProcess delegates to processUC", func(t *testing.T) {
		processed := 0
		repo := &mockArticleRepo{
			findUnprocessedFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return []*entity.Article{{ID: "1", Title: "Test", OriginalContent: "Content"}}, nil
			},
			saveFunc: func(_ context.Context, _ *entity.Article) error {
				processed++
				return nil
			},
		}
		ai := &mockAIProcessor{
			translateFunc: func(_ context.Context, title, content string) (string, string, error) {
				return "T " + title, "C " + content, nil
			},
		}

		pUC := NewProcessArticlesUseCase(ai, repo)
		uc := NewDailyPipelineUseCase(
			NewFetchArticlesUseCase([]port.ContentFetcher{}, repo),
			pUC,
			NewManageTopicsUseCase(&mockTopicRepo{
				findAllFunc: func(_ context.Context) ([]*entity.Topic, error) { return []*entity.Topic{}, nil },
			}),
		)

		count, err := uc.ExecuteProcess(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, count)
	})
}
