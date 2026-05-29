package usecase

import (
	"context"
	"testing"

	"haberbot/internal/domain/entity"
)

// MockUserRepository implements port.UserRepository for testing.
type MockUserRepository struct {
	users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: make(map[string]*entity.User)}
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	m.users[user.ID] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	return m.users[id], nil
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	for _, u := range m.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateRefreshToken(ctx context.Context, userID string, token string) error {
	return nil
}

func TestUpdateProfileUseCase_BioValidationAndSanitization(t *testing.T) {
	userRepo := NewMockUserRepository()
	uc := NewUpdateProfileUseCase(userRepo)

	ctx := context.Background()
	testUser := &entity.User{
		ID:       "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}
	_ = userRepo.Save(ctx, testUser)

	// Table-driven test cases
	tests := []struct {
		name          string
		bioInput      string
		expectedBio   string
		expectedError bool
	}{
		{
			name:          "Valid standard bio",
			bioInput:      "Merhaba, ben bir teknoloji yazarıyım.",
			expectedBio:   "Merhaba, ben bir teknoloji yazarıyım.",
			expectedError: false,
		},
		{
			name:          "Bio with XSS attempt gets escaped",
			bioInput:      "Ben bir <script>alert('xss')</script> hacker'ım.",
			expectedBio:   "Ben bir &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt; hacker&#39;ım.",
			expectedError: false,
		},
		{
			name:          "Bio exceeding 160 characters gets rejected",
			bioInput:      "Bu biyografi metni tam olarak yuz altmis karakter limitini asmak amaciyla son derece gereksiz detaylarla doldurulmus ve bilerek uzatilmistir. Bu yuzden hata firlatmasi gerekir. Fazla uzun bio!",
			expectedBio:   "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := UpdateProfileRequest{
				UserID: "test-user-id",
				Bio:    &tt.bioInput,
			}

			user, err := uc.Execute(ctx, req)
			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if user.Bio == nil || *user.Bio != tt.expectedBio {
					t.Errorf("Expected Bio %q, got %q", tt.expectedBio, *user.Bio)
				}
			}
		})
	}
}
