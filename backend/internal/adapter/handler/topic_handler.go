package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"haberbot/internal/usecase"
)

// TopicHandler handles HTTP requests related to topics.
type TopicHandler struct {
	manageTopicsUC *usecase.ManageTopicsUseCase
	listArticlesUC *usecase.ListArticlesUseCase
}

// NewTopicHandler creates a new TopicHandler.
func NewTopicHandler(manageTopicsUC *usecase.ManageTopicsUseCase, listArticlesUC *usecase.ListArticlesUseCase) *TopicHandler {
	return &TopicHandler{
		manageTopicsUC: manageTopicsUC,
		listArticlesUC: listArticlesUC,
	}
}

// GetTopics handles GET /api/v1/topics
func (h *TopicHandler) GetTopics(c *fiber.Ctx) error {
	topics, err := h.manageTopicsUC.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"topics": topics,
	})
}

// GetTopicArticles handles GET /api/v1/topics/:slug/articles
func (h *TopicHandler) GetTopicArticles(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "slug parameter is required",
		})
	}

	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := c.Query("limit", "12")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 12
	}

	topic, err := h.manageTopicsUC.GetBySlug(c.Context(), slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if topic == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "topic not found",
		})
	}

	pageData, err := h.listArticlesUC.ListByTopic(c.Context(), topic.ID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"articles":    pageData.Articles,
		"total":       pageData.Total,
		"page":        page,
		"limit":       limit,
		"total_pages": (pageData.Total + limit - 1) / limit,
	})
}
