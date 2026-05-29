# Spec: Smart TL;DR & AI Summary Redesign

## Objective
Authenticated readers can open a premium AI summary modal from a sticky CTA on the article detail page without interrupting the reading flow. The summary is fetched lazily, cached after the first open, and rendered with markdown support so `**bold**`, lists, and other formatting display correctly. The default article payload stays lean by excluding summary text from the initial page load.

## Tech Stack
- Backend: Go 1.26, Fiber
- Frontend: React 19, Vite, TanStack Query
- Markdown rendering: `react-markdown`, `remark-gfm`, `rehype-sanitize`

## Commands
- Backend build: `go build ./...`
- Backend test: `go test ./...`
- Frontend install: `cd frontend && npm install`
- Frontend build: `cd frontend && npm run build`
- Frontend test: `cd frontend && npm test`
- Frontend lint: `cd frontend && npm run lint`
- Dev backend: `go run cmd/server/main.go`
- Dev frontend: `cd frontend && npm run dev`

## Project Structure
- `backend/internal/usecase/summarize_article.go` → summary fetch/generate behavior
- `backend/internal/adapter/handler/article_handler.go` → summary endpoint + article response shaping
- `backend/internal/infrastructure/server/fiber.go` → route wiring
- `frontend/src/services/api.js` → summary API client
- `frontend/src/pages/ArticlePage/ArticlePage.jsx` → CTA, modal trigger, page integration
- `frontend/src/components/ArticleSummaryModal/` → modal, loader, markdown renderer
- `frontend/src/components/ArticleSummaryLauncher/` → sticky desktop/mobile CTA

## Code Style
Use small React components, keep data fetching in `src/services/api.js`, and render markdown through a dedicated component instead of inline HTML.

```jsx
function ArticleSummaryModal({ articleId, open, onClose }) {
  const summaryQuery = useQuery({
    queryKey: ['article-summary', articleId],
    queryFn: () => api.getArticleSummary(articleId),
    enabled: open && Boolean(articleId),
  });

  if (!open) return null;
  if (summaryQuery.isLoading) return <SummarySkeleton />;

  return (
    <Modal onClose={onClose}>
      <ReactMarkdown remarkPlugins={[remarkGfm]} rehypePlugins={[rehypeSanitize]}>
        {summaryQuery.data.summary}
      </ReactMarkdown>
    </Modal>
  );
}
```

## Testing Strategy
- Backend: unit tests for summary usecase, handler contract, and response shaping with `go test ./...`
- Frontend: component tests for CTA visibility, modal loading state, markdown rendering, and cache reuse with `cd frontend && npm test`
- Integration: article page opens modal, shows skeleton, renders markdown, and reopens from cache
- Regression: ensure article payload no longer includes summary text by default

## Boundaries
- Always: keep the article page readable, lazy-load summary content, sanitize markdown, cache the first successful summary response, and preserve the existing article layout
- Ask first: adding new markdown dependencies, changing summary auth policy, altering article response shape beyond summary omission, or introducing a background job/queue
- Never: use `dangerouslySetInnerHTML` for summary output, block article rendering on summary fetch, ship unsanitized markdown, or duplicate summary content in the initial article payload

## Success Criteria
- Article detail pages no longer render an inline summary block
- A sticky AI Summary CTA appears on desktop and mobile article views
- Opening the modal shows a skeleton while summary content loads
- Markdown formatting renders correctly in the modal
- The first open fetches summary asynchronously and later opens reuse cached data
- When a summary is missing, the backend generates it and returns it through the same summary flow
- Initial article payload is smaller because summary text is excluded by default

## Open Questions
- None at this stage.
