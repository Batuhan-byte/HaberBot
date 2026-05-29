package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"haberbot/internal/domain/port"
)

// ProfileResponse represents the public profile output.
// SECURITY: Email and Role fields are intentionally excluded to prevent data leakage (BUG-005).
type ProfileResponse struct {
	ID                string    `json:"id"`
	Username          string    `json:"username"`
	Bio               *string   `json:"bio"`
	AvatarURL         *string   `json:"avatar_url"`
	JoinDate          time.Time `json:"join_date"`
	ArticlesPublished int       `json:"articles_published"`
	IsFavorited       bool      `json:"is_favorited"`
}

// GetProfileUseCase encapsulates logic to retrieve public user profile details.
type GetProfileUseCase struct {
	userRepo    port.UserRepository
	profileRepo port.ProfileRepository
}

// NewGetProfileUseCase creates a new GetProfileUseCase.
func NewGetProfileUseCase(userRepo port.UserRepository, profileRepo port.ProfileRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepo:    userRepo,
		profileRepo: profileRepo,
	}
}

// Execute retrieves profile data for the target username and checks if the requesting user has favorited them.
func (uc *GetProfileUseCase) Execute(ctx context.Context, username string, currentUserID string) (*ProfileResponse, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	user, err := uc.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("finding user by username: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	stats, err := uc.profileRepo.GetStatsByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("getting user stats: %w", err)
	}

	isFavorited := false
	if currentUserID != "" && currentUserID != user.ID {
		fav, err := uc.profileRepo.IsFavorite(ctx, currentUserID, user.ID)
		if err == nil {
			isFavorited = fav
		}
	}

	return &ProfileResponse{
		ID:                user.ID,
		Username:          user.Username,
		Bio:               user.Bio,
		AvatarURL:         user.AvatarURL,
		JoinDate:          user.JoinDate,
		ArticlesPublished: stats.ArticlesPublished,
		IsFavorited:       isFavorited,
	}, nil
}
