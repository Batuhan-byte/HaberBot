# Codebase Concerns

**Analysis Date:** 2026-05-27

## Tech Debt

**Gemini API rate limiting and sleep logic:**
- Issue: Google Gemini API has quota limits. If many articles are processed in a loop, it throws `429 Too Many Requests`.
- Files: `backend/internal/adapter/gateway/gemini_processor.go`
- Why: Simple loops without retry logic or dynamic delays are prone to immediate failure.
- Impact: Automated updates stop processing when quotas are reached.
- Fix approach: Implement a uniform delay (`time.Sleep`) on both success and failure paths inside the processing loop to prevent quota spikes. Add a retry wrapper with exponential backoff.

**Direct API key verification for admin endpoints:**
- Issue: Custom `X-Admin-API-Key` is parsed and checked directly in hand-crafted controller handlers rather than a shared Fiber middleware.
- Files: `backend/internal/adapter/handler/admin_handler.go`
- Why: Built rapidly for standalone admin task triggers.
- Impact: If new admin endpoints are added, key verification could be accidentally omitted.
- Fix approach: Extract admin key authentication to a shared Fiber middleware under `backend/internal/infrastructure/middleware/admin_auth.go`.

---

## Security Considerations

**Static Admin API Key:**
- Risk: Using a single static API key shared globally via `ADMIN_API_KEY` in environment variables. If leaked, anyone can post/delete topics or fetch arbitrary feeds.
- Files: `backend/internal/adapter/handler/admin_handler.go`, `frontend/src/services/api.js`
- Current mitigation: Key stored on client `localStorage` and sent over HTTPS.
- Recommendations: Implement standard secure JWT authentication or time-limited admin tokens, rather than relying on a static string.

**Unchecked RSS payloads / XSS risk:**
- Risk: RSS feed summaries might contain malicious HTML or scripts. The React frontend uses React's native HTML interpreter (`dangerouslySetInnerHTML`) to display news content.
- Files: `frontend/src/pages/ArticlePage/index.jsx`
- Current mitigation: Renders raw HTML using `dangerouslySetInnerHTML`.
- Recommendations: Integrate a library like DOMPurify or sanitize HTML on the Go backend using a sanitization package before persisting raw RSS text to PostgreSQL.

---

## Performance Bottlenecks

**Sequential RSS article parsing & translation in cron loops:**
- Problem: Cron job parses RSS feeds, downloads page HTML, runs readability text extraction, translates title, and generates summaries sequentially inside a single thread.
- Files: `backend/internal/usecase/article_usecase.go`
- Measurement: Takes up to 10-15 seconds per article depending on Gemini API latency.
- Cause: Synchronous single-threaded loops.
- Improvement path: Optimize background processing using parallel Go worker pools (goroutines and channels), with strict concurrency limits to avoid Gemini API rate limit constraints.

---

## Fragile Areas

**Frontend Date Parsing Mismatches:**
- Why fragile: Dates in RSS feeds use diverse structures (RFC1123, ISO8601). Go backend stores timestamps in Postgres, but camelCase or snake_case key mismatches on fetched JSON properties can break date parsing in React.
- Common failures: Client UI crashes or shows "Invalid Date" for newly parsed articles.
- Safe modification: Standardize date conversions in a central utility checking both `article.created_at` and `article.fetched_at`:
  ```javascript
  const displayDate = new Date(article.created_at || article.fetched_at).toLocaleString();
  ```
- Test coverage: Handled via client fallback logic, but lacks comprehensive integration test coverage.

---

## Dependencies at Risk

**OpenAI Gateway Client:**
- Risk: Go OpenAI library (`github.com/sashabaranov/go-openai`) updates frequently. Major changes to client initializers can cause compilation issues on backend upgrades.
- Impact: Fallback processor fails to compile.
- Migration plan: Lock dependency versions in `go.mod` and add comprehensive integration tests covering OpenAI fallback behaviors.

---

*Concerns audit: 2026-05-27*
*Update as issues are fixed or new ones discovered*
