package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"haberbot/internal/domain/entity"
)

type mockUserRepo struct {
	saveFunc           func(ctx context.Context, user *entity.User) error
	findByIDFunc       func(ctx context.Context, id string) (*entity.User, error)
	findByUsernameFunc func(ctx context.Context, username string) (*entity.User, error)
	findByEmailFunc    func(ctx context.Context, email string) (*entity.User, error)
	updateTokenFunc    func(ctx context.Context, userID string, token string) error
}

func (m *mockUserRepo) Save(ctx context.Context, user *entity.User) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, user)
	}
	return nil
}

func (m *mockUserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepo) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	if m.findByUsernameFunc != nil {
		return m.findByUsernameFunc(ctx, username)
	}
	return nil, nil
}

func (m *mockUserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	if m.findByEmailFunc != nil {
		return m.findByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *mockUserRepo) UpdateRefreshToken(ctx context.Context, userID string, token string) error {
	if m.updateTokenFunc != nil {
		return m.updateTokenFunc(ctx, userID, token)
	}
	return nil
}

func TestAuthRegisterUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		repo := &mockUserRepo{
			findByUsernameFunc: func(_ context.Context, username string) (*entity.User, error) {
				return nil, nil
			},
			findByEmailFunc: func(_ context.Context, email string) (*entity.User, error) {
				return nil, nil
			},
			saveFunc: func(_ context.Context, user *entity.User) error {
				assert.Equal(t, "testuser", user.Username)
				assert.Equal(t, "test@haberbot.com", user.Email)
				assert.Equal(t, "User", user.Role)
				return nil
			},
		}

		uc := NewAuthRegisterUseCase(repo)
		req := RegisterRequest{
			Username: "testuser",
			Email:    "test@haberbot.com",
			Password: "securepassword123",
		}

		user, err := uc.Execute(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "test@haberbot.com", user.Email)
	})

	t.Run("missing email validation error", func(t *testing.T) {
		repo := &mockUserRepo{}
		uc := NewAuthRegisterUseCase(repo)
		req := RegisterRequest{
			Username: "testuser",
			Password: "securepassword123",
		}

		_, err := uc.Execute(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "email is required", err.Error())
	})

	t.Run("duplicate username error", func(t *testing.T) {
		repo := &mockUserRepo{
			findByUsernameFunc: func(_ context.Context, username string) (*entity.User, error) {
				return &entity.User{Username: "testuser"}, nil
			},
		}

		uc := NewAuthRegisterUseCase(repo)
		req := RegisterRequest{
			Username: "testuser",
			Email:    "test@haberbot.com",
			Password: "securepassword123",
		}

		_, err := uc.Execute(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "username is already taken", err.Error())
	})

	t.Run("duplicate email error", func(t *testing.T) {
		repo := &mockUserRepo{
			findByUsernameFunc: func(_ context.Context, username string) (*entity.User, error) {
				return nil, nil
			},
			findByEmailFunc: func(_ context.Context, email string) (*entity.User, error) {
				return &entity.User{Email: "test@haberbot.com"}, nil
			},
		}

		uc := NewAuthRegisterUseCase(repo)
		req := RegisterRequest{
			Username: "testuser",
			Email:    "test@haberbot.com",
			Password: "securepassword123",
		}

		_, err := uc.Execute(ctx, req)
		assert.Error(t, err)
		assert.Equal(t, "email is already registered", err.Error())
	})
}
