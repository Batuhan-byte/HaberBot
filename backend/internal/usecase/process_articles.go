package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"haberbot/internal/domain/port"
)

// geminiRequestDelay is the pause between AI processing requests to respect
// the Gemini API free-tier rate limit of 5 RPM.
const geminiRequestDelay = 12 * time.Second

// ProcessArticlesUseCase sends unprocessed articles through the AI pipeline
// for Turkish translation and summarization.
type ProcessArticlesUseCase struct {
	aiProcessor port.AIProcessor
	articleRepo port.ArticleRepository
}

// NewProcessArticlesUseCase creates a new ProcessArticlesUseCase.
func NewProcessArticlesUseCase(aiProcessor port.AIProcessor, articleRepo port.ArticleRepository) *ProcessArticlesUseCase {
	return &ProcessArticlesUseCase{
		aiProcessor: aiProcessor,
		articleRepo: articleRepo,
	}
}

// Execute retrieves up to batchSize unprocessed articles, sends each to the
// AI processor, and updates them with Turkish translations. Returns the count
// of successfully processed articles.
func (uc *ProcessArticlesUseCase) Execute(ctx context.Context, batchSize int) (int, error) {
	articles, err := uc.articleRepo.FindUnprocessed(ctx, batchSize)
	if err != nil {
		return 0, fmt.Errorf("finding unprocessed articles: %w", err)
	}

	var processed int
	for _, article := range articles {
		if err := ctx.Err(); err != nil {
			return processed, fmt.Errorf("context cancelled: %w", err)
		}

		// Send content to AI for full translation
		turkishTitle, turkishContent, err := uc.aiProcessor.Translate(ctx, article.Title, article.OriginalContent)
		if err != nil {
			slog.Error("failed to process article with AI",
				"article_id", article.ID,
				"error", err,
			)
			// Sleep even on error to respect rate limits before the next call
			time.Sleep(geminiRequestDelay)
			continue
		}

		now := time.Now()
		article.TurkishTitle = turkishTitle
		article.TurkishContent = turkishContent
		article.ProcessedAt = &now

		if err := uc.articleRepo.Save(ctx, article); err != nil {
			return processed, fmt.Errorf("saving processed article %s: %w", article.ID, err)
		}

		processed++
		slog.Info("processed article", "id", article.ID, "title", turkishTitle)

		// Respect Gemini rate limit between requests.
		time.Sleep(geminiRequestDelay)
	}

	return processed, nil
}
