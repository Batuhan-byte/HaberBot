package port

import (
	"context"
	"haberbot/internal/domain/entity"
)

// ProfileRepository defines the persistence contract for user profiles, favorites, stats, and reports.
type ProfileRepository interface {
	// UserStats operations
	GetStatsByUserID(ctx context.Context, userID string) (*entity.UserStats, error)
	UpsertStats(ctx context.Context, stats *entity.UserStats) error
	RecountStats(ctx context.Context, userID string) error

	// Favorites operations
	AddFavorite(ctx context.Context, favorite *entity.UserFavorite) error
	RemoveFavorite(ctx context.Context, userID, favoriteUserID string) error
	IsFavorite(ctx context.Context, userID, favoriteUserID string) (bool, error)
	ListFavorites(ctx context.Context, userID string, limit, offset int) ([]*entity.User, int, error)

	// Reports operations
	SubmitReport(ctx context.Context, report *entity.ProfileReport) error
	HasRecentReport(ctx context.Context, reporterUserID, reportedUserID string) (bool, error)
}
