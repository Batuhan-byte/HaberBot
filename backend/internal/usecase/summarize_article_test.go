package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

type mockAIProcessor struct {
	port.AIProcessor
	translateFunc func(ctx context.Context, title, content string) (string, string, error)
	summarizeFunc func(ctx context.Context, content string) (string, error)
}

func (m *mockAIProcessor) Translate(ctx context.Context, title, content string) (string, string, error) {
	return m.translateFunc(ctx, title, content)
}

func (m *mockAIProcessor) Summarize(ctx context.Context, content string) (string, error) {
	return m.summarizeFunc(ctx, content)
}

func TestSummarizeArticleUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("new summary generated", func(t *testing.T) {
		article := makeArticle("1", "Test")
		article.OriginalContent = "English content"
		article.TurkishContent = "Türkçe içerik"

		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, id string) (*entity.Article, error) {
				assert.Equal(t, "1", id)
				return article, nil
			},
			saveFunc: func(_ context.Context, a *entity.Article) error {
				assert.Equal(t, "Türkçe özet", a.TurkishSummary)
				return nil
			},
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, content string) (string, error) {
				assert.Equal(t, "Türkçe içerik", content)
				return "Türkçe özet", nil
			},
		}
		uc := NewSummarizeArticleUseCase(ai, repo)
		summary, err := uc.Execute(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "Türkçe özet", summary)
	})

	t.Run("existing summary returned immediately", func(t *testing.T) {
		article := makeArticle("1", "Test")
		article.TurkishSummary = "Mevcut özet"

		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, id string) (*entity.Article, error) {
				return article, nil
			},
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, _ string) (string, error) {
				t.Fatal("should not be called")
				return "", nil
			},
		}
		uc := NewSummarizeArticleUseCase(ai, repo)
		summary, err := uc.Execute(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "Mevcut özet", summary)
	})

	t.Run("fallback to original content when TurkishContent empty", func(t *testing.T) {
		article := makeArticle("1", "Test")
		article.OriginalContent = "English only"
		article.TurkishContent = ""

		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return article, nil
			},
			saveFunc: func(_ context.Context, _ *entity.Article) error { return nil },
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, content string) (string, error) {
				assert.Equal(t, "English only", content)
				return "Özet", nil
			},
		}
		uc := NewSummarizeArticleUseCase(ai, repo)
		summary, err := uc.Execute(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "Özet", summary)
	})

	t.Run("article not found", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, nil
			},
		}
		uc := NewSummarizeArticleUseCase(&mockAIProcessor{}, repo)
		_, err := uc.Execute(ctx, "1")
		assert.ErrorContains(t, err, "article not found")
	})

	t.Run("find error", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, errors.New("db error")
			},
		}
		uc := NewSummarizeArticleUseCase(&mockAIProcessor{}, repo)
		_, err := uc.Execute(ctx, "1")
		assert.ErrorContains(t, err, "finding article")
	})

	t.Run("ai error", func(t *testing.T) {
		article := makeArticle("1", "Test")
		article.TurkishContent = "İçerik"

		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return article, nil
			},
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, _ string) (string, error) {
				return "", errors.New("ai down")
			},
		}
		uc := NewSummarizeArticleUseCase(ai, repo)
		_, err := uc.Execute(ctx, "1")
		assert.ErrorContains(t, err, "ai summarize")
	})

	t.Run("save error", func(t *testing.T) {
		article := makeArticle("1", "Test")
		article.TurkishContent = "İçerik"

		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return article, nil
			},
			saveFunc: func(_ context.Context, _ *entity.Article) error {
				return errors.New("save failed")
			},
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, _ string) (string, error) {
				return "Özet", nil
			},
		}
		uc := NewSummarizeArticleUseCase(ai, repo)
		_, err := uc.Execute(ctx, "1")
		assert.ErrorContains(t, err, "saving summary")
	})
}
