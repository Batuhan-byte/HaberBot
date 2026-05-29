package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// Validation patterns — compiled once at package init for performance.
var (
	// emailRegex is a practical RFC 5322 subset for email format validation.
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	// usernameRegex only allows letters, digits, underscores and hyphens (BUG-014).
	// Prevents stored XSS, log injection, and spoofed admin-looking names.
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
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
	// BUG-014: username must be 3-30 chars and contain only [a-zA-Z0-9_-]
	if len(req.Username) < 3 || len(req.Username) > 30 {
		return nil, errors.New("username must be between 3 and 30 characters long")
	}
	if !usernameRegex.MatchString(req.Username) {
		return nil, errors.New("username may only contain letters, digits, underscores and hyphens")
	}

	// BUG-013: Validate email format with regex (not just empty check)
	if req.Email == "" {
		return nil, errors.New("email is required")
	}
	if !emailRegex.MatchString(req.Email) {
		return nil, errors.New("email address format is invalid")
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
