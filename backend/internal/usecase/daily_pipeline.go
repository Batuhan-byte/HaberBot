package usecase

import (
	"context"
	"fmt"
	"log/slog"
)

// DailyPipelineUseCase orchestrates the daily content pipeline:
// 1. Fetch articles for all active topics
// 2. Process all unprocessed articles through AI
type DailyPipelineUseCase struct {
	fetchUC   *FetchArticlesUseCase
	processUC *ProcessArticlesUseCase
	topicUC   *ManageTopicsUseCase
}

// defaultProcessBatchSize is the number of articles processed per pipeline run.
const defaultProcessBatchSize = 3

// NewDailyPipelineUseCase creates a new DailyPipelineUseCase.
func NewDailyPipelineUseCase(
	fetchUC *FetchArticlesUseCase,
	processUC *ProcessArticlesUseCase,
	topicUC *ManageTopicsUseCase,
) *DailyPipelineUseCase {
	return &DailyPipelineUseCase{
		fetchUC:   fetchUC,
		processUC: processUC,
		topicUC:   topicUC,
	}
}

// ExecuteFetch fetches articles for all active topics. Returns the total
// number of newly fetched articles across all topics.
func (uc *DailyPipelineUseCase) ExecuteFetch(ctx context.Context) (int, error) {
	topics, err := uc.topicUC.GetAll(ctx)
	if err != nil {
		return 0, fmt.Errorf("getting active topics: %w", err)
	}

	var totalFetched int
	for _, topic := range topics {
		if !topic.IsActive {
			continue
		}

		fetched, err := uc.fetchUC.Execute(ctx, topic)
		if err != nil {
			slog.Error("failed to fetch for topic",
				"topic", topic.Slug, "error", err,
			)
			continue
		}

		slog.Info("fetched articles for topic",
			"topic", topic.Slug, "count", fetched,
		)
		totalFetched += fetched
	}

	// Perform periodic trim cleanup for all active topics
	for _, topic := range topics {
		if !topic.IsActive {
			continue
		}
		if err := uc.fetchUC.TrimPending(ctx, topic.ID); err != nil {
			slog.Error("failed periodic trim for topic", "topic", topic.Slug, "error", err)
		}
	}

	// Perform periodic trim cleanup for NULL topic
	if err := uc.fetchUC.TrimPending(ctx, ""); err != nil {
		slog.Error("failed periodic trim for NULL topic", "error", err)
	}

	return totalFetched, nil
}

// ExecuteProcess processes all unprocessed articles through the AI pipeline.
// Returns the count of successfully processed articles.
func (uc *DailyPipelineUseCase) ExecuteProcess(ctx context.Context) (int, error) {
	processed, err := uc.processUC.Execute(ctx, defaultProcessBatchSize)
	if err != nil {
		return processed, fmt.Errorf("processing articles: %w", err)
	}

	slog.Info("daily processing complete", "processed", processed)
	return processed, nil
}
