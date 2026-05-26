package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// ListArticlesUseCase provides read-only access to articles
// with pagination support.
type ListArticlesUseCase struct {
	articleRepo port.ArticleRepository
}

// NewListArticlesUseCase creates a new ListArticlesUseCase.
func NewListArticlesUseCase(articleRepo port.ArticleRepository) *ListArticlesUseCase {
	return &ListArticlesUseCase{articleRepo: articleRepo}
}

// ListRecent retrieves the most recently fetched articles up to the given limit.
func (uc *ListArticlesUseCase) ListRecent(ctx context.Context, limit int) ([]*entity.Article, error) {
	articles, err := uc.articleRepo.FindRecent(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("listing recent articles: %w", err)
	}
	return articles, nil
}

// ArticlesPage holds a page of articles along with the total count for pagination.
type ArticlesPage struct {
	Articles []*entity.Article
	Total    int
}

// ListByTopic retrieves articles for a specific topic with pagination.
// Returns the articles page and any error encountered.
func (uc *ListArticlesUseCase) ListByTopic(ctx context.Context, topicID string, page, limit int) (*ArticlesPage, error) {
	offset := (page - 1) * limit

	articles, err := uc.articleRepo.FindByTopicID(ctx, topicID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("listing articles by topic: %w", err)
	}

	total, err := uc.articleRepo.CountByTopicID(ctx, topicID)
	if err != nil {
		return nil, fmt.Errorf("counting articles by topic: %w", err)
	}

	return &ArticlesPage{
		Articles: articles,
		Total:    total,
	}, nil
}
