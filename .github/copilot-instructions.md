# Copilot instructions for HaberBot

## Build, test, lint

**Backend (Go, from `backend/`)**
- Run migrations: `go run run_migrations.go`
- Start API server: `go run cmd/server/main.go`
- Run tests: `go test -v ./...`
- Run a single test: `go test -v -run TestArticleFetch`

**Frontend (React, from `frontend/`)**
- Dev server: `npm run dev`
- Build: `npm run build`
- Lint: `npm run lint`
- Run tests: `npm run test`
- Run a single test file: `npm run test -- src/components/ArticleCard/ArticleCard.test.jsx`

## High-level architecture

- **Clean Architecture backend** under `backend/internal/`: `domain` defines entities + ports; `usecase` implements business workflows; `adapter` contains handlers, repositories, gateways, and fetchers; `infrastructure` handles config, DB, scheduler, and the DI container.
- **Dependency wiring** lives in `backend/internal/infrastructure/container/container.go`. It chooses the AI processor (Gemini vs OpenAI), builds fetchers (HackerNews + RSS), wires use cases and handlers, and sets up the scheduler.
- **Data flow**: fetchers collect external content → use cases fetch/process (AI translation/summary) → repositories persist to Postgres → handlers expose REST endpoints.
- **Frontend** uses React Query hooks in `frontend/src/hooks/` and API helpers in `frontend/src/services/` to call the backend; pages/components are presentation-focused.

## Key conventions

- `usecase/` and handler code must **only** depend on domain ports (`internal/domain/port`), never concrete adapter implementations.
- Use cases are **stateless**; do not store request or transaction data in struct fields.
- `internal/infrastructure/container/container.go` is **wiring only**—no business logic or route definitions.
- React components must not perform direct network calls; use React Query hooks and `services/api.js`.
- UI must render fallbacks when AI results are missing: use `article.title` if `article.turkish_title` is empty, and `article.original_content` if `article.turkish_summary` is missing; avoid infinite loading states.
- Error handling: Go errors are wrapped with `fmt.Errorf("...: %w", err)` and DB rows are closed immediately after queries; frontend admin APIs should redirect to `/login` on `401` and show React Query `isError` states.
