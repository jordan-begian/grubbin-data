# Grubbin' Data

A containerized full-stack application for analyzing Grubhub delivery data over time. Built as a study project for **Go** and **React** using functional programming patterns.

## Tech Stack

| Layer           | Technology                                    |
| --------------- | --------------------------------------------- |
| Backend         | Go 1.24 + `chi` router + `slog` logging       |
| Database        | PostgreSQL 17 + Atlas (schema as code)      |
| Schema Mgmt     | Atlas (HCL-based, declarative)                |
| Task Runner     | Task (go-task)                                |
| Frontend        | React 19 + Vite + TypeScript + Bun            |
| Data Fetch      | TanStack Query                                |
| Themes          | CSS custom properties + `@catppuccin/palette` |
| DevOps          | Docker + Docker Compose                       |

## Architecture

**Backend**: MVC Architecture with Functional Core / Imperative Shell (FC/IS)

- `internal/core/` — pure functions, zero I/O
- `internal/services/` — orchestrators (core + side effects)
- `internal/models/` — domain types and data structures
- `internal/controllers/` — thin HTTP handlers (MVC controllers)
- `internal/routes/` — Chi router configuration
- `internal/repositories/` — data access layer (Phase 2+)

**Frontend**: Pure Components + Custom Hooks

- Components are pure: props → JSX
- Side effects isolated in hooks
- Theme switching via CSS custom properties

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
# Start PostgreSQL + apply schema + start backend + frontend
task up

# Or start just the database
task up:db
```

### Development Commands

```bash
# Backend only (requires postgres running)
task backend:run

# Frontend only
bun install
task frontend:dev

# Apply database schema changes after editing HCL
task db:apply

# Generate a versioned migration
task db:diff -- migration_name

# Run backend tests
task backend:test

# Stop everything
task down

# Nuclear reset (stop + remove volumes + restart)
task reset

# Regenerate theme CSS
bun run generate:themes
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

**Phase 2** 🔄 — In Progress

- PostgreSQL database added to Docker Compose
- Go model layer established (`User`, `Profile`, `Vehicle`, `Delivery`, `Location`, `Earnings`)
- Atlas HCL schema defined for all tables
- Task runner configured for database + backend + frontend workflows
- Schema management via Atlas (declarative HCL)

## Project Structure

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
│   │   ├── config/                  # Environment config (godotenv)
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
│   │   ├── hooks/                   # Custom hooks
│   │   ├── services/                # API wrappers
│   │   ├── themes/                  # Theme registry
│   │   └── styles/                  # Generated CSS
│   ├── scripts/generate-themes.ts   # Theme CSS generator
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
