package port

import (
	"context"
	"haberbot/internal/domain/entity"
)

// UserRepository defines the persistence contract for User entities.
type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	UpdateRefreshToken(ctx context.Context, userID string, token string) error
}
