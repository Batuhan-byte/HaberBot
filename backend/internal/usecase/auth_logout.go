package usecase

import (
	"context"
	"fmt"

	"haberbot/internal/domain/port"
)

// AuthLogoutUseCase invalidates the user's session by clearing the refresh token in the repository.
type AuthLogoutUseCase struct {
	userRepo port.UserRepository
}

// NewAuthLogoutUseCase creates a new AuthLogoutUseCase.
func NewAuthLogoutUseCase(userRepo port.UserRepository) *AuthLogoutUseCase {
	return &AuthLogoutUseCase{userRepo: userRepo}
}

// Execute clears the refresh token associated with the user ID.
func (uc *AuthLogoutUseCase) Execute(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}

	if err := uc.userRepo.UpdateRefreshToken(ctx, userID, ""); err != nil {
		return fmt.Errorf("clearing session refresh token: %w", err)
	}

	return nil
}
