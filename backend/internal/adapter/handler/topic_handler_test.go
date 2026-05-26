package handler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"haberbot/internal/domain/entity"
	"haberbot/internal/domain/valueobject"
	"haberbot/internal/usecase"
)

type mockTopicRepo struct {
	saveFunc       func(ctx context.Context, topic *entity.Topic) error
	findAllFunc    func(ctx context.Context) ([]*entity.Topic, error)
	findBySlugFunc func(ctx context.Context, slug string) (*entity.Topic, error)
	deleteFunc     func(ctx context.Context, id string) error
	findByIDFunc   func(ctx context.Context, id string) (*entity.Topic, error)
	findActiveFunc func(ctx context.Context) ([]*entity.Topic, error)
}

func (m *mockTopicRepo) Save(ctx context.Context, topic *entity.Topic) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, topic)
	}
	return nil
}
func (m *mockTopicRepo) FindByID(ctx context.Context, id string) (*entity.Topic, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockTopicRepo) FindBySlug(ctx context.Context, slug string) (*entity.Topic, error) {
	if m.findBySlugFunc != nil {
		return m.findBySlugFunc(ctx, slug)
	}
	return nil, nil
}
func (m *mockTopicRepo) FindAll(ctx context.Context) ([]*entity.Topic, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}
func (m *mockTopicRepo) FindActive(ctx context.Context) ([]*entity.Topic, error) {
	if m.findActiveFunc != nil {
		return m.findActiveFunc(ctx)
	}
	return nil, nil
}
func (m *mockTopicRepo) Delete(ctx context.Context, id string) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, id)
	}
	return nil
}

type mockArticleRepo struct {
	findRecentFunc      func(ctx context.Context, limit int) ([]*entity.Article, error)
	findByIDFunc        func(ctx context.Context, id string) (*entity.Article, error)
	findByTopicIDFunc   func(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error)
	findUnprocessedFunc func(ctx context.Context, limit int) ([]*entity.Article, error)
	existsByURLFunc     func(ctx context.Context, url string) (bool, error)
	countByTopicIDFunc  func(ctx context.Context, topicID string) (int, error)
	saveFunc            func(ctx context.Context, article *entity.Article) error
	searchFunc          func(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error)
}

