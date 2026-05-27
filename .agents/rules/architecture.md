---
trigger: always_on
glob: **/*.{go,js,jsx,css}
description: Architectural guidelines, decoupling boundaries, and structural patterns for HaberBot established by Principal Architect code mapping.
---

# HaberBot Architectural Standards & Guidelines

This document outlines the official architectural guidelines, separation boundaries, and software decoupling standards for the HaberBot application. All developers and AI agents must strictly adhere to these rules when modifying or extending the system.

---

## 1. Clean Architecture Boundaries (Backend)

The Go backend strictly implements **Clean Architecture** (Ports & Adapters). To prevent this from degrading into spaghetti code, follow these boundaries:

### Somut Sınıf İthalat Yasağı (Concrete Import Restriction):
* **Rule:** UseCase structures (`internal/usecase/`) and Handler structures (`internal/adapter/handler/`) must **never** import concrete repository or gateway adapter implementations (e.g. `PostgresArticleRepo` or `GeminiProcessor`).
* **Interface Decoupling:** They must depend strictly on domain ports (`internal/domain/port/`) as interfaces.
  * *Why:* Decoupling domain logic from the database and external APIs ensures that changes to infrastructure (e.g., migrating from PostgreSQL to MongoDB, or Gemini to Claude) do not affect core business logic.

### Durumsuz UseCase'ler (Stateless UseCases):
* **Rule:** All UseCase implementations must be completely stateless. They must not store request-scoped data, database records, or transaction states within their struct fields.
* **Orchestration Only:** UseCases should only coordinate domain entities and interact with the database/API services purely through the defined ports.

### Dependency Injection (DI) Sınırları (`NewContainer` Limit):
* **Rule:** The `internal/infrastructure/container/container.go` file is strictly reserved for wiring up dependencies (Dependency Injection).
* **No Logic in Container:** Absolutely no business logic, route definitions, database query executions, or active environment configurations should exist in `container.go`. Keep it clean as a pure wiring factory.

---

## 2. Decoupled Component & State Architecture (Frontend)

To keep the Vite-React frontend highly responsive, clean, and modular, follow these presentation rules:

### Arayüzde Mantık ve Sunum Ayrımı (Presentation & Logic Separation):
* **Rule:** React components (`frontend/src/components/` and `frontend/src/pages/`) must not perform direct HTTP/fetch network requests or contain heavy data filtering/caching logic.
* **Custom Hooks & React Query:** Always delegate data fetching lifecycle, mutations, and caching to TanStack React Query custom hooks (`frontend/src/hooks/` and `frontend/src/services/api.js`).
* **Visual Only:** Components must remain strictly presentation-oriented, focusing on rendering UI and handling user click/input interactions.

### Edge-Case & B-Plan (Fallback) Gösterimi:
* **Rule:** The frontend must remain robust and fully functional even when backend AI summarization or translation fails (due to Gemini API rate limits/quota exhaustion).
* **Resilient Data Fields:** Always render fallback strings for dynamic properties:
  * Show original English title (`article.title`) if Turkish translation (`article.title_tr`) is missing.
  * Show original English content preview (`article.original_content`) in card summaries if the Turkish AI summary (`article.summary_tr`) is not yet generated. Never leave components in static, infinite loading states (e.g., "loading...").

---

## 3. Testing & Validation Standards

* **TDD & Unit Testing:** All complex data parsing (such as AI response processing) and boundary handlers must have accompanying table-driven unit tests using standard library `testing` in Go.
* **Zero Console Error Standard:** The frontend must maintain a clean browser console. No console warnings or uncaught exceptions should occur during standard E2E workflows.
