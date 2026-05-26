package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
)

func TestSearchArticlesUseCase_Search(t *testing.T) {
	ctx := context.Background()

	t.Run("search with query", func(t *testing.T) {
		expected := []*entity.Article{makeArticle("1", "AI News")}
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
				assert.Equal(t, "AI", query)
				assert.Equal(t, valueobject.SourceType(""), source)
				assert.Equal(t, 10, limit)
				assert.Equal(t, 0, offset)
				return expected, nil
			},
		}
		uc := NewSearchArticlesUseCase(repo)
		got, err := uc.Search(ctx, "AI", "", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("search with source filter", func(t *testing.T) {
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, query string, source valueobject.SourceType, _, _ int) ([]*entity.Article, error) {
				assert.Equal(t, "AI", query)
				assert.Equal(t, "hackernews", source.String())
				return []*entity.Article{}, nil
			},
		}
		uc := NewSearchArticlesUseCase(repo)
		_, err := uc.Search(ctx, "AI", valueobject.SourceHackerNews, 10, 0)
		require.NoError(t, err)
	})

	t.Run("empty query and source falls back to FindRecent", func(t *testing.T) {
		expected := []*entity.Article{makeArticle("1", "Recent")}
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, limit int) ([]*entity.Article, error) {
				assert.Equal(t, 10, limit)
				return expected, nil
			},
		}
		uc := NewSearchArticlesUseCase(repo)
		got, err := uc.Search(ctx, "", "", 10, 0)
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("search repo error", func(t *testing.T) {
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, _ string, _ valueobject.SourceType, _, _ int) ([]*entity.Article, error) {
				return nil, errors.New("search failed")
			},
		}
		uc := NewSearchArticlesUseCase(repo)
		_, err := uc.Search(ctx, "AI", "", 10, 0)
		assert.Error(t, err)
	})
}
