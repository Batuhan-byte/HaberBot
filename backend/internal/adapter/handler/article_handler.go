package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"haberbot/internal/usecase"
	"haberbot/internal/domain/valueobject"
)

// ArticleHandler handles HTTP requests related to articles.
type ArticleHandler struct {
	listArticlesUC     *usecase.ListArticlesUseCase
	getArticleUC       *usecase.GetArticleUseCase
	searchArticlesUC   *usecase.SearchArticlesUseCase
	summarizeArticleUC *usecase.SummarizeArticleUseCase
	adminAPIKey        string
}

// NewArticleHandler creates a new ArticleHandler.
func NewArticleHandler(
	listArticlesUC *usecase.ListArticlesUseCase,
	getArticleUC *usecase.GetArticleUseCase,
	searchArticlesUC *usecase.SearchArticlesUseCase,
	summarizeArticleUC *usecase.SummarizeArticleUseCase,
	adminAPIKey string,
) *ArticleHandler {
	return &ArticleHandler{
		listArticlesUC:     listArticlesUC,
		getArticleUC:       getArticleUC,
		searchArticlesUC:   searchArticlesUC,
		summarizeArticleUC: summarizeArticleUC,
		adminAPIKey:        adminAPIKey,
	}
}

// GetRecentArticles handles GET /api/v1/articles
func (h *ArticleHandler) GetRecentArticles(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "12")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 12
	}

	articles, err := h.listArticlesUC.ListRecent(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"articles": articles,
	})
}

// SearchArticles handles GET /api/v1/articles/search?q=...
func (h *ArticleHandler) SearchArticles(c *fiber.Ctx) error {
	query := c.Query("q", "")
	if query == "" {
		// If no query, return recent articles
		limitStr := c.Query("limit", "12")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 12
		}
		articles, err := h.listArticlesUC.ListRecent(c.Context(), limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.JSON(fiber.Map{
			"articles": articles,
		})
	}

	sourceStr := c.Query("source", "")
	var source valueobject.SourceType
	if sourceStr != "" {
		source = valueobject.SourceType(sourceStr)
	}

	limitStr := c.Query("limit", "12")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 12
	}
	offsetStr := c.Query("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	articles, err := h.searchArticlesUC.Search(c.Context(), query, source, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"articles": articles,
		"query":    query,
	})
}

// GetArticleByID handles GET /api/v1/articles/:id
func (h *ArticleHandler) GetArticleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	article, err := h.getArticleUC.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if article == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "article not found",
		})
	}

	// Admin authorization check
	adminKey := c.Get("X-Admin-API-Key")
	isAdmin := adminKey != "" && adminKey == h.adminAPIKey

	// If the article is hidden or unapproved, only allow the admin to view it (otherwise return 404)
	if (!article.IsApproved || article.IsHidden) && !isAdmin {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "article not found",
		})
	}

	// Exclude summary text by default from the article detail response.
	// Users will retrieve it lazily via the GET /articles/:id/summary endpoint.
	responseArticle := *article
	responseArticle.TurkishSummary = ""

	return c.JSON(responseArticle)
}

// SummarizeArticle handles POST /api/v1/articles/:id/summary
func (h *ArticleHandler) SummarizeArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id parameter is required",
		})
	}

	summary, err := h.summarizeArticleUC.Execute(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"summary":    summary,
		"article_id": id,
	})
}
