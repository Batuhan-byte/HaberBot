# Codebase Structure

**Analysis Date:** 2026-05-27

## Directory Layout

```
HaberBot/
├── .agent/             # GSD internal environment, settings, and skills
├── .agents/            # IDE agent caching and integration files
├── .git/               # Git version control metadata
├── backend/            # Go Clean Architecture Web API Server
│   ├── cmd/            # Command line entry points
│   │   └── server/     # Main web API bootstrap
│   ├── internal/       # Clean architecture components
│   │   ├── adapter/    # Repositories, API gateways, HTTP controllers, and RSS fetchers
│   │   ├── domain/     # Core domain entities, value objects, and repository ports (interfaces)
│   │   ├── infrastructure/ # Fiber server router initialization, cron triggers, and DI container
│   │   └── usecase/    # Stateless core application business logic
│   └── migrations/     # PostgreSQL SQL database migration scripts
├── frontend/           # Vite-React Single Page Application (SPA)
│   ├── public/         # Static assets and browser tab icons
│   └── src/            # JavaScript React codebase
│       ├── assets/     # Images and vector styling icons
│       ├── components/ # Presentational visual components
│       ├── hooks/      # TanStack Query custom data-fetching hooks
│       ├── pages/      # Routeable view screens
│       ├── services/   # Fetch client abstractions (api.js)
│       └── utils/      # Client-side helpers and constant configurations
├── docker-compose.yml  # Local multi-container Docker composition configuration
├── DESIGN.md           # User visual identity specifications
└── README.md           # Developer onboarding document
```

---

## Directory Purposes

**backend/cmd/**
- Purpose: Entry point binary builders.
- Contains: `server/main.go` - parses `.env` and initializes server.

**backend/internal/domain/**
- Purpose: Pure domain entities and adapter contracts.
- Contains: `entity/article.go` (models), `port/article_repository.go` (interfaces).
- Key files: `port/ai_processor.go` - abstract Gemini SDK wrapper definitions.

**backend/internal/usecase/**
- Purpose: Coordinates business logic orchestrations.
- Contains: Pure stateless Go usecase handlers.
- Key files: `article_usecase.go` - handles fetching, summary trigger, and persistence flows.

**backend/internal/adapter/**
- Purpose: Concrete implementations of external communication layers (DB, RSS, REST, Gemini API).
- Contains: `repository/` (SQL queries), `gateway/` (Gemini SDK integration), `fetcher/` (RSS downloaders), `handler/` (Fiber routes).
- Key files: `gateway/gemini_processor.go` - concrete Google Gemini translation and summary wrapper.
- Key files: `repository/postgres_article.go` - pgxpool PostgreSQL repository queries.

**backend/internal/infrastructure/**
- Purpose: Application routing, scheduling and Dependency Injection.
- Contains: `container/` (wires up adapters with usecases), `router.go` (Fiber routes), `cron/` (triggers background fetcher loops).
- Key files: `container/container.go` - dependency injection factory.

**frontend/src/components/**
- Purpose: Presentational UI design layouts.
- Contains: React components with semantic layouts, premium hover animations, and loading skeletons.

**frontend/src/hooks/**
- Purpose: Remote caching and data state sync using TanStack Query.
- Contains: Hooks like `useArticles.js` that decouple presentational rendering from backend network calls.

**frontend/src/pages/**
- Purpose: Screen-sized page routing endpoints.
- Contains: Page wrappers for `HomePage`, `ArticlePage`, `TopicPage`, `AdminPage`, `LoginPage`.

---

## Key File Locations

**Entry Points:**
- `backend/cmd/server/main.go` - Go Web API Server bootstrap entry point.
- `frontend/src/main.jsx` - Vite-React client compilation entry point.

**Configuration:**
- `backend/.env` - Backend database and API keys.
- `frontend/.env` - Frontend API endpoint configurations.
- `docker-compose.yml` - Multi-service local environment launcher.

**Core Logic:**
- `backend/internal/usecase/` - Go Clean Architecture stateless business logic.
- `frontend/src/services/api.js` - Centralized fetch client orchestrations.

**Testing:**
- `backend/internal/domain/entity/article_test.go` - Domain unit tests.
- `backend/internal/adapter/handler/admin_handler_test.go` - Route handler tests.
- `frontend/src/test/` - Frontend Vitest and React testing library scripts.

---

## Naming Conventions

**Files:**
- `snake_case.go` - Go modules and standard library source structures.
- `snake_case_test.go` - Go test modules.
- `PascalCase.jsx` - React components and view page classes.
- `camelCase.js` - JavaScript services, hooks, utilities, and configs.
- `XXX_migration_name.up.sql` / `down.sql` - Database SQL migrations in chronological order.

**Directories:**
- `snake_case` - All backend packaging subdirectories.
- `camelCase` - All frontend source subdirectories.

---

## Where to Add New Code

**New Feature (Backend):**
1. Add new database tables: Create SQL migrations under `backend/migrations/`.
2. Define models & ports: Create entities under `backend/internal/domain/entity/` and interfaces under `backend/internal/domain/port/`.
3. Code Usecase: Implement stateless business flow under `backend/internal/usecase/`.
4. Code Repository/Gateway Adapter: Implement SQL methods under `backend/internal/adapter/repository/`.
5. Code Route Handler: Implement controller under `backend/internal/adapter/handler/`.
6. Wire it up: Add wiring calls inside the DI Factory in `backend/internal/infrastructure/container/container.go`.

**New Component (Frontend):**
- Implementation: `frontend/src/components/{ComponentFolder}/index.jsx`
- Styling: Styled within `index.css` under modern typography selectors or inside local CSS/glassmorphic wrappers.
- Custom Hooks: Delegate remote operations under `frontend/src/hooks/`.

---

## Special Directories

**graphify-out/**
- Purpose: Automatically generated code dependency and design graph.
- Source: Generated by CLI tool `graphify`.
- Committed: No (Ignored under `.gitignore`).

**.agent/**
- Purpose: GSD Redux CLI active execution context and custom skills.
- Source: Set up on workspace environment initialization.
- Committed: No (Ignored under `.gitignore`).

---

*Structure analysis: 2026-05-27*
*Update when directory structure changes*
