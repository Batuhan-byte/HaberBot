# Technology Stack

**Analysis Date:** 2026-05-27

## Languages

**Primary:**
- Go 1.26.0 - All backend application code (Clean Architecture structure, HTTP routers, RSS fetchers, processors).
- JavaScript (ES6+ Modules) - All frontend application code (React components, routing, services).

**Secondary:**
- SQL (PostgreSQL Dialect) - Database schema migrations under `backend/migrations/` and queries in repository adapters.
- HTML5 / CSS3 - Frontend layout and responsive magazine-quality styling under `frontend/src/index.css`.

## Runtime

**Environment:**
- **Backend:** Native Go binary compiled for production.
- **Frontend:** Modern Web Browser, compiled and bundled via Vite.
- **Local Dev Server:** Node.js 20.x environment running Vite for frontend development and standard Go environment for backend development.

**Package Manager:**
- **Go Mod:** Go modular system using `backend/go.mod` and `backend/go.sum`.
- **npm 10.x:** Frontend package management using `frontend/package.json` and `frontend/package-lock.json`.

## Frameworks

**Core:**
- **Go Fiber v2 (`github.com/gofiber/fiber/v2`):** Modern and high-performance Go web server framework for building the REST API.
- **React 19 (`^19.2.6`):** Standard Single-Page Application (SPA) UI framework.
- **React Router v7 (`^7.15.1`):** Frontend client-side routing library.

**Data Fetching & State:**
- **TanStack React Query v5 (`^5.100.14`):** Frontend remote data fetching, caching, query states, and automatic UI synchronizations.

**Testing:**
- **Go testing standard library:** Used with `github.com/stretchr/testify` for table-driven backend tests.
- **Vitest v4 (`^4.1.7`):** Frontend JavaScript unit testing framework.
- **React Testing Library (`^16.3.2`):** Frontend component and UI testing.

**Build/Dev:**
- **Vite v8 (`^8.0.12`):** Fast frontend bundler and dev server.
- **ESLint v10 (`^10.3.0`):** Frontend code quality linting.

## Key Dependencies

**Critical:**
- `github.com/google/generative-ai-go v0.11.0` - Official Google Gemini SDK for article summarization and translation.
- `github.com/jackc/pgx/v5 v5.5.5` - PostgreSQL driver and connection pool (`pgxpool.Pool`) management.
- `github.com/mmcdole/gofeed v1.3.0` - RSS feed fetching and parsing utility.
- `github.com/go-shiori/go-readability v0.0.0` - Article extraction engine to parse clean reading content from raw HTML payloads.
- `github.com/robfig/cron/v3 v3.0.1` - Cron scheduler for managing automated background news fetch and process tasks.

**Infrastructure:**
- `github.com/joho/godotenv v1.5.1` - Environment file loader (.env).
- `github.com/PuerkitoBio/goquery v1.8.0` - HTML document parsing.

## Configuration

**Environment:**
- **Backend:** Configured via `backend/.env` file. Requires `DATABASE_URL` (Postgres connection string), `GEMINI_API_KEY` (Gemini API access), and `ADMIN_API_KEY` (admin route authorization).
- **Frontend:** Configured via `frontend/.env` file. Defines `VITE_API_BASE_URL` pointing to Go backend.

**Build:**
- `frontend/vite.config.js` - Vite compiler/bundler configurations.
- `frontend/eslint.config.js` - Lint rules for JavaScript.

## Platform Requirements

**Development:**
- Cross-platform: Runs on Windows, macOS, and Linux.
- Requires local Go 1.26+ installation, Node.js 20+ installation, and PostgreSQL database.

**Production:**
- Native Docker support via `backend/Dockerfile` and root `docker-compose.yml`.
- Deployable to cloud platforms (e.g., Railway, Heroku, AWS) with serverless PostgreSQL.

---

*Stack analysis: 2026-05-27*
*Update after major dependency changes*
