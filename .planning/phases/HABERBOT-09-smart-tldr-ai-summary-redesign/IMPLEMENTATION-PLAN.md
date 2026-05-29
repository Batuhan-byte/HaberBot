# Implementation Plan: Smart TL;DR & AI Summary Redesign (Phase 9)

**Spec:** `docs/superpowers/specs/2026-05-29-smart-tldr-ai-summary-redesign-design.md`  
**Phase:** 9  
**Date Created:** 2026-05-29

---

## Overview

This phase moves AI summaries out of the article body and into a premium summary experience: a sticky CTA on article pages, a modal with skeleton loading, and markdown-safe rendering. The backend becomes the source of truth for summary delivery and generation, while the frontend keeps the article payload light and caches the first summary response.

## Architecture Decisions

1. **Lazy summary loading:** summary text is fetched only when the modal opens so the article payload stays smaller.
2. **Modal-first UX:** the summary is presented in a dedicated modal instead of inside the article body to preserve reading flow.
3. **Markdown-safe rendering:** summary output uses a markdown renderer with sanitization instead of raw HTML.
4. **Client cache:** once fetched, the summary stays cached in the client for the current article session.
5. **Backend fallback generation:** if the summary is missing, the backend generates it and returns the same response shape.

---

## Task List

### Phase 1: Backend Contract

#### Task 1: Redefine the article summary response contract
**Description:** Split the article detail payload from the AI summary payload so the default article response no longer carries summary text. Add a dedicated summary response shape for the modal flow.

**Acceptance Criteria:**
- [ ] Article detail response excludes summary text by default
- [ ] Summary endpoint returns a compact payload with the summary text and article id
- [ ] Existing article fields continue to work for the page body and metadata

**Verification:**
- [ ] Backend tests pass: `go test ./...`
- [ ] Manual API check confirms `/articles/:id` is smaller and `/articles/:id/summary` returns summary data

**Dependencies:** None

**Files Likely Touched:**
- `backend/internal/domain/entity/article.go`
- `backend/internal/adapter/handler/article_handler.go`
- `backend/internal/usecase/get_article.go`

**Estimated Scope:** Medium (3 files)

---

#### Task 2: Implement lazy summary fetch and generate flow
**Description:** Update the summary usecase/handler path so the summary endpoint returns cached content when present and generates/persists it when missing.

**Acceptance Criteria:**
- [ ] Summary endpoint returns saved summary immediately when one exists
- [ ] Missing summary triggers backend generation and persistence
- [ ] Errors return explicit HTTP status codes and messages
- [ ] Existing auth behavior for the summary route remains consistent

**Verification:**
- [ ] Unit tests cover cached summary, generated summary, and error cases
- [ ] Manual request confirms the first miss generates content and the second hit reuses it

**Dependencies:** Task 1

**Files Likely Touched:**
- `backend/internal/usecase/summarize_article.go`
- `backend/internal/adapter/handler/article_handler.go`
- `backend/internal/infrastructure/server/fiber.go`
- `backend/internal/usecase/summarize_article_test.go`

**Estimated Scope:** Medium (4 files)

---

### Checkpoint: Backend Contract
- [ ] Article payload is lean by default
- [ ] Summary endpoint works for cached and uncached articles
- [ ] Backend tests pass

---

### Phase 2: Frontend Summary Experience

#### Task 3: Add summary API client and markdown dependencies
**Description:** Extend the frontend API service with a dedicated summary fetch method and add the markdown rendering dependencies needed for rich summary output.

**Acceptance Criteria:**
- [ ] `api.getArticleSummary(id)` exists and calls the summary endpoint
- [ ] Frontend dependencies support markdown rendering and sanitization
- [ ] API client handles auth headers consistently with the existing article page flow

**Verification:**
- [ ] Frontend tests still pass after dependency install: `cd frontend && npm test`
- [ ] Frontend build still succeeds: `cd frontend && npm run build`

**Dependencies:** Task 2

**Files Likely Touched:**
- `frontend/src/services/api.js`
- `frontend/package.json`
- `frontend/package-lock.json`

**Estimated Scope:** Small (3 files)

---

#### Task 4: Build the summary modal and skeleton loader
**Description:** Create the premium summary modal with a skeleton loader, markdown rendering, and a cached query keyed by article id.

