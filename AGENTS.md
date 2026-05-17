# AGENTS.md — Grubbin Data

> Repo for a Grubhub delivery data analysis app. Go backend + React (Bun) frontend.
> Developer background: backend-focused, TypeScript/Node/Bun.

---

## Purpose

Accept and analyze Grubhub delivery data over time. Starting as a study project for Go and React with functional programming patterns.

---

## Architecture

- **Backend**: Go 1.25, MVC Architecture with Functional Core / Imperative Shell (FC/IS).
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
| Backend | Go 1.25 |
| Router | `go-chi/chi/v5` |
| Config | `godotenv` |
| Logging | `log/slog` (stdlib) |
| Database | PostgreSQL 17 |
| Schema Mgmt | Atlas (HCL-based, declarative) |
| Task Runner | Task (go-task) |
| Frontend | React 19 + Vite + TypeScript |
| Frontend runtime | Bun (package install + script runner) |
| Frontend data | TanStack Query |
| Frontend themes | CSS custom properties (Catppuccin, TokyoNight) |
| Styling | Tailwind CSS v4 |
| UI Components | shadcn/ui (manual install) + Radix UI |
| DevOps | Docker + Docker Compose |

---

## Monorepo Layout

```
grubbin-data/
├── backend/
│   ├── atlas/                       # Atlas schema-as-code
│   │   ├── schemas/                 # Schema definition files
│   │   │   ├── schema.pg.hcl        # Schema container (public)
│   │   │   ├── users.pg.hcl         # users, profiles, vehicles tables
│   │   │   └── deliveries.pg.hcl    # deliveries, locations, earnings tables
│   │   └── atlas.hcl                # Atlas project config (envs, vars)
│   ├── cmd/api/main.go              # Entry point
│   ├── internal/
│   │   ├── config/                  # Env loader
│   │   ├── controllers/             # MVC controllers (HTTP handlers)
│   │   ├── core/                    # Pure business logic
│   │   ├── models/                  # Domain types and structs
│   │   │   ├── deliveries.go        # Delivery, PickupLocation, DropoffLocation, DeliveryEarnings
│   │   │   ├── greeting.go          # Greeting response
│   │   │   ├── requests.go          # CreateDeliveryRequest, UpdateDeliveryRequest, RegisterUserRequest, LoginRequest
│   │   │   ├── responses.go         # DeliveryResponse, DeliveryListResponse, DeliveryStats, UserResponse
│   │   │   └── user.go              # User, Profile, Vehicle
│   │   ├── repositories/            # Data access layer (Phase 2+)
│   │   ├── routes/                  # Chi router configuration
│   │   └── services/                # Business logic orchestrators
│   ├── Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── components/              # UI components
│   │   │   ├── ui/                  # shadcn/ui components (manual install)
│   │   │   ├── delivery-form.tsx    # Create delivery form
│   │   │   ├── delivery-edit-form.tsx # Edit delivery modal form
│   │   │   ├── delivery-list.tsx    # Delivery list with selection/edit/delete
│   │   │   ├── delivery-stats.tsx   # Stats overview cards
│   │   │   ├── user-sidebar.tsx     # Right-edge user profile sidebar
│   │   │   ├── login-form.tsx       # Login form
│   │   │   ├── register-form.tsx    # Registration form
│   │   │   ├── password-requirements.tsx # Password validation UI
│   │   │   └── ThemeSwitcher.tsx    # Theme selector
│   │   ├── hooks/                   # Custom hooks (data fetch, side effects)
│   │   │   ├── useAuth.tsx          # Auth context + login/register/logout
│   │   │   ├── useDeliveries.ts    # TanStack Query for delivery CRUD
│   │   │   └── useTheme.tsx         # Theme context + persistence
│   │   ├── lib/                     # Utilities (cn, password validation)
│   │   ├── services/                # API fetch wrappers
│   │   │   ├── auth.ts              # Auth API calls
│   │   │   └── deliveries.ts        # Delivery CRUD API calls
│   │   ├── types/                   # TS interfaces (mirror Go structs)
│   │   │   ├── auth.ts              # User, Profile, Vehicle types
│   │   │   └── delivery.ts          # Delivery request/response types
│   │   ├── themes/                  # Theme CSS files
│   │   ├── styles/                  # Global styles (index.css)
│   │   ├── App.tsx
│   │   └── main.tsx
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
2. **Phase 2 (Complete)**: PostgreSQL + domain model (`Delivery`). CRUD + stats endpoints. Atlas schema management. Task runner workflows. Auth endpoints (register/login) with atomic user+profile+vehicle creation. Delivery CRUD with partial update (PATCH). Database seeding for local dev. Hurl integration tests.
3. **Phase 3 (In Progress)**: CSV/JSON data ingestion endpoint.
4. **Phase 4**: Charts and trend visualization.

---

## Commands

### Task Runner (Primary)

```bash
# Start all services in Docker (postgres + backend + frontend)
task up

