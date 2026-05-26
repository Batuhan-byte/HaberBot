package usecase

import (
	"context"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"haberbot/internal/domain/valueobject"
)

// SearchArticlesUseCase provides search functionality for articles.
type SearchArticlesUseCase struct {
	articleRepo port.ArticleRepository
}

// NewSearchArticlesUseCase creates a new SearchArticlesUseCase.
func NewSearchArticlesUseCase(articleRepo port.ArticleRepository) *SearchArticlesUseCase {
	return &SearchArticlesUseCase{articleRepo: articleRepo}
}

// Search searches for articles matching the query.
// Returns the search results and any error encountered.
func (uc *SearchArticlesUseCase) Search(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
	if query == "" && source == "" {
		// If both query and source are empty, return recent articles
		return uc.articleRepo.FindRecent(ctx, limit)
	}

	articles, err := uc.articleRepo.Search(ctx, query, source, limit, offset)
	if err != nil {
		return nil, err
	}
	return articles, nil
}