# AGENTS.md — Grubbin Data

> Repo for a Grubhub delivery data analysis app. Go backend + React (Bun) frontend.
> Developer background: backend-focused, TypeScript/Node/Bun.

---

## Purpose

Accept and analyze Grubhub delivery data over time. Starting as a study project for Go and React with functional programming patterns.

---

## Architecture

- **Backend**: Go 1.24, MVC Architecture with Functional Core / Imperative Shell (FC/IS).
  - `internal/core/` — pure functions, zero I/O, zero external imports.
  - `internal/services/` — orchestrators: compose core + side effects (DB, time).
  - `internal/models/` — domain types and data structures.
  - `internal/controllers/` — thin HTTP handlers (MVC controllers).
  - `internal/routes/` — Chi router configuration.
  - `internal/repositories/` — DB interfaces + implementations (Phase 2+).

> **Note**: The backend uses MVC terminology (models, controllers, routes) while maintaining FC/IS principles. Controllers are thin HTTP adapters, services orchestrate between pure core logic and side effects, and models define domain types. This provides familiar MVC structure with functional programming benefits.

- **Frontend**: React 19 + Vite + TypeScript, executed via Bun.
  - Pure components (props → JSX, no side effects).
  - Custom hooks isolate data fetching / side effects.
  - TanStack Query for server state.

- **Database**: PostgreSQL 17 with Atlas schema management.
  - Schema defined as HCL in `backend/atlas/`.
  - Declarative workflow: edit HCL → `task db:apply` → Atlas generates SQL.
  - Versioned migrations available via `task db:diff`.

---

## Tech Stack

| Layer | Tech |
|-------|------|
| Backend | Go 1.24 |
| Router | `go-chi/chi/v5` |
| Config | `godotenv` |
| Logging | `log/slog` (stdlib) |
| Database | PostgreSQL 17 |
| Schema Mgmt | Atlas (HCL-based, declarative) |
| Task Runner | Task (go-task) |
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
│   ├── atlas/                       # Atlas schema-as-code
│   │   ├── schema.pg.hcl            # Schema container (public)
│   │   ├── users.pg.hcl             # users, profiles, vehicles tables
│   │   ├── deliveries.pg.hcl        # deliveries, locations, earnings tables
│   │   └── atlas.hcl                # Atlas project config (envs, vars)
│   ├── cmd/api/main.go              # Entry point
│   ├── internal/
│   │   ├── config/                  # Env loader
│   │   ├── controllers/             # MVC controllers (HTTP handlers)
│   │   ├── core/                    # Pure business logic
│   │   ├── models/                  # Domain types and structs
│   │   │   ├── deliveries.go        # Delivery, PickupLocation, DropoffLocation, DeliveryEarnings
│   │   │   ├── greeting.go          # Greeting response
│   │   │   ├── requests.go          # CreateDeliveryRequest
│   │   │   ├── responses.go         # DeliveryResponse, DeliveryListResponse, DeliveryStats
│   │   │   └── user.go              # User, Profile, Vehicle
│   │   ├── repositories/            # Data access layer (Phase 2+)
│   │   ├── routes/                  # Chi router configuration
│   │   └── services/                # Business logic orchestrators
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
├── tasks/                           # Component-based Taskfiles
│   ├── backend.yml                  # Go build, test, run
│   ├── db.yml                     # Atlas schema, migrations, wait
│   ├── docker.yml                 # Compose up/down
│   └── frontend.yml               # Bun install, dev, build
├── docker-compose.yml
├── Taskfile.yml                     # Root task runner (includes tasks/*.yml)
└── .gitignore
```

---

## Phases

1. **Phase 1 (Complete)**: Hello-world scaffold. Go API with `GET /api/v1/hello`. React fetches and displays it. Docker Compose runs both.
2. **Phase 2 (In Progress)**: PostgreSQL + domain model (`Delivery`). CRUD + stats endpoints. Atlas schema management. Task runner workflows.
3. **Phase 3**: CSV/JSON data ingestion endpoint.
4. **Phase 4**: Charts and trend visualization.

---

## Commands

### Task Runner (Primary)

```bash
# Start everything (postgres + schema apply + backend + frontend)
task up

# Start only database
task up:db

# Run backend locally (requires postgres running)
task backend:run

# Run backend tests
task backend:test

# Apply database schema changes after editing HCL
task db:apply

# Generate a versioned migration
task db:diff -- migration_name

# Check migration status
task db:status

# Run frontend dev server
task frontend:dev

# Stop everything
task down

# Nuclear reset (stop + remove volumes + restart)
task reset
```

### Legacy (Direct)

```bash
# Run everything via Docker Compose directly
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
- PostgreSQL port: `:5432`
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
- **Database schema**: edit HCL in `backend/atlas/` → run `task db:apply` → Atlas handles SQL generation.
- **Money as cents**: all monetary values stored as `integer` (cents) to avoid float errors.
- **Time as UTC**: all timestamps use `timestamptz` in PostgreSQL; Go normalizes to UTC before storage.

---

## Notes for Future Agents

- The developer knows Bun/TypeScript/Node but is learning Go. Draw parallels when explaining Go idioms (interfaces, error handling, goroutines).
- Keep backend handlers thin. Business logic belongs in `core/` or `services/`.
- Keep React components pure. Side effects live in hooks or event handlers, never during render.
- **Schema changes workflow**: edit `backend/atlas/*.pg.hcl` → run `task db:apply` (declarative) or `task db:diff -- name` (versioned). Never write raw migration SQL by hand.
- Atlas uses a dev database (`atlas_dev`) for planning diffs. It's created automatically by `task db:apply`.
- If adding auth later, evaluate JWT middleware in chi; do not add to Phase 1.
- The Task runner (`task`) is the primary dev tool. Prefer `task up` over `docker-compose up --build`.
