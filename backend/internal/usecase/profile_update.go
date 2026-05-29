package usecase

import (
	"context"
	"errors"
	"fmt"
	"html"
	"strings"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// UpdateProfileRequest defines the input payload for updating profile details.
type UpdateProfileRequest struct {
	UserID    string
	Bio       *string
	AvatarURL *string
}

// UpdateProfileUseCase coordinates updating own profile bio and avatar url.
type UpdateProfileUseCase struct {
	userRepo port.UserRepository
}

// NewUpdateProfileUseCase creates a new UpdateProfileUseCase.
func NewUpdateProfileUseCase(userRepo port.UserRepository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo}
}

// Execute performs validation and updates the user's profile information.
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, req UpdateProfileRequest) (*entity.User, error) {
	if req.UserID == "" {
		return nil, errors.New("user ID is required")
	}

	user, err := uc.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 1. Bio validation and sanitization
	if req.Bio != nil {
		bioText := *req.Bio
		// Truncate/validate bio length of 160 characters
		if len([]rune(bioText)) > 160 {
			return nil, errors.New("biyografi en fazla 160 karakter olabilir")
		}
		// Sanitize bio to prevent XSS: strip out tags or escape
		sanitizedBio := html.EscapeString(strings.TrimSpace(bioText))
		user.Bio = &sanitizedBio
	}

	// 2. AvatarURL update
	if req.AvatarURL != nil {
		avatarURL := *req.AvatarURL
		if avatarURL != "" {
			// Basic link check or local storage path validation
			if !strings.HasPrefix(avatarURL, "/public/avatars/") && !strings.HasPrefix(avatarURL, "http") {
				return nil, errors.New("geçersiz avatar URL biçimi")
			}
		}
		user.AvatarURL = &avatarURL
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("saving updated profile: %w", err)
	}

	return user, nil
}
