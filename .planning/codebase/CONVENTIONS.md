# Coding Conventions

**Analysis Date:** 2026-05-27

## Naming Patterns

### Go Backend:
- **Files:** `snake_case.go` for all source files (e.g. `article_repository.go`, `postgres_article.go`). Test files named `snake_case_test.go` in the same directory.
- **Packages:** All lowercase single words (e.g. `package entity`, `package usecase`, `package gateway`).
- **Structs & Interfaces:** `PascalCase` (e.g. `Article`, `PostgresArticleRepo`, `AIProcessor`). Interfaces do not have custom prefix (e.g. `ArticleRepository` interface, not `IArticleRepository`).
- **Functions:** `PascalCase` (e.g. `NewArticle`, `FetchAndSaveArticles`, `GetArticles`).
- **Variables:** `camelCase` for local variables and package-private parameters.

### React Frontend:
- **Files:** `PascalCase.jsx` for React components and Pages (e.g. `HomePage.jsx`, `LoginPage.jsx`, `NewsCard.jsx`). `camelCase.js` for hooks, services, and utilities (e.g. `api.js`, `useArticles.js`).
- **Components & Pages:** `PascalCase` (e.g. `HomePage`, `ArticlePage`).
- **Functions & Variables:** `camelCase` (e.g. `fetchArticles`, `isLoading`, `getHeaders`).
- **Constants:** `UPPER_SNAKE_CASE` (e.g. `API_BASE_URL`, `DEFAULT_LIMIT`).

---

## Code Style

### Go Backend:
- Standard `gofmt` code formatting.
- 2-space indentation or tabs (standard Go compiler formatting).
- Early returns for guard clauses:
  ```go
  if err != nil {
      return nil, fmt.Errorf("error details: %w", err)
  }
  ```

### React Frontend:
- ESLint configuration specified in `frontend/eslint.config.js`.
- 2-space indentation.
- Single quotes for strings, semicolons required.
- Zero console error standard (use structured error alerts, not `console.error` in committed code).

---

## Clean Architecture Boundaries

### Somut Sınıf İthalat Yasağı (Concrete Import Restriction):
* Usecase structures (`internal/usecase/`) and Handler structures (`internal/adapter/handler/`) must **never** import concrete repository or gateway adapter implementations (e.g. `PostgresArticleRepo` or `GeminiProcessor`).
* They must depend strictly on domain ports (`internal/domain/port/`) as interfaces.

### Durumsuz UseCase'ler (Stateless UseCases):
* All UseCase implementations must be completely stateless. They must not store request-scoped data, database records, or transaction states within their struct fields.

### DI Container Boundary:
* The `internal/infrastructure/container/container.go` file is strictly reserved for wiring up dependencies. No business logic, database queries, or route definitions allowed inside `container.go`.

---

## Presentation & Logic Separation (Frontend)

* React components must not perform direct HTTP/fetch network requests or contain heavy data filtering/caching logic.
* Always delegate data fetching lifecycle, mutations, and caching to TanStack React Query custom hooks (`frontend/src/hooks/` and `frontend/src/services/api.js`).

---

## Edge-Case & Fallback Rendering (B-Plan)

* The frontend must remain robust and fully functional even when backend AI summarization or translation fails (Gemini API rate limits/quota).
* Always render fallback strings for dynamic properties:
  * Show original English title (`article.title`) if Turkish translation (`article.turkish_title`) is missing.
  * Show original English content preview (`article.original_content`) in card summaries if the Turkish AI summary (`article.turkish_summary`) is not yet ready. Never leave components in infinite loading states.

---

## Error Handling

### Go Backend:
- Errors returned explicitly as the last return parameter.
- Use `fmt.Errorf("...: %w", err)` to wrap and preserve original errors.
- Always defer `rows.Close()` immediately after database query execution (`pool.Query`) to prevent connection pooling leaks.

### React Frontend:
- Use React Query `isError` blocks to display user-friendly warnings.
- Graceful API key verification error handling: Redirect to `/login` if admin APIs return `401 Unauthorized`.

---

*Convention analysis: 2026-05-27*
*Update when patterns change*
