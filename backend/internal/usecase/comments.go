package usecase

import (
	"context"
	"errors"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// CreateCommentRequest defines the input payload for creating a new comment.
type CreateCommentRequest struct {
	ArticleID string `json:"article_id"`
	UserID    string `json:"-"`
	Content   string `json:"content"`
}

// CreateCommentUseCase handles creation and persistence of comments.
type CreateCommentUseCase struct {
	commentRepo port.CommentRepository
	articleRepo port.ArticleRepository
}

// NewCreateCommentUseCase creates a new CreateCommentUseCase.
func NewCreateCommentUseCase(commentRepo port.CommentRepository, articleRepo port.ArticleRepository) *CreateCommentUseCase {
	return &CreateCommentUseCase{
		commentRepo: commentRepo,
		articleRepo: articleRepo,
	}
}

// Execute persists a new comment on an article.
func (uc *CreateCommentUseCase) Execute(ctx context.Context, req CreateCommentRequest) (*entity.Comment, error) {
	if req.ArticleID == "" {
		return nil, errors.New("article ID is required")
	}
	if req.Content == "" {
		return nil, errors.New("comment content cannot be empty")
	}

	// Validate article exists
	article, err := uc.articleRepo.FindByID(ctx, req.ArticleID)
	if err != nil {
		return nil, fmt.Errorf("checking article existence: %w", err)
	}
	if article == nil {
		return nil, errors.New("article not found")
	}

	comment := &entity.Comment{
		ArticleID: req.ArticleID,
		UserID:    req.UserID,
		Content:   req.Content,
	}

	if err := uc.commentRepo.Save(ctx, comment); err != nil {
		return nil, fmt.Errorf("saving comment: %w", err)
	}

	return comment, nil
}

// ListCommentsUseCase retrieves comments for a given article.
type ListCommentsUseCase struct {
	commentRepo port.CommentRepository
}

// NewListCommentsUseCase creates a new ListCommentsUseCase.
func NewListCommentsUseCase(commentRepo port.CommentRepository) *ListCommentsUseCase {
	return &ListCommentsUseCase{commentRepo: commentRepo}
}

// Execute fetches and returns all comments for the given article ID.
func (uc *ListCommentsUseCase) Execute(ctx context.Context, articleID string) ([]*entity.Comment, error) {
	if articleID == "" {
		return nil, errors.New("article ID is required")
	}

	comments, err := uc.commentRepo.FindByArticleID(ctx, articleID)
	if err != nil {
		return nil, fmt.Errorf("listing comments: %w", err)
	}

	return comments, nil
}
