# AGENTS.md — Grubbin Data

> Repo for a Grubhub delivery data analysis app. Go backend + React (Bun) frontend.
> Developer background: backend-focused, TypeScript/Node/Bun.

---

## Purpose

Accept and analyze Grubhub delivery data over time. Starting as a study project for Go and React with functional programming patterns.

---

## Architecture

- **Backend**: Go 1.24, Functional Core / Imperative Shell (FC/IS).
  - `internal/core/` — pure functions, zero I/O, zero external imports.
  - `internal/services/` — orchestrators: compose core + side effects (DB, time).
  - `internal/adapters/http/` — thin HTTP handlers (chi router).
  - `internal/repositories/` — DB interfaces + implementations.
- **Frontend**: React 19 + Vite + TypeScript, executed via Bun.
  - Pure components (props → JSX, no side effects).
  - Custom hooks isolate data fetching / side effects.
  - TanStack Query for server state.

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go 1.24 |
| Router | `go-chi/chi/v5` |
| Config | `godotenv` |
| Logging | `log/slog` (stdlib) |
| DB (Phase 2+) | PostgreSQL + `pgx` + `golang-migrate` |
| Frontend | React 19 + Vite + TypeScript |
| Frontend runtime | Bun (package install + script runner) |
| Frontend data | TanStack Query |
| Frontend themes | `@catppuccin/palette` + CSS custom properties |
| Styling | CSS Modules |
| DevOps | Docker + Docker Compose |

---

## Monorepo Layout

```
grubbin-data/
├── backend/
│   ├── cmd/api/main.go              # Entry point
│   ├── internal/
│   │   ├── config/                  # Env loader
│   │   ├── core/                    # Pure business logic
│   │   ├── services/                # Orchestrators
│   │   └── adapters/http/           # Chi router + handlers
│   │       ├── router.go
│   │       └── handlers/            # Thin HTTP adapters
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/              # Pure UI components
│   │   ├── hooks/                   # Custom hooks (data fetch, side effects)
│   │   ├── services/                # API fetch wrappers
│   │   ├── types/                   # TS interfaces (mirror Go structs)
│   │   ├── themes/                  # Theme registry & definitions
│   │   ├── styles/                  # Generated theme CSS
│   │   ├── index.css                # Global theme-aware styles
│   │   ├── App.tsx
│   │   └── main.tsx
│   ├── scripts/
│   │   └── generate-themes.ts       # Build-time CSS generator
│   ├── Dockerfile
│   ├── package.json
│   ├── bun.lock
│   └── vite.config.ts
├── docker-compose.yml
└── .gitignore
```

---

## Phases

1. **Phase 1 (Current)**: Hello-world scaffold. Go API with `GET /api/v1/hello`. React fetches and displays it. Docker Compose runs both.
2. **Phase 2**: PostgreSQL + domain model (`Delivery`). CRUD + stats endpoints.
3. **Phase 3**: CSV/JSON data ingestion endpoint.
4. **Phase 4**: Charts and trend visualization.

---

## Commands

```bash
# Run everything
docker-compose up --build

# Backend only (inside container or locally)
go run cmd/api/main.go

# Frontend only (Bun)
bun install
bun run dev        # Vite HMR via Bun

# Regenerate theme CSS (run when adding/modifying themes)
bun run generate:themes
```

- Backend port: `:8080`
- Frontend port: `:5173`
- Frontend → Backend in Docker: `http://backend:8080`

---

## Conventions

- **No auth in Phase 1.** Keep scaffold lean.
- **No hot-reload tool for Go** (no `air`). Use `go run` inside Docker for simplicity.
- **Bun replaces npm** on the frontend. Use `bun install`, `bun run`, `bun.lockb`.
- **Vite proxy** (recommended): forward `/api` to backend in `vite.config.ts` to avoid CORS.
- **Go errors are values**: every function returns `(T, error)`. Check `err != nil` immediately.
- **Go interfaces are implicit**: structural typing, like TypeScript.
- **Frontend fetch**: prefer TanStack Query over raw `useEffect` + `fetch`.

---

## Notes for Future Agents

- The developer knows Bun/TypeScript/Node but is learning Go. Draw parallels when explaining Go idioms (interfaces, error handling, goroutines).
- Keep backend handlers thin. Business logic belongs in `core/` or `services/`.
- Keep React components pure. Side effects live in hooks or event handlers, never during render.
- When adding PostgreSQL (Phase 2), use `pgx` driver + `golang-migrate` for migrations.
- If adding auth later, evaluate JWT middleware in chi; do not add to Phase 1.
