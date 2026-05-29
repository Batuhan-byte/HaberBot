package usecase

import (
	"context"
	"errors"
	"fmt"

	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/port"
)

// AddFavoriteUseCase handles adding a user to favorites list.
type AddFavoriteUseCase struct {
	profileRepo port.ProfileRepository
	userRepo    port.UserRepository
}

// NewAddFavoriteUseCase creates a new AddFavoriteUseCase.
func NewAddFavoriteUseCase(profileRepo port.ProfileRepository, userRepo port.UserRepository) *AddFavoriteUseCase {
	return &AddFavoriteUseCase{
		profileRepo: profileRepo,
		userRepo:    userRepo,
	}
}

// Execute adds favorite_user_id to user_id's favorites list.
func (uc *AddFavoriteUseCase) Execute(ctx context.Context, userID, favoriteUserID string) error {
	if userID == "" || favoriteUserID == "" {
		return errors.New("kullanıcı kimlikleri gereklidir")
	}

	if userID == favoriteUserID {
		return errors.New("kendi profilinizi favorilere ekleyemezsiniz")
	}

	// Validate target user exists
	targetUser, err := uc.userRepo.FindByID(ctx, favoriteUserID)
	if err != nil {
		return fmt.Errorf("checking target user existence: %w", err)
	}
	if targetUser == nil {
		return errors.New("favoriye eklenecek kullanıcı bulunamadı")
	}

	// Check if already favorited
	exists, err := uc.profileRepo.IsFavorite(ctx, userID, favoriteUserID)
	if err != nil {
		return fmt.Errorf("checking duplicate favorite: %w", err)
	}
	if exists {
		return errors.New("kullanıcı zaten favorilerinizde ekli")
	}

	favorite := &entity.UserFavorite{
		UserID:         userID,
		FavoriteUserID: favoriteUserID,
	}

	return uc.profileRepo.AddFavorite(ctx, favorite)
}

// RemoveFavoriteUseCase handles removing a user from favorites list.
type RemoveFavoriteUseCase struct {
	profileRepo port.ProfileRepository
}

// NewRemoveFavoriteUseCase creates a new RemoveFavoriteUseCase.
func NewRemoveFavoriteUseCase(profileRepo port.ProfileRepository) *RemoveFavoriteUseCase {
	return &RemoveFavoriteUseCase{profileRepo: profileRepo}
}

// Execute deletes the favorite connection.
func (uc *RemoveFavoriteUseCase) Execute(ctx context.Context, userID, favoriteUserID string) error {
	if userID == "" || favoriteUserID == "" {
		return errors.New("kullanıcı kimlikleri gereklidir")
	}

	return uc.profileRepo.RemoveFavorite(ctx, userID, favoriteUserID)
}

// ListFavoritesUseCase retrieves paginated favorites.
type ListFavoritesUseCase struct {
	profileRepo port.ProfileRepository
}

// NewListFavoritesUseCase creates a new ListFavoritesUseCase.
func NewListFavoritesUseCase(profileRepo port.ProfileRepository) *ListFavoritesUseCase {
	return &ListFavoritesUseCase{profileRepo: profileRepo}
}

// Execute fetches favorited users.
func (uc *ListFavoritesUseCase) Execute(ctx context.Context, userID string, limit, offset int) ([]*entity.User, int, error) {
	if userID == "" {
		return nil, 0, errors.New("kullanıcı kimliği gereklidir")
	}

	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	return uc.profileRepo.ListFavorites(ctx, userID, limit, offset)
}
