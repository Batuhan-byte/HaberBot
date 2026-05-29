package port

import (
	"context"
	"haberbot/internal/domain/entity"
)

// CommentRepository defines the persistence contract for Comment entities.
type CommentRepository interface {
	Save(ctx context.Context, comment *entity.Comment) error
	FindByArticleID(ctx context.Context, articleID string) ([]*entity.Comment, error)
}
