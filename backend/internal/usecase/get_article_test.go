package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
)

func TestGetArticleUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("found", func(t *testing.T) {
		expected := makeArticle("1", "Test Article")
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, id string) (*entity.Article, error) {
				assert.Equal(t, "1", id)
				return expected, nil
			},
		}
		uc := NewGetArticleUseCase(repo)
		got, err := uc.Execute(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("not found returns nil", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, nil
			},
		}
		uc := NewGetArticleUseCase(repo)
		got, err := uc.Execute(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, errors.New("db error")
			},
		}
		uc := NewGetArticleUseCase(repo)
		_, err := uc.Execute(ctx, "1")
		assert.ErrorContains(t, err, "getting article")
	})
}
