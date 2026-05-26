package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"haberbot/internal/domain/valueobject"
	"haberbot/internal/usecase"
)

func TestAdminHandler_AuthMiddleware(t *testing.T) {
	t.Run("valid key passes", func(t *testing.T) {
		h := NewAdminHandler(nil, nil, "secret-key")
		app := fiber.New()
		adminGroup := app.Group("/admin")
		adminGroup.Use(h.AuthMiddleware())
		adminGroup.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("X-Admin-API-Key", "secret-key")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("missing key returns 401", func(t *testing.T) {
		h := NewAdminHandler(nil, nil, "secret-key")
		app := fiber.New()
		adminGroup := app.Group("/admin")
		adminGroup.Use(h.AuthMiddleware())
		adminGroup.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("wrong key returns 401", func(t *testing.T) {
		h := NewAdminHandler(nil, nil, "secret-key")
		app := fiber.New()
		adminGroup := app.Group("/admin")
		adminGroup.Use(h.AuthMiddleware())
		adminGroup.Get("/test", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"ok": true})
		})

		req := httptest.NewRequest(http.MethodGet, "/admin/test", nil)
		req.Header.Set("X-Admin-API-Key", "wrong-key")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAdminHandler_CreateTopic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, topic *entity.Topic) error {
				assert.Equal(t, "AI", topic.Name)
				return nil
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")

		app := fiber.New()
		app.Post("/admin/topics", h.CreateTopic)

		body := `{"name":"AI","slug":"ai","keywords":["ai"],"sources":["hackernews"],"is_active":true}`
		req := httptest.NewRequest(http.MethodPost, "/admin/topics", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("missing name returns 400", func(t *testing.T) {
		repo := &mockTopicRepo{}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Post("/admin/topics", h.CreateTopic)

		body := `{"slug":"ai"}`
		req := httptest.NewRequest(http.MethodPost, "/admin/topics", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("invalid JSON returns 400", func(t *testing.T) {
		h := NewAdminHandler(usecase.NewManageTopicsUseCase(&mockTopicRepo{}), nil, "key")
		app := fiber.New()
		app.Post("/admin/topics", h.CreateTopic)

		req := httptest.NewRequest(http.MethodPost, "/admin/topics", strings.NewReader("{invalid}"))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("usecase error returns 500", func(t *testing.T) {
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, _ *entity.Topic) error {
				return errors.New("create failed")
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Post("/admin/topics", h.CreateTopic)

		body := `{"name":"AI","slug":"ai"}`
		req := httptest.NewRequest(http.MethodPost, "/admin/topics", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAdminHandler_UpdateTopic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockTopicRepo{
			saveFunc: func(_ context.Context, topic *entity.Topic) error {
				assert.Equal(t, "1", topic.ID)
				assert.Equal(t, "Updated", topic.Name)
				return nil
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Put("/admin/topics/:id", h.UpdateTopic)

		body := `{"name":"Updated","slug":"updated"}`
		req := httptest.NewRequest(http.MethodPut, "/admin/topics/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("missing id returns 400", func(t *testing.T) {
		repo := &mockTopicRepo{}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Put("/admin/topics/:id", h.UpdateTopic)

		body := `{}`
		req := httptest.NewRequest(http.MethodPut, "/admin/topics/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestAdminHandler_DeleteTopic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockTopicRepo{
			deleteFunc: func(_ context.Context, id string) error {
				assert.Equal(t, "1", id)
				return nil
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Delete("/admin/topics/:id", h.DeleteTopic)

		req := httptest.NewRequest(http.MethodDelete, "/admin/topics/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var result map[string]interface{}
		json.Unmarshal(body, &result)
		assert.Equal(t, "topic deleted successfully", result["message"])
	})

	t.Run("usecase error returns 500", func(t *testing.T) {
		repo := &mockTopicRepo{
			deleteFunc: func(_ context.Context, _ string) error {
				return errors.New("delete failed")
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(repo)
		h := NewAdminHandler(manageUC, nil, "key")
		app := fiber.New()
		app.Delete("/admin/topics/:id", h.DeleteTopic)

		req := httptest.NewRequest(http.MethodDelete, "/admin/topics/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAdminHandler_TriggerFetch(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fetchUC := usecase.NewFetchArticlesUseCase(
			[]port.ContentFetcher{&mockContentFetcher{
				sourceTypeFunc: func() valueobject.SourceType { return valueobject.SourceHackerNews },
				fetchByKeywordsFunc: func(_ context.Context, _ []valueobject.TopicKeyword) ([]*entity.Article, error) {
					return nil, nil
				},
			}},
			&mockArticleRepo{},
		)
		processUC := usecase.NewProcessArticlesUseCase(&mockAIProcessor{}, &mockArticleRepo{})
		topicRepo := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return []*entity.Topic{}, nil
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(topicRepo)
		pipelineUC := usecase.NewDailyPipelineUseCase(fetchUC, processUC, manageUC)
		h := NewAdminHandler(manageUC, pipelineUC, "key")

		app := fiber.New()
		app.Post("/admin/fetch", h.TriggerFetch)
		req := httptest.NewRequest(http.MethodPost, "/admin/fetch", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("error returns 500", func(t *testing.T) {
		topicRepo := &mockTopicRepo{
			findAllFunc: func(_ context.Context) ([]*entity.Topic, error) {
				return nil, errors.New("fetch error")
			},
		}
		manageUC := usecase.NewManageTopicsUseCase(topicRepo)
		fetchUC := usecase.NewFetchArticlesUseCase([]port.ContentFetcher{&mockContentFetcher{}}, &mockArticleRepo{})
		processUC := usecase.NewProcessArticlesUseCase(&mockAIProcessor{}, &mockArticleRepo{})
		pipelineUC := usecase.NewDailyPipelineUseCase(fetchUC, processUC, manageUC)
		h := NewAdminHandler(manageUC, pipelineUC, "key")
		app := fiber.New()
		app.Post("/admin/fetch", h.TriggerFetch)
		req := httptest.NewRequest(http.MethodPost, "/admin/fetch", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
