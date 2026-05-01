# Go for TypeScript Developers

A practical guide mapping TypeScript/Node concepts to Go idioms.

---

## Official Resources

| Resource | Link | Purpose |
|----------|------|---------|
| **A Tour of Go** | https://go.dev/tour/ | Interactive introduction (start here) |
| **Effective Go** | https://go.dev/doc/effective_go | Idiomatic Go patterns |
| **Go Documentation** | https://go.dev/doc/ | Language spec, standard library |
| **Go by Example** | https://gobyexample.com/ | Quick reference with examples |
| **Go Modules** | https://go.dev/ref/mod | Dependency management |

---

## Mental Model Shift

| TypeScript/Node | Go | Key Difference |
|-----------------|----|----------------|
| Async/await + Promises | Goroutines + channels | Go schedules goroutines; no `await` keyword |
| `throw` / `try/catch` | `return (T, error)` | Errors are values, not exceptions |
| `interface` (structural) | `interface` (implicit) | Same concept, Go enforces at compile time |
| `class` + `this` | `struct` + receiver methods | No inheritance, composition only |
| `npm` + `package.json` | `go mod` + `go.mod` | Built-in, no external tool needed |
| `console.log` | `fmt.Println` / `slog` | Structured logging built into stdlib |

---

## Error Handling

**TypeScript:**
```typescript
async function getUser(id: string): Promise<User> {
  const response = await fetch(`/api/users/${id}`)
  if (!response.ok) {
    throw new Error(`Failed to fetch user: ${response.statusText}`)
  }
  return response.json()
}

try {
  const user = await getUser("123")
  console.log(user.name)
} catch (error) {
  console.error("Error:", error.message)
}
```

**Go:**
```go
func GetUser(id string) (User, error) {
    response, requestError := http.Get(fmt.Sprintf("/api/users/%s", id))
    if requestError != nil {
        return User{}, requestError
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        return User{}, fmt.Errorf("failed to fetch user: %s", response.Status)
    }

    var user User
    decodeError := json.NewDecoder(response.Body).Decode(&user)
    if decodeError != nil {
        return User{}, decodeError
    }

    return user, nil
}

// Usage
user, getUserError := GetUser("123")
if getUserError != nil {
    log.Printf("Error: %v", getUserError)
    return
}
fmt.Println(user.Name)
```

**Key takeaway:** Go errors are values you check immediately. No stack traces, no try/catch blocks. The `defer` keyword is like a `finally` block that runs when the function returns.

---

## Interfaces (Structural Typing)

**TypeScript:**
```typescript
interface Logger {
  log(message: string): void
  error(message: string): void
}

class ConsoleLogger implements Logger {
  log(message: string) { console.log(message) }
  error(message: string) { console.error(message) }
}

function process(logger: Logger) {
  logger.log("Processing...")
}
```

**Go:**
```go
type Logger interface {
    Log(message string)
    Error(message string)
}

type ConsoleLogger struct{}

func (l ConsoleLogger) Log(message string)   { fmt.Println(message) }
func (l ConsoleLogger) Error(message string) { fmt.Println(message) }

func Process(logger Logger) {
    logger.Log("Processing...")
}

// No "implements" keyword needed.
// If ConsoleLogger has Log() and Error(), it satisfies Logger.
```

**Key takeaway:** Go interfaces are implicit. If your type has the methods, it satisfies the interface. This is identical to TypeScript's structural typing.

---

## Concurrency

**TypeScript:**
```typescript
async function fetchMultiple(urls: string[]): Promise<string[]> {
  const promises = urls.map(url => fetch(url).then(r => r.text()))
  return Promise.all(promises)
}
```

**Go:**
```go
func FetchMultiple(urls []string) ([]string, error) {
    results := make([]string, len(urls))
    waitGroup := sync.WaitGroup{}

    for index, url := range urls {
        waitGroup.Add(1)
        go func(url string, index int) {
            defer waitGroup.Done()
            response, requestError := http.Get(url)
            if requestError != nil {
                return
            }
            defer response.Body.Close()
            body, readError := io.ReadAll(response.Body)
            if readError != nil {
                return
            }
            results[index] = string(body)
        }(url, index)
    }

    waitGroup.Wait()
    return results, nil
}
```

**Key takeaway:** `go func()` spins up a lightweight thread (goroutine). `sync.WaitGroup` waits for all goroutines to finish. Go's scheduler handles the rest—no event loop, no Promise.all.

---

## Modules & Dependencies

**TypeScript:**
```bash
npm install express
# package.json updated automatically
```

**Go:**
```bash
go get github.com/go-chi/chi/v5
# go.mod and go.sum updated automatically
```

**Key takeaway:** Go modules are built into the language. `go.mod` is like `package.json`, `go.sum` is like `package-lock.json`. No separate package manager needed.

---

## Project Structure

**TypeScript/Express:**
```
src/
├── controllers/
├── middleware/
├── routes/
├── services/
└── app.ts
```

**Go (FC/IS Pattern):**
```
internal/
├── core/           # Pure functions (no I/O)
├── services/       # Orchestrators (core + side effects)
├── adapters/http/  # Thin HTTP handlers
└── repositories/   # DB interfaces + implementations
```

**Key takeaway:** Go's `internal/` directory prevents other packages from importing your code. The FC/IS pattern keeps business logic pure and testable.

---

## Quick Reference: Syntax Differences

| Concept | TypeScript | Go |
|---------|-----------|----|
| Variable declaration | `const x = 5` | `x := 5` |
| Function | `function add(a: number, b: number): number` | `func Add(a int, b int) int` |
| Export | `export function` | Capitalized name: `func Add()` |
| Private | `function _helper()` | Lowercase name: `func helper()` |
| Null/undefined | `null` / `undefined` | `nil` |
| Template string | `` `Hello ${name}` `` | `` fmt.Sprintf("Hello %s", name) `` |
| Array length | `arr.length` | `len(arr)` |
| Map/Dict | `Map<string, number>` | `map[string]int` |
| Destructuring | `const { a, b } = obj` | No direct equivalent; access fields |
| Optional chaining | `obj?.prop?.method()` | Check `nil` explicitly |

---

## Common Pitfalls for TS Developers

1. **Don't use `this`**: Go methods use explicit receivers. No `this` binding issues.
2. **Don't throw errors**: Return them. Check `err != nil` immediately.
3. **Don't mutate slices**: Use `append()` which returns a new slice.
4. **Don't ignore errors**: `_ = someFunction()` hides errors. Handle them.
5. **Don't over-abstract**: Go favors simplicity. Interfaces only when needed.

---

## Next Steps

1. Complete [A Tour of Go](https://go.dev/tour/) (interactive, ~2 hours)
2. Read [Effective Go](https://go.dev/doc/effective_go) (idioms and conventions)
3. Practice with [Go by Example](https://gobyexample.com/) (quick reference)
4. Build a small CLI tool to internalize error handling patterns
