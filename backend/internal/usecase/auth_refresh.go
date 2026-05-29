package usecase

import (
	"context"
	"errors"
	"fmt"

	"haberbot/internal/domain/port"
	"haberbot/internal/infrastructure/auth"
)

// AuthRefreshUseCase validates the HttpOnly refresh token and generates a rotated pair.
type AuthRefreshUseCase struct {
	userRepo            port.UserRepository
	jwtSecret           string
	jwtAccessTTLMinutes int
	jwtRefreshTTLDays   int
}

// NewAuthRefreshUseCase creates a new AuthRefreshUseCase.
func NewAuthRefreshUseCase(
	userRepo port.UserRepository,
	jwtSecret string,
	jwtAccessTTLMinutes int,
	jwtRefreshTTLDays int,
) *AuthRefreshUseCase {
	return &AuthRefreshUseCase{
		userRepo:            userRepo,
		jwtSecret:           jwtSecret,
		jwtAccessTTLMinutes: jwtAccessTTLMinutes,
		jwtRefreshTTLDays:   jwtRefreshTTLDays,
	}
}

// Execute performs refresh token validation, reuse detection, and rotation.
func (uc *AuthRefreshUseCase) Execute(ctx context.Context, refreshTokenStr string) (*LoginResponse, error) {
	if refreshTokenStr == "" {
		return nil, errors.New("refresh token is required")
	}

	// Validate refresh token
	claims, err := auth.ValidateToken(refreshTokenStr, uc.jwtSecret)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return nil, errors.New("invalid token claims")
	}

	// Fetch user to verify stored token
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Token rotation and reuse detection:
	// If the provided refresh token does not match the stored token,
	// it means the token was already used or compromised!
	if user.RefreshToken != refreshTokenStr {
		// Invalidate and delete refresh token to protect the account
		_ = uc.userRepo.UpdateRefreshToken(ctx, user.ID, "")
		return nil, errors.New("refresh token compromise detected; logging out")
	}

	// Generate new Access Token
	newAccessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, uc.jwtSecret, uc.jwtAccessTTLMinutes)
	if err != nil {
		return nil, fmt.Errorf("generating rotated access token: %w", err)
	}

	// Generate new Refresh Token (Rotation)
	newRefreshToken, err := auth.GenerateRefreshToken(user.ID, uc.jwtSecret, uc.jwtRefreshTTLDays)
	if err != nil {
		return nil, fmt.Errorf("generating rotated refresh token: %w", err)
	}

	// Update refresh token in DB
	if err := uc.userRepo.UpdateRefreshToken(ctx, user.ID, newRefreshToken); err != nil {
		return nil, fmt.Errorf("updating rotated refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}
