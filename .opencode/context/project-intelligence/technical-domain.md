<!-- Context: project-intelligence/technical | Priority: critical | Version: 1.5 | Updated: 2026-05-17 -->

# Technical Domain

**Purpose**: Tech stack, architecture, and coding patterns for Grubbin Data.
**Last Updated**: 2026-05-16
**Update Triggers**: Tech stack changes | New patterns | Architecture decisions

---

## Quick Reference

| Layer | Technology | Role |
|-------|-----------|------|
| Backend | Go 1.25 | API server, business logic |
| Router | `go-chi/chi/v5` | HTTP routing & middleware |
| Config | `godotenv` | `.env` file loading |
| Logging | `log/slog` (stdlib) | Structured JSON logs |
| Database | PostgreSQL 17 | Persistence |
| Schema Mgmt | Atlas (HCL-based) | Declarative schema as code |
| Task Runner | Task (go-task) | Developer workflow automation |
| Frontend | React 19 + Vite + TypeScript | UI layer |
| Runtime | Bun | Package manager & script runner |
| Data Fetch | TanStack Query | Server-state caching |
| Styling | Tailwind CSS v4 | Utility-first CSS |
| UI Components | shadcn/ui (manual) | Accessible UI primitives |
| DevOps | Docker + Docker Compose | Containerization |

---

## Architecture Patterns

**Backend: MVC Architecture with Functional Core / Imperative Shell (FC/IS)**
- `internal/core/` — pure functions, zero I/O, zero external imports.
- `internal/services/` — orchestrators: compose core + side effects (DB, time).
- `internal/models/` — domain types and data structures.
- `internal/controllers/` — thin HTTP handlers (MVC controllers).
- `internal/routes/` — Chi router configuration.
- `internal/repositories/` — DB interfaces + implementations (Phase 2+).

> **Note**: The backend uses MVC terminology while maintaining FC/IS principles. Controllers are thin HTTP adapters, services orchestrate between pure core logic and side effects.

**Frontend: Pure Components + Custom Hooks**
- Components are pure: props → JSX, no side effects during render.
- Custom hooks isolate data fetching and browser side effects.
- TanStack Query handles server-state caching and background updates.

**Database: Declarative Schema Management**
- Schema defined as HCL in `backend/atlas/`.
- Atlas compares desired state (HCL) to actual database, generates migrations.
- Two workflows: declarative (`task db:apply`) and versioned (`task db:diff`).
- Dev database (`atlas_dev`) used for planning diffs.

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
| DB tables | snake_case, plural | `delivery_records`, `pickup_locations` |
| Atlas HCL | snake_case for resources | `table "users"`, `column "created_at"` |
| Taskfiles | kebab-case, domain-based | `backend.yml`, `db.yml`, `docker.yml` |

**Avoid**: single-letter vars (except `i` in tight loops), abbreviations, Hungarian notation, type-encoded names (`strName`, `intCount`).

---

## API Pattern (Go)

```go
func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
    var req models.RegisterUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, err)
        return
    }

    user, err := c.authService.RegisterUser(r.Context(), req.Username, req.Password, req.FirstName, req.LastName, req.VehicleName, req.VehicleMPG)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, err)
        return
    }

    respondWithJSON(w, http.StatusCreated, user)
}
```

**Rules**: 
- Controllers are thin — parse request, call service, encode response. Never put business logic in controllers.
- Use a unified `Controller` struct with service interfaces (`AuthService`, `GreetingService`, `DeliveryService`) for testability.
- Services are injected via interfaces; mocks are used in controller tests.
- Name variables after their domain meaning.
- Delivery endpoints: `POST/GET/PATCH/DELETE /api/users/{userId}/deliveries` with date filtering via query params.

---

## Component Pattern (React)

```tsx
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

interface UserProfileProps {
  username: string
  fullName: string
  memberSince: string
}

export function UserProfile({ username, fullName, memberSince }: UserProfileProps) {
  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Welcome, {fullName}</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex justify-between">
          <span className="text-muted-foreground">Username</span>
          <span className="font-medium">{username}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-muted-foreground">Member Since</span>
          <span className="font-medium">{memberSince}</span>
        </div>
        <Button variant="outline" className="w-full">Sign Out</Button>
      </CardContent>
    </Card>
  )
}
```

**Rules**: Props in, JSX out. Side effects live in hooks, never during render. Compose shadcn/ui components with Tailwind utility classes. Name props to mirror the backend domain model.

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
| DB tables | snake_case, plural | `delivery_records` |
| Atlas HCL files | snake_case, `.pg.hcl` suffix | `users.pg.hcl`, `deliveries.pg.hcl` |
| Tailwind classes | kebab-case, semantic | `text-muted-foreground`, `bg-primary` |
| Taskfiles | kebab-case, domain-based | `backend.yml`, `db.yml` |
| Env vars | UPPER_SNAKE_CASE | `DATABASE_URL`, `DB_HOST` |

---

## Code Standards

