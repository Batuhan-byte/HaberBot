package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// ManageArticlesUseCase handles administrative actions on articles.
type ManageArticlesUseCase struct {
	articleRepo port.ArticleRepository
}

// NewManageArticlesUseCase creates a new ManageArticlesUseCase.
func NewManageArticlesUseCase(articleRepo port.ArticleRepository) *ManageArticlesUseCase {
	return &ManageArticlesUseCase{articleRepo: articleRepo}
}

// ListArticlesAdmin retrieves all articles for the admin dashboard with total counts and pagination.
func (uc *ManageArticlesUseCase) ListArticlesAdmin(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, int, error) {
	articles, err := uc.articleRepo.FindAllAdmin(ctx, topicID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("listing admin articles: %w", err)
	}

	total, err := uc.articleRepo.CountAllAdmin(ctx, topicID)
	if err != nil {
		return nil, 0, fmt.Errorf("counting admin articles: %w", err)
	}

	return articles, total, nil
}

// Approve sets the approval status of an article.
func (uc *ManageArticlesUseCase) Approve(ctx context.Context, id string, approved bool) error {
	if err := uc.articleRepo.UpdateApprovalStatus(ctx, id, approved); err != nil {
		return fmt.Errorf("updating article approval for %s: %w", id, err)
	}
	return nil
}

// Hide sets the hiding (soft-delete) status of an article.
func (uc *ManageArticlesUseCase) Hide(ctx context.Context, id string, hidden bool) error {
	if err := uc.articleRepo.UpdateHidingStatus(ctx, id, hidden); err != nil {
		return fmt.Errorf("updating article hiding status for %s: %w", id, err)
	}
	return nil
}

// Delete hard-deletes an article from the database.
func (uc *ManageArticlesUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.articleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting article %s: %w", id, err)
	}
	return nil
}
