<!-- Context: project-intelligence/technical | Priority: critical | Version: 1.1 | Updated: 2026-05-01 -->

# Technical Domain

**Purpose**: Tech stack, architecture, and coding patterns for Grubbin Data.
**Last Updated**: 2026-05-01
**Update Triggers**: Tech stack changes | New patterns | Architecture decisions

---

## Quick Reference

| Layer | Technology | Role |
|-------|-----------|------|
| Backend | Go 1.24 | API server, business logic |
| Router | `go-chi/chi/v5` | HTTP routing & middleware |
| Config | `godotenv` | `.env` file loading |
| Logging | `log/slog` (stdlib) | Structured JSON logs |
| DB (Phase 2+) | PostgreSQL + `pgx` + `golang-migrate` | Persistence & migrations |
| Frontend | React 19 + Vite + TypeScript | UI layer |
| Runtime | Bun | Package manager & script runner |
| Data Fetch | TanStack Query | Server-state caching |
| Styling | CSS Modules | Scoped styles |
| DevOps | Docker + Docker Compose | Containerization |

---

## Architecture Patterns

**Backend: Functional Core / Imperative Shell (FC/IS)**
- `internal/core/` — pure functions, zero I/O, zero external imports.
- `internal/services/` — orchestrators: compose core + side effects (DB, time).
- `internal/adapters/http/` — thin HTTP handlers (chi router).
- `internal/repositories/` — DB interfaces + implementations (Phase 2+).

**Frontend: Pure Components + Custom Hooks**
- Components are pure: props → JSX, no side effects during render.
- Custom hooks isolate data fetching and browser side effects.
- TanStack Query handles server-state caching and background updates.

---

## Naming & Readability Standards

**Self-documenting code over comments.** Use descriptive names so the code explains itself. Comments are reserved for "why," not "what."

| Element | Rule | Example |
|---------|------|---------|
| Variables | Describe content/purpose, not type | `deliveryCount` not `n`; `earningsPerHour` not `eph` |
| Functions | Verb + noun, describe action | `CalculateEarningsPerHour` not `CalcEPH` |
| Parameters | Full words, avoid abbreviations | `deliveryDate` not `dt`; `restaurantName` not `rName` |
| Struct fields | Clear domain language | `BasePayCents` not `bp`; `TipAmountCents` not `tip` |
| Booleans | `Is`/`Has`/`Can` prefix | `IsCompleted` not `done`; `HasValidTip` not `valid` |
| Errors | Prefix with failed action | `ErrInvalidDeliveryDate` not `ErrBadDate` |
| React props | Mirror domain model | `deliveryDate` not `date`; `earningsTotal` not `total` |
| CSS Modules | Descriptive class names | `.deliveryCard` not `.dc`; `.earningsHighlight` not `.eh` |

**Avoid**: single-letter vars (except `i` in tight loops), abbreviations, Hungarian notation, type-encoded names (`strName`, `intCount`).

---

## API Pattern (Go)

```go
func (handler *HelloWorldHandler) GetGreeting(
    responseWriter http.ResponseWriter,
    request *http.Request,
) {
    greetingMessage, greetingError := handler.helloWorldService.GenerateGreeting(request.Context())
    if greetingError != nil {
        handler.respondWithError(responseWriter, greetingError)
        return
    }
    handler.respondWithJSON(responseWriter, http.StatusOK, greetingMessage)
}
```

**Rules**: Handlers are thin — parse request, call service, encode response. Never put business logic in handlers. Name variables after their domain meaning.

---

## Component Pattern (React)

```tsx
interface HelloWorldProps {
  greetingMessage: string;
  generatedAtTimestamp: string;
}

export function HelloWorld({ greetingMessage, generatedAtTimestamp }: HelloWorldProps) {
  return (
    <article className={styles.greetingContainer}>
      <h1 className={styles.greetingMessage}>{greetingMessage}</h1>
      <time className={styles.generatedAtTimestamp} dateTime={generatedAtTimestamp}>
        {generatedAtTimestamp}
      </time>
    </article>
  );
}
```

**Rules**: Props in, JSX out. Side effects live in hooks, never during render. Name props to mirror the backend domain model.

---

## Naming Conventions

| Type | Convention | Example |
|------|-----------|---------|
| Go files | lowercase, descriptive | `hello_world_handler.go` |
| Go exported | PascalCase, descriptive | `HelloWorldHandler` |
| Go unexported | camelCase, descriptive | `helloWorldService` |
| TS/TSX files | kebab-case, descriptive | `hello-world.tsx` |
| React components | PascalCase, descriptive | `HelloWorld` |
| Custom hooks | camelCase, `use` prefix, descriptive | `useHelloWorldGreeting` |
| DB tables (Phase 2+) | snake_case, plural | `delivery_records` |
| CSS Modules | camelCase or kebab-case, descriptive | `.greetingContainer` or `.greeting-container` |

---

## Code Standards

- **Go**: Every function returns `(result, error)`. Check `error != nil` immediately.
- **Go interfaces are implicit**: structural typing, like TypeScript.
- **No auth in Phase 1**. Keep scaffold lean.
- **No hot-reload tool for Go** (no `air`). Use `go run` inside Docker.
- **Bun replaces npm**: use `bun install`, `bun run`, `bun.lockb`.
- **Vite proxy**: forward `/api` to backend in `vite.config.ts` to avoid CORS.
- **Frontend fetch**: prefer TanStack Query over raw `useEffect` + `fetch`.
- **Self-documenting names**: prefer `delivery.EarningsPerHour` over `d.Eph` + comment.
- **Theme system**: CSS custom properties + `data-theme` attribute on `<html>`.
- **Theme sources**: `@catppuccin/palette` for Catppuccin; official repos for Tokyo Night/Gruvbox.

---

## Security Requirements

- Validate all user input with `go-playground/validator` (Phase 2+).
- Use parameterized queries via `pgx` (Phase 2+).
- No auth in Phase 1; evaluate JWT middleware in chi if added later.
- CORS handled via Vite proxy in dev; chi CORS middleware if needed.

---

## 📂 Codebase References

**Backend entry**: `backend/cmd/api/main.go` — wires config, router, server.
**Frontend entry**: `frontend/src/main.tsx` — React root + QueryClientProvider.
**API client**: `frontend/src/services/api.ts` — pure fetch wrapper.
**Config**: `backend/internal/config/config.go` — typed env loader.
**Theme registry**: `frontend/src/themes/registry.ts` — centralized theme definitions.
**Theme hook**: `frontend/src/hooks/useTheme.tsx` — context provider + persistence.
**Theme UI**: `frontend/src/components/ThemeSwitcher.tsx` — dropdown selector.
**Global styles**: `frontend/src/index.css` — theme-aware body/html styles.
**Theme CSS**: `frontend/src/styles/themes.css` — auto-generated variables.
**Theme script**: `frontend/scripts/generate-themes.ts` — build-time generator.

---

## Related Files

- `AGENTS.md` — Project overview, conventions, commands.
- `docker-compose.yml` — Service orchestration.
