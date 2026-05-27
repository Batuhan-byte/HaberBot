# External Integrations

**Analysis Date:** 2026-05-27

## APIs & External Services

**Artificial Intelligence / LLM Providers:**
- **Google Gemini API** - Used for automated news translation (English to Turkish) and Turkish summary generation.
  - SDK/Client: `github.com/google/generative-ai-go` v0.11.0 (Official Go Client SDK)
  - Auth: Authenticated using API Key via the `GEMINI_API_KEY` environment variable.
  - Model Used: `gemini-1.5-flash` or newer (for summarization/translation tasks).
  
- **OpenAI API** - Alternative/fallback LLM processor for translation and summarization.
  - SDK/Client: `github.com/sashabaranov/go-openai` v1.41.2
  - Auth: Authenticated using API Key via the `OPENAI_API_KEY` environment variable.

**News RSS Feeds:**
- **Feeds RSS Fetching** - Periodically scrapes designated external tech/general news RSS feeds (such as Hacker News, TechCrunch, etc.).
  - Client: `github.com/mmcdole/gofeed` v1.3.0
  - Auth: Publicly available feeds (no authorization required).

## Data Storage

**Databases:**
- **PostgreSQL Database** - Primary persistent data store for articles, topics, and admin settings.
  - Connection: Multi-client connection pooling managed via `pgxpool.Pool` (PostgreSQL Driver: `github.com/jackc/pgx/v5`).
  - Connection Config: Configured via `DATABASE_URL` environment variable.
  - Migrations: SQL migration scripts located in `backend/migrations/` (.up.sql / .down.sql).

## Authentication & Identity

**Admin Access Authorization:**
- **Custom Header-Based API Key** - Secured routes on Go backend use a custom authentication middleware that validates requests against a static token.
  - Implementation: Backend checks for header `'X-Admin-API-Key'` which must match the backend's `ADMIN_API_KEY` environment variable.
  - Token Storage: Stored securely on the React client in `localStorage` under `admin_api_key`.
  - Session Management: Checked on API request lifecycle. Returns `401 Unauthorized` on mismatch and redirects client to the `/login` view.

## Monitoring & Observability

**Error Tracking:**
- **Console Log Stream** - Standard stdout/stderr logging outputs. No external telemetry service like Sentry or Datadog is currently integrated.

## CI/CD & Deployment

**Hosting:**
- **Railway / Docker Support** - The project contains configuration for multi-container deployments.
  - Deployment Configuration: Configured via root `docker-compose.yml`, backend `Dockerfile`, and `railway.toml` at root.
  - Environment variables: Set securely on the hosting platform's dashboard.

## Environment Configuration

**Development:**
- Required env vars:
  - **Backend (`backend/.env`):**
    - `DATABASE_URL` - Local PostgreSQL connection string.
    - `GEMINI_API_KEY` - Developer API key for Google Gemini services.
    - `ADMIN_API_KEY` - Developer-selected string for authenticating admin routes.
  - **Frontend (`frontend/.env`):**
    - `VITE_API_BASE_URL` - Go backend local host address (e.g. `http://localhost:8080`).
- Mock/stub services:
  - Unit tests for LLM processors mock the `AIProcessor` port to avoid network traffic and API billing costs during builds.

**Production:**
- Managed via production dashboard environment variables. Uses robust serverless Postgres with SSL.

---

*Integration audit: 2026-05-27*
*Update when adding/removing external services*
