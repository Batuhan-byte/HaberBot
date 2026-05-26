package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
	"haberbot/internal/domain/port"
)

type mockArticleRepo struct {
	port.ArticleRepository
	findRecentFunc    func(ctx context.Context, limit int) ([]*entity.Article, error)
	findByIDFunc      func(ctx context.Context, id string) (*entity.Article, error)
	findByTopicIDFunc func(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error)
	findUnprocessedFunc func(ctx context.Context, limit int) ([]*entity.Article, error)
	existsByURLFunc   func(ctx context.Context, url string) (bool, error)
	countByTopicIDFunc func(ctx context.Context, topicID string) (int, error)
	saveFunc          func(ctx context.Context, article *entity.Article) error
	searchFunc        func(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error)
}

func (m *mockArticleRepo) FindRecent(ctx context.Context, limit int) ([]*entity.Article, error) {
	return m.findRecentFunc(ctx, limit)
}
func (m *mockArticleRepo) FindByID(ctx context.Context, id string) (*entity.Article, error) {
	return m.findByIDFunc(ctx, id)
}
func (m *mockArticleRepo) FindByTopicID(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
	return m.findByTopicIDFunc(ctx, topicID, limit, offset)
}
func (m *mockArticleRepo) FindUnprocessed(ctx context.Context, limit int) ([]*entity.Article, error) {
	return m.findUnprocessedFunc(ctx, limit)
}
func (m *mockArticleRepo) ExistsByURL(ctx context.Context, url string) (bool, error) {
	return m.existsByURLFunc(ctx, url)
}
func (m *mockArticleRepo) CountByTopicID(ctx context.Context, topicID string) (int, error) {
	return m.countByTopicIDFunc(ctx, topicID)
}
func (m *mockArticleRepo) Save(ctx context.Context, article *entity.Article) error {
	return m.saveFunc(ctx, article)
}
func (m *mockArticleRepo) Search(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
	return m.searchFunc(ctx, query, source, limit, offset)
}

func makeArticle(id, title string) *entity.Article {
	return &entity.Article{
		ID:    id,
		Title: title,
	}
}

func TestListArticlesUseCase_ListRecent(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expected := []*entity.Article{makeArticle("1", "Test")}
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, limit int) ([]*entity.Article, error) {
				assert.Equal(t, 10, limit)
				return expected, nil
			},
		}
		uc := NewListArticlesUseCase(repo)
		got, err := uc.ListRecent(ctx, 10)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return nil, errors.New("db error")
			},
		}
		uc := NewListArticlesUseCase(repo)
		_, err := uc.ListRecent(ctx, 10)
		assert.ErrorContains(t, err, "listing recent articles")
	})
}

func TestListArticlesUseCase_ListByTopic(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		articles := []*entity.Article{makeArticle("1", "Article 1"), makeArticle("2", "Article 2")}

		repo := &mockArticleRepo{
			findByTopicIDFunc: func(_ context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
				assert.Equal(t, "topic-1", topicID)
				assert.Equal(t, 10, limit)
				assert.Equal(t, 0, offset)
				return articles, nil
			},
			countByTopicIDFunc: func(_ context.Context, topicID string) (int, error) {
				assert.Equal(t, "topic-1", topicID)
				return 5, nil
			},
		}
		uc := NewListArticlesUseCase(repo)
		page, err := uc.ListByTopic(ctx, "topic-1", 1, 10)
		require.NoError(t, err)
		require.NotNil(t, page)
		assert.Equal(t, articles, page.Articles)
		assert.Equal(t, 5, page.Total)
	})

	t.Run("find error", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByTopicIDFunc: func(_ context.Context, _ string, _, _ int) ([]*entity.Article, error) {
				return nil, errors.New("find error")
			},
		}
		uc := NewListArticlesUseCase(repo)
		_, err := uc.ListByTopic(ctx, "topic-1", 1, 10)
		assert.ErrorContains(t, err, "listing articles by topic")
	})

	t.Run("count error", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByTopicIDFunc: func(_ context.Context, _ string, _, _ int) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
			countByTopicIDFunc: func(_ context.Context, _ string) (int, error) {
				return 0, errors.New("count error")
			},
		}
		uc := NewListArticlesUseCase(repo)
		_, err := uc.ListByTopic(ctx, "topic-1", 1, 10)
		assert.ErrorContains(t, err, "counting articles by topic")
	})

	t.Run("page 2 offset calculation", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByTopicIDFunc: func(_ context.Context, _ string, limit, offset int) ([]*entity.Article, error) {
				assert.Equal(t, 10, limit)
				assert.Equal(t, 10, offset)
				return []*entity.Article{}, nil
			},
			countByTopicIDFunc: func(_ context.Context, _ string) (int, error) {
				return 0, nil
			},
		}
		uc := NewListArticlesUseCase(repo)
		page, err := uc.ListByTopic(ctx, "topic-1", 2, 10)
		require.NoError(t, err)
		assert.Equal(t, 0, page.Total)
	})
}
