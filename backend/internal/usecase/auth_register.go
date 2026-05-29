package usecase

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// RegisterRequest defines the input payload for registering a new user.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthRegisterUseCase handles registering new users with password validation and hashing.
type AuthRegisterUseCase struct {
	userRepo port.UserRepository
}

// NewAuthRegisterUseCase creates a new AuthRegisterUseCase.
func NewAuthRegisterUseCase(userRepo port.UserRepository) *AuthRegisterUseCase {
	return &AuthRegisterUseCase{userRepo: userRepo}
}

// Execute registers a new user with standard User role.
func (uc *AuthRegisterUseCase) Execute(ctx context.Context, req RegisterRequest) (*entity.User, error) {
	if len(req.Username) < 3 {
		return nil, errors.New("username must be at least 3 characters long")
	}
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	}

	// Check if user already exists by username
	existing, err := uc.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("checking username uniqueness: %w", err)
	}
	if existing != nil {
		return nil, errors.New("username is already taken")
	}

	// Check if user already exists by email
	existingEmail, err := uc.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("checking email uniqueness: %w", err)
	}
	if existingEmail != nil {
		return nil, errors.New("email is already registered")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashed),
		Role:         "User",
	}

	if err := uc.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("saving registered user: %w", err)
	}

	return user, nil
}
