package port

import (
	"context"

	"haberbot/internal/domain/entity"
)

// TopicRepository defines the contract for topic persistence operations.
// Implementations may use PostgreSQL, in-memory storage, or any other backend.
type TopicRepository interface {
	// Save persists a new topic or updates an existing one.
	Save(ctx context.Context, topic *entity.Topic) error

	// FindByID retrieves a single topic by its unique identifier.
	// Returns nil and no error if the topic is not found.
	FindByID(ctx context.Context, id string) (*entity.Topic, error)

	// FindBySlug retrieves a single topic by its URL-friendly slug.
	// Returns nil and no error if the topic is not found.
	FindBySlug(ctx context.Context, slug string) (*entity.Topic, error)

	// FindAll retrieves all topics regardless of their active status.
	FindAll(ctx context.Context) ([]*entity.Topic, error)

	// FindActive retrieves only topics that are currently active.
	FindActive(ctx context.Context) ([]*entity.Topic, error)

	// Delete removes a topic by its unique identifier.
	Delete(ctx context.Context, id string) error
}