# Start only database
task up:db

# Docker DB + local backend + local frontend
task up:local

# Run backend locally (requires postgres running)
task backend:run

# Run backend tests
task backend:test

# Apply database schema changes after editing HCL
task db:apply

# Seed database with test user (local dev only)
task db:seed

# Generate a versioned migration
task db:diff -- migration_name

# Check migration status
task db:status

# Run frontend dev server
task frontend:dev

# Follow Docker logs
task logs

# Stop all Docker services
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

- **No hot-reload tool for Go** (no `air`). Use `go run` inside Docker for simplicity.
- **Bun replaces npm** on the frontend. Use `bun install`, `bun run`, `bun.lock`.
- **Vite proxy** (recommended): forward `/api` to backend in `vite.config.ts` to avoid CORS.
- **Go errors are values**: every function returns `(T, error)`. Check `err != nil` immediately.
- **Go interfaces are implicit**: structural typing, like TypeScript.
- **Frontend fetch**: prefer TanStack Query over raw `useEffect` + `fetch`.
- **shadcn/ui**: Manual installation (copy-paste into `src/components/ui/`). Not installed via npm package.
- **Database schema**: edit HCL in `backend/atlas/schemas/` → run `task db:apply` → Atlas handles SQL generation.
- **Money as cents**: all monetary values stored as `integer` (cents) to avoid float errors.
- **Time as UTC**: all timestamps use `timestamptz` in PostgreSQL; Go normalizes to UTC before storage.

---

## Notes for Future Agents

- The developer knows Bun/TypeScript/Node but is learning Go. Draw parallels when explaining Go idioms (interfaces, error handling, goroutines).
- Keep backend handlers thin. Business logic belongs in `core/` or `services/`.
- Keep React components pure. Side effects live in hooks or event handlers, never during render.
- **Schema changes workflow**: edit `backend/atlas/schemas/*.pg.hcl` → run `task db:apply` (declarative) or `task db:diff -- name` (versioned). Never write raw migration SQL by hand.
- Atlas uses a dev database (`atlas_dev`) for planning diffs. It's created automatically by `task db:apply`.
- Auth is implemented in Phase 2 with bcrypt + pepper password hashing. JWT middleware evaluation is deferred to Phase 3+ if needed.
- The Task runner (`task`) is the primary dev tool. Prefer `task up` over `docker-compose up --build`.
- **Database seeding**: `task db:seed` creates a test user (`test` / `T3stPass135!`) with profile and vehicle for local development. Only runs if user doesn't exist (idempotent).
- **Hurl tests**: Integration tests in `backend/tests/hurl/` test the full delivery CRUD lifecycle. Run with `task backend:test:hurl` (requires backend running).
- **Delivery endpoints**: `POST/GET/PATCH/DELETE /api/users/{userId}/deliveries`. Date filtering via `?start_date=` and `?end_date=` query params (ISO 8601 UTC).
- **Partial updates**: PATCH accepts an array of `{id, ...fields}`. Only non-nil fields are updated. Pointer fields (`*string`, `*time.Time`) distinguish between "don't update" (omitted) and "update to value" (provided).
- **Frontend date-time picker**: Calendar + Popover pattern from shadcn/ui. Use `captionLayout="dropdown"` for month/year selectors. Include Today/Yesterday preset buttons inside the popover footer. Time input below presets.
- **Frontend dollar inputs**: Use `type="text"` with `inputMode="decimal"` (removes spinner arrows). Wrap in relative container with absolute `$` prefix and `pl-7` padding. `parseFloat()` handles string-to-number conversion on submit.
- **Frontend delete confirmation**: AlertDialog for both bulk and single delete operations. Query params for delete: `?id=1&id=2`.
- **Frontend toast notifications**: sonner Toaster in main.tsx, triggered from TanStack Query mutation callbacks in hooks.
