# Official Resources

Curated links to verified documentation for Go and React.

---

## Go

### Core Language
| Resource | Link | Description |
|----------|------|-------------|
| **A Tour of Go** | https://go.dev/tour/ | Interactive introduction (start here) |
| **Effective Go** | https://go.dev/doc/effective_go | Idiomatic Go patterns and conventions |
| **Go Documentation** | https://go.dev/doc/ | Language spec, release notes, guides |
| **Go by Example** | https://gobyexample.com/ | Quick reference with runnable examples |
| **Go Playground** | https://go.dev/play/ | Online code editor for testing |

### Standard Library
| Resource | Link | Description |
|----------|------|-------------|
| **net/http** | https://pkg.go.dev/net/http | HTTP client and server |
| **encoding/json** | https://pkg.go.dev/encoding/json | JSON encoding/decoding |
| **log/slog** | https://pkg.go.dev/log/slog | Structured logging (Go 1.21+) |
| **context** | https://pkg.go.dev/context | Request-scoped values and cancellation |
| **sync** | https://pkg.go.dev/sync | Concurrency primitives |

### Ecosystem
| Resource | Link | Description |
|----------|------|-------------|
| **Go Modules** | https://go.dev/ref/mod | Dependency management |
| **chi Router** | https://go-chi.io/ | Lightweight HTTP router |
| **godotenv** | https://github.com/joho/godotenv | `.env` file loading |
| **pgx** | https://github.com/jackc/pgx | PostgreSQL driver |
| **golang-migrate** | https://github.com/golang-migrate/migrate | Database migrations |

---

## React

### Core Documentation
| Resource | Link | Description |
|----------|------|-------------|
| **React Documentation** | https://react.dev/ | Complete reference (start here) |
| **Thinking in React** | https://react.dev/learn/thinking-in-react | Mental model for React |
| **Components & Props** | https://react.dev/learn/passing-props-to-a-component | Data flow fundamentals |
| **State Management** | https://react.dev/learn/managing-state | When and how to use state |
| **Hooks Reference** | https://react.dev/reference/react/hooks | All built-in hooks |

### Key Concepts
| Resource | Link | Description |
|----------|------|-------------|
| **useState** | https://react.dev/reference/react/useState | Local component state |
| **useEffect** | https://react.dev/reference/react/useEffect | Side effects after render |
| **useContext** | https://react.dev/reference/react/useContext | Share data without prop drilling |
| **Custom Hooks** | https://react.dev/learn/reusing-logic-with-custom-hooks | Reusable stateful logic |
| **Pure Components** | https://react.dev/learn/keeping-components-pure | Components as pure functions |

### Ecosystem
| Resource | Link | Description |
|----------|------|-------------|
| **TanStack Query** | https://tanstack.com/query/latest | Server-state fetching and caching |
| **React Router** | https://reactrouter.com/ | Client-side routing |
| **Vite** | https://vitejs.dev/ | Build tool and dev server |
| **TypeScript** | https://www.typescriptlang.org/docs/ | Type-safe JavaScript |
| **Bun** | https://bun.sh/docs | Fast JavaScript runtime |

---

## Learning Path

### Week 1: Go Fundamentals
1. Complete [A Tour of Go](https://go.dev/tour/) (interactive, ~2 hours)
2. Read [Effective Go](https://go.dev/doc/effective_go) (idioms and conventions)
3. Practice with [Go by Example](https://gobyexample.com/) (quick reference)
4. Build a small CLI tool to internalize error handling patterns

### Week 2: Go Web Development
1. Study [net/http](https://pkg.go.dev/net/http) (standard library)
2. Read [chi Documentation](https://go-chi.io/) (router features)
3. Build a CRUD API with chi to practice routing and middleware
4. Add structured logging with `slog`

### Week 3: React Fundamentals
1. Read [Thinking in React](https://react.dev/learn/thinking-in-react) (mental model)
2. Complete [React Tutorial](https://react.dev/learn/tutorial-tic-tac-toe) (interactive)
3. Study [Hooks Reference](https://react.dev/reference/react/hooks) (API reference)
4. Build a small form with validation to practice state management

### Week 4: React + Go Integration
1. Set up TanStack Query for data fetching
2. Connect React frontend to Go backend
3. Add theme switching with CSS custom properties
4. Deploy with Docker Compose

---

## Quick Tips

- **Go**: Errors are values. Check `err != nil` immediately. No exceptions.
- **React**: Components are pure functions. Side effects live in hooks or event handlers.
- **Both**: Favor composition over inheritance. Keep things simple.
- **Docker**: Use multi-stage builds for smaller images. Mount volumes for hot reload.

---

## Community

| Platform | Go | React |
|----------|----|----|
| **Official Forum** | https://forum.golangbridge.org/ | https://react.dev/community |
| **Reddit** | r/golang | r/reactjs |
| **Discord** | https://discord.gg/golang | https://discord.gg/reactiflux |
| **Stack Overflow** | [go] tag | [reactjs] tag |
| **GitHub** | https://github.com/golang/go | https://github.com/facebook/react |