func (m *mockArticleRepo) FindRecent(ctx context.Context, limit int) ([]*entity.Article, error) {
	if m.findRecentFunc != nil {
		return m.findRecentFunc(ctx, limit)
	}
	return nil, nil
}
func (m *mockArticleRepo) FindByID(ctx context.Context, id string) (*entity.Article, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
func (m *mockArticleRepo) FindByTopicID(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
	if m.findByTopicIDFunc != nil {
		return m.findByTopicIDFunc(ctx, topicID, limit, offset)
	}
	return nil, nil
}
func (m *mockArticleRepo) FindUnprocessed(ctx context.Context, limit int) ([]*entity.Article, error) {
	if m.findUnprocessedFunc != nil {
		return m.findUnprocessedFunc(ctx, limit)
	}
	return nil, nil
}
func (m *mockArticleRepo) ExistsByURL(ctx context.Context, url string) (bool, error) {
	if m.existsByURLFunc != nil {
		return m.existsByURLFunc(ctx, url)
	}
	return false, nil
}
func (m *mockArticleRepo) CountByTopicID(ctx context.Context, topicID string) (int, error) {
	if m.countByTopicIDFunc != nil {
		return m.countByTopicIDFunc(ctx, topicID)
	}
	return 0, nil
}
func (m *mockArticleRepo) Save(ctx context.Context, article *entity.Article) error {
	if m.saveFunc != nil {
		return m.saveFunc(ctx, article)
	}
	return nil
}
func (m *mockArticleRepo) Search(ctx context.Context, query string, source valueobject.SourceType, limit, offset int) ([]*entity.Article, error) {
	if m.searchFunc != nil {
		return m.searchFunc(ctx, query, source, limit, offset)
	}
	return nil, nil
}

type mockContentFetcher struct {
	fetchByKeywordsFunc func(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error)
	sourceTypeFunc      func() valueobject.SourceType
}

func (m *mockContentFetcher) FetchByKeywords(ctx context.Context, keywords []valueobject.TopicKeyword) ([]*entity.Article, error) {
	if m.fetchByKeywordsFunc != nil {
		return m.fetchByKeywordsFunc(ctx, keywords)
	}
	return nil, nil
}
func (m *mockContentFetcher) SourceType() valueobject.SourceType {
	if m.sourceTypeFunc != nil {
		return m.sourceTypeFunc()
	}
	return ""
}

type mockAIProcessor struct {
	translateFunc func(ctx context.Context, title, content string) (string, string, error)
	summarizeFunc func(ctx context.Context, content string) (string, error)
}

func (m *mockAIProcessor) Translate(ctx context.Context, title, content string) (string, string, error) {
	if m.translateFunc != nil {
		return m.translateFunc(ctx, title, content)
	}
	return "", "", nil
}
func (m *mockAIProcessor) Summarize(ctx context.Context, content string) (string, error) {
	if m.summarizeFunc != nil {
		return m.summarizeFunc(ctx, content)
	}
	return "", nil
}

func TestGetTopicArticles_EdgeCases(t *testing.T) {
	app := fiber.New()

	mockTopic := &mockTopicRepo{}
	mockArticle := &mockArticleRepo{}

	manageTopicsUC := usecase.NewManageTopicsUseCase(mockTopic)
	listArticlesUC := usecase.NewListArticlesUseCase(mockArticle)
	h := NewTopicHandler(manageTopicsUC, listArticlesUC)

	app.Get("/api/v1/topics/:slug/articles", h.GetTopicArticles)

	tests := []struct {
		name           string
		slug           string
		query          string
		setupMocks     func()
		expectedStatus int
		verifyResponse func(t *testing.T, body []byte)
	}{
		{
			name: "Topic Not Found - 404",
			slug: "non-existent",
			setupMocks: func() {
				mockTopic.findBySlugFunc = func(ctx context.Context, slug string) (*entity.Topic, error) {
					return nil, nil
				}
			},
			expectedStatus: http.StatusNotFound,
			verifyResponse: func(t *testing.T, body []byte) {
				var resp map[string]string
				json.Unmarshal(body, &resp)
				if resp["error"] != "topic not found" {
					t.Errorf("expected error 'topic not found', got %v", resp["error"])
				}
			},
		},
		{
			name: "Database Error Getting Topic - 500",
			slug: "error-slug",
			setupMocks: func() {
				mockTopic.findBySlugFunc = func(ctx context.Context, slug string) (*entity.Topic, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:  "Invalid Page and Limit Query Params - Fallback Sanitization",
			slug:  "valid-topic",
			query: "?page=abc&limit=-5",
			setupMocks: func() {
				mockTopic.findBySlugFunc = func(ctx context.Context, slug string) (*entity.Topic, error) {
					return &entity.Topic{ID: "123", Name: "Valid Topic", Slug: "valid-topic"}, nil
				}
				mockArticle.findByTopicIDFunc = func(ctx context.Context, topicID string, limit, offset int) ([]*entity.Article, error) {
					if limit != 12 || offset != 0 {
						t.Errorf("expected fallback limit 12, got %d; offset 0, got %d", limit, offset)
					}
					return []*entity.Article{}, nil
				}
				mockArticle.countByTopicIDFunc = func(ctx context.Context, topicID string) (int, error) {
					return 0, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/api/v1/topics/"+tt.slug+"/articles"+tt.query, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("failed to run test: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			if tt.verifyResponse != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					t.Fatalf("failed to read response body: %v", err)
				}
				tt.verifyResponse(t, body)
			}
		})
	}
}
