package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// ManageTopicsUseCase handles CRUD operations for topics.
type ManageTopicsUseCase struct {
	topicRepo port.TopicRepository
}

// NewManageTopicsUseCase creates a new ManageTopicsUseCase.
func NewManageTopicsUseCase(topicRepo port.TopicRepository) *ManageTopicsUseCase {
	return &ManageTopicsUseCase{topicRepo: topicRepo}
}

// Create persists a new topic.
func (uc *ManageTopicsUseCase) Create(ctx context.Context, topic *entity.Topic) error {
	if err := uc.topicRepo.Save(ctx, topic); err != nil {
		return fmt.Errorf("creating topic: %w", err)
	}
	return nil
}

// GetAll retrieves all topics regardless of active status.
func (uc *ManageTopicsUseCase) GetAll(ctx context.Context) ([]*entity.Topic, error) {
	topics, err := uc.topicRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing all topics: %w", err)
	}
	return topics, nil
}

// GetBySlug retrieves a single topic by its URL-friendly slug.
func (uc *ManageTopicsUseCase) GetBySlug(ctx context.Context, slug string) (*entity.Topic, error) {
	topic, err := uc.topicRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("finding topic by slug %s: %w", slug, err)
	}
	return topic, nil
}

// Update persists changes to an existing topic.
func (uc *ManageTopicsUseCase) Update(ctx context.Context, topic *entity.Topic) error {
	if err := uc.topicRepo.Save(ctx, topic); err != nil {
		return fmt.Errorf("updating topic: %w", err)
	}
	return nil
}

// Delete removes a topic by its unique identifier.
func (uc *ManageTopicsUseCase) Delete(ctx context.Context, id string) error {
	if err := uc.topicRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting topic %s: %w", id, err)
	}
	return nil
}
