package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// GetArticleUseCase retrieves a single article by its ID.
type GetArticleUseCase struct {
	articleRepo port.ArticleRepository
}

// NewGetArticleUseCase creates a new GetArticleUseCase.
func NewGetArticleUseCase(articleRepo port.ArticleRepository) *GetArticleUseCase {
	return &GetArticleUseCase{articleRepo: articleRepo}
}

// Execute retrieves a single article by its unique identifier.
// Returns nil and no error if the article does not exist.
func (uc *GetArticleUseCase) Execute(ctx context.Context, id string) (*entity.Article, error) {
	article, err := uc.articleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting article %s: %w", id, err)
	}
	return article, nil
}
