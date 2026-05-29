package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"haberbot/internal/usecase"
)

// CommentHandler handles HTTP requests for user comments.
type CommentHandler struct {
	createCommentUC *usecase.CreateCommentUseCase
	listCommentsUC  *usecase.ListCommentsUseCase
}

// NewCommentHandler creates a new CommentHandler.
func NewCommentHandler(
	createCommentUC *usecase.CreateCommentUseCase,
	listCommentsUC *usecase.ListCommentsUseCase,
) *CommentHandler {
	return &CommentHandler{
		createCommentUC: createCommentUC,
		listCommentsUC:  listCommentsUC,
	}
}

// CreateComment handles POST /api/comments
func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	var req usecase.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Extract authenticated User ID from JWT claims
	userClaims, ok := c.Locals("user").(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	userID, ok := userClaims["sub"].(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized: missing user identifier",
		})
	}

	req.UserID = userID

	comment, err := h.createCommentUC.Execute(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Fetch username for response
	username, _ := userClaims["username"].(string)
	comment.Username = username

	return c.Status(fiber.StatusCreated).JSON(comment)
}

// ListComments handles GET /api/comments
func (h *CommentHandler) ListComments(c *fiber.Ctx) error {
	articleID := c.Query("article_id")
	if articleID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "article_id query parameter is required",
		})
	}

	comments, err := h.listCommentsUC.Execute(c.Context(), articleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"comments": comments,
	})
}
