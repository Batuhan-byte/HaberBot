package handler

import (
	"github.com/gofiber/fiber/v2"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
	"haberbot/internal/usecase"
)

// AdminHandler handles administrator endpoints, protected by API key auth.
type AdminHandler struct {
	manageTopicsUC *usecase.ManageTopicsUseCase
	pipelineUC     *usecase.DailyPipelineUseCase
	adminAPIKey    string
}

// NewAdminHandler creates a new AdminHandler.
func NewAdminHandler(manageTopicsUC *usecase.ManageTopicsUseCase, pipelineUC *usecase.DailyPipelineUseCase, adminAPIKey string) *AdminHandler {
	return &AdminHandler{
		manageTopicsUC: manageTopicsUC,
		pipelineUC:     pipelineUC,
		adminAPIKey:    adminAPIKey,
	}
}

// AuthMiddleware authenticates requests against the admin API key.
func (h *AdminHandler) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Admin-API-Key")
		if key == "" || key != h.adminAPIKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		return c.Next()
	}
}

type topicRequest struct {
	Name     string   `json:"name"`
	Slug     string   `json:"slug"`
	Keywords []string `json:"keywords"`
	Sources  []string `json:"sources"`
	RSSFeeds []string `json:"rss_feeds"`
	IsActive bool     `json:"is_active"`
}

// CreateTopic handles POST /api/v1/admin/topics
func (h *AdminHandler) CreateTopic(c *fiber.Ctx) error {
	var req topicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.Name == "" || req.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and slug are required",
		})
	}

	keywords := make([]valueobject.TopicKeyword, len(req.Keywords))
	for i, kw := range req.Keywords {
		keywords[i] = valueobject.TopicKeyword(kw)
	}

	sources := make([]valueobject.SourceType, len(req.Sources))
	for i, src := range req.Sources {
		sources[i] = valueobject.SourceType(src)
	}

	topic := &entity.Topic{
		Name:     req.Name,
		Slug:     req.Slug,
		Keywords: keywords,
		Sources:  sources,
		RSSFeeds: req.RSSFeeds,
		IsActive: req.IsActive,
	}

	if err := h.manageTopicsUC.Create(c.Context(), topic); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(topic)
}

// UpdateTopic handles PUT /api/v1/admin/topics/:id
func (h *AdminHandler) UpdateTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	var req topicRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	keywords := make([]valueobject.TopicKeyword, len(req.Keywords))
	for i, kw := range req.Keywords {
		keywords[i] = valueobject.TopicKeyword(kw)
	}

	sources := make([]valueobject.SourceType, len(req.Sources))
	for i, src := range req.Sources {
		sources[i] = valueobject.SourceType(src)
	}

	topic := &entity.Topic{
		ID:       id,
		Name:     req.Name,
		Slug:     req.Slug,
		Keywords: keywords,
		Sources:  sources,
		RSSFeeds: req.RSSFeeds,
		IsActive: req.IsActive,
	}

	if err := h.manageTopicsUC.Update(c.Context(), topic); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(topic)
}

// DeleteTopic handles DELETE /api/v1/admin/topics/:id
func (h *AdminHandler) DeleteTopic(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	if err := h.manageTopicsUC.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "topic deleted successfully",
	})
}

// TriggerFetch handles POST /api/v1/admin/fetch
func (h *AdminHandler) TriggerFetch(c *fiber.Ctx) error {
	fetched, err := h.pipelineUC.ExecuteFetch(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "fetching completed",
		"fetched": fetched,
	})
}

// TriggerProcess handles POST /api/v1/admin/process
func (h *AdminHandler) TriggerProcess(c *fiber.Ctx) error {
	processed, err := h.pipelineUC.ExecuteProcess(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":   "processing completed",
		"processed": processed,
	})
}