- **Go**: Every function returns `(result, error)`. Check `error != nil` immediately.
- **Go interfaces are implicit**: structural typing, like TypeScript.
- **No auth in Phase 1**. Keep scaffold lean.
- **No hot-reload tool for Go** (no `air`). Use `go run` inside Docker.
- **Bun replaces npm**: use `bun install`, `bun run`, `bun.lock`.
- **Vite proxy**: forward `/api` to backend in `vite.config.ts` to avoid CORS.
- **Frontend fetch**: prefer TanStack Query over raw `useEffect` + `fetch`.
- **Self-documenting names**: prefer `delivery.EarningsPerHour` over `d.Eph` + comment.
- **shadcn/ui**: Manual installation (copy-paste into `src/components/ui/`). Not installed via npm package.
- **Theme system**: CSS custom properties defined in `src/styles/index.css` via `@theme inline` blocks.
- **Database schema**: edit HCL in `backend/atlas/schemas/` → run `task db:apply` → Atlas handles SQL generation.
- **Date-time picker pattern**: Calendar (shadcn/ui) inside Popover with `captionLayout="dropdown"` for month/year selectors. Preset buttons (Today/Yesterday) above time input in popover footer. Trigger button shows combined date + time string.
- **Dollar input pattern**: `type="text"` with `inputMode="decimal"` (no spinner arrows). Relative container with absolute `$` prefix span and `pl-7` padding on input. `parseFloat()` for conversion on submit.
- **Money as cents**: all monetary values stored as `integer` (cents) to avoid float errors.
- **Time as UTC**: all timestamps use `timestamptz` in PostgreSQL; Go normalizes to UTC before storage.
- **Task runner**: `task` is the primary dev tool. Prefer `task up` over `docker-compose up --build`.
- **Database seeding**: `task db:seed` creates test user for local dev (idempotent).
- **Integration tests**: Hurl tests in `backend/tests/hurl/` test full API lifecycle.

---

## Security Requirements

- Validate all user input with `go-playground/validator` (Phase 2+).
- Use parameterized queries via `pgx` (Phase 2+).
- Auth implemented in Phase 2: register/login with password validation.
- Passwords hashed with bcrypt + pepper; never return password in responses.
- Frontend password requirements mirror backend validation rules.
- CORS handled via Vite proxy in dev; chi CORS middleware if needed.

---

## 📂 Codebase References

**Backend entry**: `backend/cmd/api/main.go` — wires config, router, server.
**Seed command**: `backend/cmd/seed/main.go` — seeds test user for local development.
**Frontend entry**: `frontend/src/main.tsx` — React root + QueryClientProvider.
**API client**: `frontend/src/services/api.ts` — pure fetch wrapper.
**Config**: `backend/internal/config/config.go` — typed env loader.
**UI components**: `frontend/src/components/ui/` — shadcn/ui components (Button, Card, Input, Label, PasswordInput, Calendar, Popover, Dialog, AlertDialog, DropdownMenu, ScrollArea, Skeleton, Checkbox, Badge, Avatar, Separator, Textarea).
**Hurl tests**: `backend/tests/hurl/` — integration tests for delivery CRUD.

**Models**:
- `backend/internal/models/user.go` — `User`, `Profile`, `Vehicle`
- `backend/internal/models/deliveries.go` — `Delivery`, `PickupLocation`, `DropoffLocation`, `DeliveryEarnings`
- `backend/internal/models/requests.go` — `CreateDeliveryRequest`, `UpdateDeliveryRequest`, `RegisterUserRequest`, `LoginRequest`
- `backend/internal/models/responses.go` — `DeliveryResponse`, `DeliveryListResponse`, `DeliveryStats`, `UserResponse`

**Atlas Schema**:
- `backend/atlas/schemas/schema.pg.hcl` — Schema container (`public`)
- `backend/atlas/schemas/users.pg.hcl` — `users`, `profiles`, `vehicles` tables
- `backend/atlas/schemas/deliveries.pg.hcl` — `deliveries`, `pickup_locations`, `dropoff_locations`, `earnings` tables
- `backend/atlas/atlas.hcl` — Atlas project config (environments, variables)

**Services**:
- `backend/internal/services/auth_service.go` — `AuthService` (register, login)
- `backend/internal/services/delivery_service.go` — `DeliveryService` (CRUD + stats)
- `backend/internal/services/greeting_service.go` — `GreetingService`

**Task Runner**:
- `Taskfile.yml` — Root taskfile with shared vars and includes
- `tasks/backend.yml` — Go build, test, run, hurl tests
- `tasks/db.yml` — Atlas schema apply, migrations, seed, wait
- `tasks/docker.yml` — Docker Compose up/down
- `tasks/frontend.yml` — Bun install, dev, build

**Theme hook**: `frontend/src/hooks/useTheme.tsx` — context provider + persistence.
**Theme UI**: `frontend/src/components/ThemeSwitcher.tsx` — dropdown selector.
**Global styles**: `frontend/src/styles/index.css` — theme-aware body/html styles with `@theme inline` blocks.
**shadcn utils**: `frontend/src/lib/utils.ts` — `cn()` helper (clsx + tailwind-merge).
**Password validation**: `frontend/src/lib/password-validation.ts` — mirrors backend password rules.

**Delivery Dashboard Components**:
- `frontend/src/components/delivery-form.tsx` — Create delivery form with Calendar + Popover date-time picker, `$` prefixed dollar inputs, form validation.
- `frontend/src/components/delivery-edit-form.tsx` — Edit delivery modal form with same date-time picker pattern, partial update logic.
- `frontend/src/components/delivery-list.tsx` — Delivery list with checkbox selection, inline edit/delete actions, AlertDialog confirmations, ScrollArea for overflow.
- `frontend/src/components/delivery-stats.tsx` — Stats overview cards with Skeleton loading states.
- `frontend/src/components/user-sidebar.tsx` — Right-edge collapsible sidebar with Avatar, profile/vehicle info, logout button.

**Delivery Hooks**:
- `frontend/src/hooks/useDeliveries.ts` — TanStack Query hook with mutations for create/update/delete, toast notifications on success/error, query invalidation.

---

## Related Files

- `AGENTS.md` — Project overview, conventions, commands.
- `README.md` — Quick start, project structure, usage.
- `docker-compose.yml` — Service orchestration.