**Acceptance Criteria:**
- [ ] Modal opens and closes cleanly
- [ ] Skeleton loader shows while the summary request is in flight
- [ ] Markdown renders correctly with lists and bold text
- [ ] Successful fetches are cached for later modal opens

**Verification:**
- [ ] Component tests pass for loading, error, and markdown rendering states
- [ ] Manual browser test confirms the first open loads and the second open reuses cache

**Dependencies:** Task 3

**Files Likely Touched:**
- `frontend/src/components/ArticleSummaryModal/ArticleSummaryModal.jsx`
- `frontend/src/components/ArticleSummaryModal/ArticleSummaryModal.test.jsx`
- `frontend/src/components/ArticleSummaryModal/SummarySkeleton.jsx`

**Estimated Scope:** Medium (3 files)

---

#### Task 5: Add the sticky article summary CTA
**Description:** Add the premium sticky CTA that sits on the right side on desktop and pins to the bottom on mobile, then opens the summary modal.

**Acceptance Criteria:**
- [ ] CTA is visible on article pages in the intended responsive positions
- [ ] CTA opens the summary modal
- [ ] CTA state reflects loading/error/ready states cleanly

**Verification:**
- [ ] Responsive browser test confirms desktop sticky and mobile bottom placement
- [ ] CTA click opens the modal and loads the summary

**Dependencies:** Task 4

**Files Likely Touched:**
- `frontend/src/components/ArticleSummaryLauncher/ArticleSummaryLauncher.jsx`
- `frontend/src/components/ArticleSummaryLauncher/ArticleSummaryLauncher.test.jsx`
- `frontend/src/pages/ArticlePage/ArticlePage.jsx`

**Estimated Scope:** Medium (3 files)

---

### Checkpoint: Frontend Experience
- [ ] CTA is responsive and premium-looking
- [ ] Modal loads and renders markdown
- [ ] Summary cache works on repeat opens

---

### Phase 3: Integration and Polish

#### Task 6: Wire the article page to the new summary flow
**Description:** Remove the inline summary block from the article body, mount the sticky CTA, and keep the article reading experience uninterrupted.

**Acceptance Criteria:**
- [ ] No summary block appears inline in the article body
- [ ] Article content continues to render normally
- [ ] Summary is only reachable through the CTA/modal path

**Verification:**
- [ ] Manual article-page review confirms the body is clean and summary lives in the modal
- [ ] Regression tests still pass

**Dependencies:** Tasks 1-5

**Files Likely Touched:**
- `frontend/src/pages/ArticlePage/ArticlePage.jsx`

**Estimated Scope:** Small (1 file)

---

#### Task 7: Add regression tests and handoff notes
**Description:** Add final backend/frontend regression coverage and update documentation so the new summary experience is easy to maintain.

**Acceptance Criteria:**
- [ ] Backend tests cover summary contract and lazy generation
- [ ] Frontend tests cover CTA, modal, skeleton, markdown, and cache behavior
- [ ] Docs note the new lazy summary flow and markdown rendering dependency

**Verification:**
- [ ] `go test ./...` passes
- [ ] `cd frontend && npm test` passes
- [ ] `cd frontend && npm run build` passes

**Dependencies:** Tasks 1-6

**Files Likely Touched:**
- `backend/internal/usecase/summarize_article_test.go`
- `frontend/src/components/ArticleSummaryModal/ArticleSummaryModal.test.jsx`
- `frontend/src/components/ArticleSummaryLauncher/ArticleSummaryLauncher.test.jsx`
- `README.md` or relevant docs

**Estimated Scope:** Medium (4 files)

---

### Checkpoint: Complete
- [ ] Summary UX is modal-only and markdown-safe
- [ ] Lazy load + cache behavior verified
- [ ] Backend and frontend tests pass
- [ ] Documentation updated

---

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Markdown renders unsafely | High | Use `react-markdown` with sanitization and avoid raw HTML rendering |
| Summary payload still leaks into article detail | Medium | Split article detail and summary response contracts early |
| Modal feels slow on open | Medium | Show skeleton immediately and cache the first successful result |
| Auth behavior diverges between page and summary endpoint | Medium | Keep the existing summary-route auth policy explicit in the contract and tests |
| Frontend dependency addition causes build regressions | Low | Add the dependency in a dedicated task and verify build/test immediately |

## Open Questions
- None after the current lock-in; the remaining work is implementation.
