package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
	"haberbot/internal/infrastructure/auth"
)

// LoginRequest defines the input payload for user login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains the authentication tokens and basic user info.
type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"-"` // Passed via cookie in HTTP handler
	User         *entity.User `json:"user"`
}

// AuthLoginUseCase orchestrates user verification, bcrypt validation, and JWT generation.
type AuthLoginUseCase struct {
	userRepo            port.UserRepository
	jwtSecret           string
	jwtAccessTTLMinutes int
	jwtRefreshTTLDays   int
}

// NewAuthLoginUseCase creates a new AuthLoginUseCase.
func NewAuthLoginUseCase(
	userRepo port.UserRepository,
	jwtSecret string,
	jwtAccessTTLMinutes int,
	jwtRefreshTTLDays int,
) *AuthLoginUseCase {
	return &AuthLoginUseCase{
		userRepo:            userRepo,
		jwtSecret:           jwtSecret,
		jwtAccessTTLMinutes: jwtAccessTTLMinutes,
		jwtRefreshTTLDays:   jwtRefreshTTLDays,
	}
}

// Execute validates credentials and returns new JWT access and refresh tokens.
func (uc *AuthLoginUseCase) Execute(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}

	user, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JWT Access Token
	accessToken, err := auth.GenerateAccessToken(user.ID, user.Username, user.Role, uc.jwtSecret, uc.jwtAccessTTLMinutes)
	if err != nil {
		return nil, fmt.Errorf("generating access token: %w", err)
	}

	// Generate JWT Refresh Token
	refreshToken, err := auth.GenerateRefreshToken(user.ID, uc.jwtSecret, uc.jwtRefreshTTLDays)
	if err != nil {
		return nil, fmt.Errorf("generating refresh token: %w", err)
	}

	// Update refresh token in repository (rotation)
	if err := uc.userRepo.UpdateRefreshToken(ctx, user.ID, refreshToken); err != nil {
		return nil, fmt.Errorf("updating refresh token: %w", err)
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
