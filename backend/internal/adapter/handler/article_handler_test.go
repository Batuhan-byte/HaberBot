package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
	"haberbot/internal/usecase"
)

func TestArticleHandler_GetRecentArticles(t *testing.T) {
	t.Run("returns articles", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, limit int) ([]*entity.Article, error) {
				assert.Equal(t, 12, limit)
				return []*entity.Article{{ID: "1", Title: "Test"}}, nil
			},
		}
		listUC := usecase.NewListArticlesUseCase(repo)
		h := NewArticleHandler(listUC, nil, nil, nil)

		app := fiber.New()
		app.Get("/articles", h.GetRecentArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Contains(t, result, "articles")
	})

	t.Run("custom limit query param", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, limit int) ([]*entity.Article, error) {
				assert.Equal(t, 5, limit)
				return []*entity.Article{}, nil
			},
		}
		listUC := usecase.NewListArticlesUseCase(repo)
		h := NewArticleHandler(listUC, nil, nil, nil)

		app := fiber.New()
		app.Get("/articles", h.GetRecentArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles?limit=5", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("invalid limit falls back to 12", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, limit int) ([]*entity.Article, error) {
				assert.Equal(t, 12, limit)
				return []*entity.Article{}, nil
			},
		}
		listUC := usecase.NewListArticlesUseCase(repo)
		h := NewArticleHandler(listUC, nil, nil, nil)

		app := fiber.New()
		app.Get("/articles", h.GetRecentArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles?limit=-1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("usecase error returns 500", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return nil, errors.New("db error")
			},
		}
		listUC := usecase.NewListArticlesUseCase(repo)
		h := NewArticleHandler(listUC, nil, nil, nil)

		app := fiber.New()
		app.Get("/articles", h.GetRecentArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestArticleHandler_SearchArticles(t *testing.T) {
	t.Run("search with query", func(t *testing.T) {
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
				assert.Equal(t, "AI", query)
				return []*entity.Article{{ID: "1", Title: "AI News"}}, nil
			},
		}
		searchUC := usecase.NewSearchArticlesUseCase(repo)
		h := NewArticleHandler(nil, nil, searchUC, nil)

		app := fiber.New()
		app.Get("/articles/search", h.SearchArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles/search?q=AI", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Equal(t, "AI", result["query"])
	})

	t.Run("search with source filter", func(t *testing.T) {
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, query string, source valueobject.SourceType, _, _ int) ([]*entity.Article, error) {
				assert.Equal(t, "hackernews", source.String())
				return []*entity.Article{}, nil
			},
		}
		searchUC := usecase.NewSearchArticlesUseCase(repo)
		h := NewArticleHandler(nil, nil, searchUC, nil)

		app := fiber.New()
		app.Get("/articles/search", h.SearchArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles/search?q=AI&source=hackernews", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("empty query returns recent articles", func(t *testing.T) {
		repo := &mockArticleRepo{
			findRecentFunc: func(_ context.Context, _ int) ([]*entity.Article, error) {
				return []*entity.Article{}, nil
			},
		}
		searchUC := usecase.NewSearchArticlesUseCase(repo)
		listUC := usecase.NewListArticlesUseCase(repo)
		h := NewArticleHandler(listUC, nil, searchUC, nil)

		app := fiber.New()
		app.Get("/articles/search", h.SearchArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles/search", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("search error returns 500", func(t *testing.T) {
		repo := &mockArticleRepo{
			searchFunc: func(_ context.Context, _ string, _ valueobject.SourceType, _, _ int) ([]*entity.Article, error) {
				return nil, errors.New("search error")
			},
		}
		searchUC := usecase.NewSearchArticlesUseCase(repo)
		h := NewArticleHandler(nil, nil, searchUC, nil)

		app := fiber.New()
		app.Get("/articles/search", h.SearchArticles)
		req := httptest.NewRequest(http.MethodGet, "/articles/search?q=AI", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestArticleHandler_GetArticleByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, id string) (*entity.Article, error) {
				assert.Equal(t, "1", id)
				return &entity.Article{ID: "1", Title: "Test"}, nil
			},
		}
		getUC := usecase.NewGetArticleUseCase(repo)
		h := NewArticleHandler(nil, getUC, nil, nil)

		app := fiber.New()
		app.Get("/articles/:id", h.GetArticleByID)
		req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("not found returns 404", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, nil
			},
		}
		getUC := usecase.NewGetArticleUseCase(repo)
		h := NewArticleHandler(nil, getUC, nil, nil)

		app := fiber.New()
		app.Get("/articles/:id", h.GetArticleByID)
		req := httptest.NewRequest(http.MethodGet, "/articles/999", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("repo error returns 500", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, errors.New("db error")
			},
		}
		getUC := usecase.NewGetArticleUseCase(repo)
		h := NewArticleHandler(nil, getUC, nil, nil)

		app := fiber.New()
		app.Get("/articles/:id", h.GetArticleByID)
		req := httptest.NewRequest(http.MethodGet, "/articles/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestArticleHandler_SummarizeArticle(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		article := &entity.Article{ID: "1", Title: "Test", TurkishContent: "İçerik"}
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, id string) (*entity.Article, error) {
				return article, nil
			},
			saveFunc: func(_ context.Context, _ *entity.Article) error { return nil },
		}
		ai := &mockAIProcessor{
			summarizeFunc: func(_ context.Context, content string) (string, error) {
				return "Özet metni", nil
			},
		}
		sumUC := usecase.NewSummarizeArticleUseCase(ai, repo)
		h := NewArticleHandler(nil, nil, nil, sumUC)

		app := fiber.New()
		app.Post("/articles/:id/summary", h.SummarizeArticle)
		req := httptest.NewRequest(http.MethodPost, "/articles/1/summary", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Equal(t, "Özet metni", result["summary"])
	})

	t.Run("error returns 500", func(t *testing.T) {
		repo := &mockArticleRepo{
			findByIDFunc: func(_ context.Context, _ string) (*entity.Article, error) {
				return nil, errors.New("find error")
			},
		}
		sumUC := usecase.NewSummarizeArticleUseCase(&mockAIProcessor{}, repo)
		h := NewArticleHandler(nil, nil, nil, sumUC)

		app := fiber.New()
		app.Post("/articles/:id/summary", h.SummarizeArticle)
		req := httptest.NewRequest(http.MethodPost, "/articles/1/summary", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
