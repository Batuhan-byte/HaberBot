package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/port"
)

type SummarizeArticleUseCase struct {
	aiProcessor port.AIProcessor
	articleRepo port.ArticleRepository
}

func NewSummarizeArticleUseCase(aiProcessor port.AIProcessor, articleRepo port.ArticleRepository) *SummarizeArticleUseCase {
	return &SummarizeArticleUseCase{
		aiProcessor: aiProcessor,
		articleRepo: articleRepo,
	}
}

func (uc *SummarizeArticleUseCase) Execute(ctx context.Context, articleID string) (string, error) {
	article, err := uc.articleRepo.FindByID(ctx, articleID)
	if err != nil {
		return "", fmt.Errorf("finding article: %w", err)
	}
	if article == nil {
		return "", fmt.Errorf("article not found")
	}

	// If already summarized, return immediately
	if article.TurkishSummary != "" {
		return article.TurkishSummary, nil
	}

	contentToSummarize := article.TurkishContent
	if contentToSummarize == "" {
		contentToSummarize = article.OriginalContent
	}

	summary, err := uc.aiProcessor.Summarize(ctx, contentToSummarize)
	if err != nil {
		return "", fmt.Errorf("ai summarize: %w", err)
	}

	article.TurkishSummary = summary
	if err := uc.articleRepo.Save(ctx, article); err != nil {
		return "", fmt.Errorf("saving summary: %w", err)
	}

	return summary, nil
}
