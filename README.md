# Grubbin' Data

A containerized full-stack application for analyzing Grubhub delivery data over time. Built as a study project for **Go** and **React** using functional programming patterns.

## Tech Stack

| Layer      | Technology                                    |
| ---------- | --------------------------------------------- |
| Backend    | Go 1.24 + `chi` router + `slog` logging       |
| Frontend   | React 19 + Vite + TypeScript + Bun            |
| Data Fetch | TanStack Query                                |
| Themes     | CSS custom properties + `@catppuccin/palette` |
| DevOps     | Docker + Docker Compose                       |

## Architecture

**Backend**: MVC Architecture with Functional Core / Imperative Shell (FC/IS)

- `internal/core/` — pure functions, zero I/O
- `internal/services/` — orchestrators (core + side effects)
- `internal/models/` — domain types and data structures
- `internal/controllers/` — thin HTTP handlers (MVC controllers)
- `internal/routes/` — Chi router configuration

**Frontend**: Pure Components + Custom Hooks

- Components are pure: props → JSX
- Side effects isolated in hooks
- Theme switching via CSS custom properties

## Quick Start

```bash
# Run everything
docker compose up --build

# Backend only
go run cmd/api/main.go

# Frontend only
bun install
bun run dev

# Regenerate theme CSS
bun run generate:themes
```

- **Backend**: `http://localhost:8080`
- **Frontend**: `http://localhost:5173`

## Current Status

**Phase 1** ✅ — Hello World scaffold

- `GET /api/v1/hello` returns greeting with timestamp
- Frontend fetches and displays greeting
- Multi-theme support (Catppuccin, Tokyo Night, Gruvbox)
- Docker Compose orchestration

## Project Structure

```
grubbin-data/
├── backend/
│   ├── cmd/api/main.go              # Entry point
│   ├── internal/
│   │   ├── config/                  # Environment config
│   │   ├── controllers/             # MVC controllers (HTTP handlers)
│   │   ├── core/                    # Pure business logic
│   │   ├── models/                  # Domain types and structs
│   │   ├── repositories/            # Data access layer (Phase 2+)
│   │   ├── routes/                  # Chi router configuration
│   │   └── services/                # Business logic orchestrators
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── components/              # Pure UI components
│   │   ├── hooks/                   # Custom hooks
│   │   ├── services/                # API wrappers
│   │   ├── themes/                  # Theme registry
│   │   └── styles/                  # Generated CSS
│   ├── scripts/generate-themes.ts   # Theme CSS generator
│   └── Dockerfile
└── docker-compose.yml
```

## License

MIT — see [LICENSE](LICENSE) for details.
