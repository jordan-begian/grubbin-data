# Grubbin' Data

A containerized full-stack application for analyzing Grubhub delivery data over time. Built as a study project for **Go** and **React** using functional programming patterns.

## Tech Stack

| Layer           | Technology                                    |
| --------------- | --------------------------------------------- |
| Backend         | Go 1.25 + `chi` router + `slog` logging       |
| Database        | PostgreSQL 17 + Atlas (schema as code)      |
| Schema Mgmt     | Atlas (HCL-based, declarative)                |
| Task Runner     | Task (go-task)                                |
| Frontend        | React 19 + Vite + TypeScript + Bun            |
| Data Fetch      | TanStack Query                                |
| UI Components   | shadcn/ui (manual install) + Radix UI primitives |
| Themes          | CSS custom properties (Catppuccin, TokyoNight) |
| DevOps          | Docker + Docker Compose                       |

## Architecture

**Backend**: MVC Architecture with Functional Core / Imperative Shell (FC/IS)

- `internal/core/` — pure functions, zero I/O
- `internal/services/` — orchestrators (core + side effects)
- `internal/models/` — domain types and data structures
- `internal/controllers/` — thin HTTP handlers (MVC controllers)
- `internal/routes/` — Chi router configuration
- `internal/repositories/` — data access layer (Phase 2+)

**Frontend**: Pure Components + Custom Hooks + shadcn/ui

- Components are pure: props → JSX
- Side effects isolated in hooks
- shadcn/ui components manually installed in `src/components/ui/`
- Tailwind CSS v4 with CSS custom properties for theming

**Database**: Declarative schema management with Atlas

- Schema defined as HCL in `backend/atlas/`
- Atlas generates and applies migrations automatically
- PostgreSQL 17 via Docker Compose

## Quick Start

### Prerequisites

- [Docker](https://docker.com) + Docker Compose
- [Task](https://taskfile.dev) (`brew install go-task`)
- [Atlas](https://atlasgo.io) (`brew install ariga/tap/atlas`)
- [Bun](https://bun.sh) for frontend

### Environment Setup

```bash
# Copy environment template
cp backend/.env.example backend/.env
# Edit backend/.env with your values (defaults work for local dev)
```

### Run Everything

```bash
# Start all services in Docker (postgres + backend + frontend)
task up

# Or start just the database
task up:db

# Or run Docker DB + local backend + local frontend
task up:local
```

### Development Commands

```bash
# All services in Docker (recommended)
task up

# Docker DB + local backend + local frontend
task up:local

# Backend only (requires postgres running)
task backend:run

# Frontend only
task frontend:dev

# Apply database schema changes after editing HCL
task db:apply

# Generate a versioned migration
task db:diff -- migration_name

# Check migration status
task db:status

# Run backend tests
task backend:test

# Follow Docker logs
task logs

# Stop all Docker services
task down

# Nuclear reset (stop + remove volumes + restart)
task reset
```

- **Backend**: `http://localhost:8080`
- **Frontend**: `http://localhost:5173`
- **PostgreSQL**: `localhost:5432`

## Current Status

**Phase 1** ✅ — Hello World scaffold

- `GET /api/v1/hello` returns greeting with timestamp
- Frontend fetches and displays greeting
- Multi-theme support (Catppuccin, Tokyo Night, Gruvbox)
- Docker Compose orchestration

**Phase 2** ✅ — Complete

- PostgreSQL database added to Docker Compose
- Go model layer established (`User`, `Profile`, `Vehicle`, `Delivery`, `Location`, `Earnings`)
- Atlas HCL schema defined for all tables (in `backend/atlas/schemas/`)
- Task runner configured for database + backend + frontend workflows
- Schema management via Atlas (declarative HCL)
- Auth endpoints: `POST /api/v1/auth/register` and `POST /api/v1/auth/login`
- Password validation with real-time frontend feedback
- Unified `Controller` with service interfaces for testability
- ULID generation for sortable, unique IDs
- Atomic user+profile+vehicle creation via transactions
- All services run in Docker via `task up`

## Project Structure

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
│   │   ├── config/                  # Environment config (godotenv)
│   │   ├── controllers/             # MVC controllers (HTTP handlers)
│   │   ├── core/                    # Pure business logic
│   │   ├── models/                  # Domain types and structs
│   │   │   ├── deliveries.go        # Delivery, PickupLocation, DropoffLocation, DeliveryEarnings
│   │   │   ├── greeting.go          # Greeting response
│   │   │   ├── requests.go          # CreateDeliveryRequest, RegisterUserRequest, LoginRequest
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
│   │   │   └── ui/                  # shadcn/ui components (manual install)
│   │   ├── hooks/                   # Custom hooks
│   │   ├── lib/                     # Utilities (cn, password validation)
│   │   ├── services/                # API wrappers
│   │   ├── themes/                  # Theme CSS files
│   │   └── styles/                  # Global styles (index.css)
│   └── Dockerfile
├── tasks/                           # Component-based Taskfiles
│   ├── backend.yml                  # Go build, test, run
│   ├── db.yml                     # Atlas schema, migrations, wait
│   ├── docker.yml                 # Compose up/down
│   └── frontend.yml               # Bun install, dev, build
├── docker-compose.yml
├── Taskfile.yml                     # Root task runner (includes tasks/*.yml)
└── README.md
```

## License

MIT — see [LICENSE](LICENSE) for details.
