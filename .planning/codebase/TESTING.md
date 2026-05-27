# Testing Patterns

**Analysis Date:** 2026-05-27

## Test Framework

### Go Backend:
- **Runner:** Native Go `testing` standard library.
- **Assertion Library:** `github.com/stretchr/testify/assert` and `github.com/stretchr/testify/require`.
- **Run Commands:**
  ```bash
  go test -v ./...                 # Run all backend tests
  go test -v ./internal/domain/... # Run domain layer tests only
  go test -v -run TestArticleFetch # Run a specific test suite
  ```

### React Frontend:
- **Runner:** Vitest v4 (`vitest`)
- **Config:** `frontend/vite.config.js` or `frontend/package.json` config.
- **Assertion Library:** Vitest built-in `expect`.
- **Testing Library:** `@testing-library/react` and `@testing-library/jest-dom`.
- **Run Commands:**
  ```bash
  npm test                        # Run all frontend tests (Vitest run mode)
  npm run test:watch              # Watch mode
  ```

---

## Test File Organization

### Go Backend:
- Collocated with the implementation. For every source file `module_name.go`, the corresponding test file `module_name_test.go` sits in the exact same directory.
- Example layout:
  ```
  backend/internal/domain/entity/
    ├── article.go
    ├── article_test.go
    ├── topic.go
    └── topic_test.go
  ```

### React Frontend:
- Located under a dedicated `frontend/src/test/` directory, or collocated alongside the component folder structure.

---

## Test Structure

### Go Backend (Table-Driven Tests):
Almost all critical parsers, business logic validators, and handler layers in Go are verified using standard **table-driven tests** to cover success and error execution branches:

```go
func TestTopicValidation(t *testing.T) {
    tests := []struct {
        name    string
        title   string
        slug    string
        wantErr bool
    }{
        {
            name:    "valid topic",
            title:   "Teknoloji",
            slug:    "teknoloji",
            wantErr: false,
        },
        {
            name:    "empty title",
            title:   "",
            slug:    "teknoloji",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            topic, err := entity.NewTopic(uuid.NewString(), tt.title, tt.slug)
            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, topic)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, topic)
            }
        })
    }
}
```

---

## Mocking

### Go Backend Mocking Patterns:
- Mocking is strictly interface-based. External ports (e.g. `port.ArticleRepository`, `port.AIProcessor`) are mocked during Usecase layer unit testing.
- Mocks are manually coded or mock structures are declared directly inside the test files to keep testing boundaries completely isolated from Postgres databases or external Google Gemini endpoints.

Example mock implementation in tests:
```go
type mockArticleRepository struct {
    articles []*entity.Article
    err      error
}

func (m *mockArticleRepository) Save(ctx context.Context, a *entity.Article) error {
    if m.err != nil {
        return m.err
    }
    m.articles = append(m.articles, a)
    return nil
}
```

---

## What to Mock

**Always Mock:**
- Network / HTTP APIs (Google Gemini translation requests, OpenAI API triggers).
- Raw RSS Feed fetches (use static XML fixtures).
- PostgreSQL database connections (mock the database repository ports).

**Do NOT Mock:**
- Clean Architecture Domain Entities (`entity.Article`, `entity.Topic`).
- Pure parsing helper utilities (e.g., `go-readability` extractions, title cleaning functions).

---

*Testing analysis: 2026-05-27*
*Update when test patterns change*
