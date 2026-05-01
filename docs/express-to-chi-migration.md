# Express to Chi: Router & Middleware Comparison

Mapping Express.js patterns to Go's `chi` router.

---

## Official Resources

| Resource | Link | Purpose |
|----------|------|---------|
| **chi Documentation** | https://go-chi.io/ | Router features and examples |
| **chi GitHub** | https://github.com/go-chi/chi | Source code and middleware |
| **net/http Documentation** | https://pkg.go.dev/net/http | Go's standard HTTP library |
| **Express Documentation** | https://expressjs.com/ | Reference for comparison |

---

## Mental Model Shift

| Express | Chi | Key Difference |
|---------|-----|----------------|
| `app.get("/path", handler)` | `router.Get("/path", handler)` | Same pattern, different syntax |
| `app.use(middleware)` | `router.Use(middleware)` | Middleware stack works identically |
| `req.params.id` | `chi.URLParam(request, "id")` | URL params extracted differently |
| `req.query.name` | `request.URL.Query().Get("name")` | Query params use stdlib |
| `res.json(data)` | `json.NewEncoder(w).Encode(data)` | Response encoding is manual |
| `next()` | Implicit | No `next()`—middleware calls handler directly |

---

## Basic Route

**Express:**
```typescript
app.get("/api/users/:id", (request, response) => {
  const userId = request.params.id
  const user = getUserById(userId)
  response.json(user)
})
```

**Chi:**
```go
router.Get("/api/users/{id}", func(responseWriter http.ResponseWriter, request *http.Request) {
    userId := chi.URLParam(request, "id")
    user := GetUserById(userId)
    json.NewEncoder(responseWriter).Encode(user)
})
```

**Key takeaway:** Chi uses `{param}` syntax (like Express `:param`). URL params are extracted via `chi.URLParam()`.

---

## Middleware Stack

**Express:**
```typescript
app.use(cors())
app.use(express.json())
app.use(logger)

app.get("/api/users", (request, response) => {
  response.json(users)
})
```

**Chi:**
```go
router := chi.NewRouter()

router.Use(middleware.CORS)
router.Use(middleware.Logger)
router.Use(middleware.Recoverer)

router.Get("/api/users", func(responseWriter http.ResponseWriter, request *http.Request) {
    json.NewEncoder(responseWriter).Encode(users)
})
```

**Key takeaway:** Middleware order matters in both. Chi's middleware is compatible with `net/http`, so you can use any standard Go middleware.

---

## Route Groups

**Express:**
```typescript
const userRouter = express.Router()
userRouter.get("/", listUsers)
userRouter.post("/", createUser)
app.use("/api/users", userRouter)
```

**Chi:**
```go
userRouter := chi.NewRouter()
userRouter.Get("/", ListUsers)
userRouter.Post("/", CreateUser)

router.Mount("/api/users", userRouter)
```

**Key takeaway:** `router.Mount()` is like `app.use()` for sub-routers. It prefixes all routes with the mount path.

---

## Error Handling

**Express:**
```typescript
app.get("/api/users/:id", async (request, response, next) => {
  try {
    const user = await getUserById(request.params.id)
    if (!user) {
      return response.status(404).json({ error: "User not found" })
    }
    response.json(user)
  } catch (error) {
    next(error)
  }
})

app.use((error, request, response, next) => {
  response.status(500).json({ error: error.message })
})
```

**Chi:**
```go
router.Get("/api/users/{id}", func(responseWriter http.ResponseWriter, request *http.Request) {
    userId := chi.URLParam(request, "id")
    user, getUserError := GetUserById(userId)
    if getUserError != nil {
        http.Error(responseWriter, getUserError.Error(), http.StatusInternalServerError)
        return
    }
    if user == nil {
        http.Error(responseWriter, "User not found", http.StatusNotFound)
        return
    }
    json.NewEncoder(responseWriter).Encode(user)
})
```

**Key takeaway:** Go doesn't have a global error handler like Express. Each handler checks errors and writes responses directly. You can create a helper function to reduce repetition.

---

## Request Parsing

**Express:**
```typescript
app.use(express.json())

app.post("/api/users", (request, response) => {
  const { name, email } = request.body
  // ...
})
```

**Chi:**
```go
router.Post("/api/users", func(responseWriter http.ResponseWriter, request *http.Request) {
    var requestBody UserRequest
    decodeError := json.NewDecoder(request.Body).Decode(&requestBody)
    if decodeError != nil {
        http.Error(responseWriter, "Invalid JSON", http.StatusBadRequest)
        return
    }
    // ...
})
```

**Key takeaway:** Go doesn't have body parsing middleware. You decode JSON directly in the handler using `json.NewDecoder()`.

---

## Quick Reference: Common Patterns

| Pattern | Express | Chi |
|---------|---------|-----|
| GET route | `app.get("/path", handler)` | `router.Get("/path", handler)` |
| POST route | `app.post("/path", handler)` | `router.Post("/path", handler)` |
| URL param | `req.params.id` | `chi.URLParam(request, "id")` |
| Query param | `req.query.name` | `request.URL.Query().Get("name")` |
| JSON response | `res.json(data)` | `json.NewEncoder(w).Encode(data)` |
| Error response | `res.status(500).json({error})` | `http.Error(w, error, 500)` |
| Middleware | `app.use(middleware)` | `router.Use(middleware)` |
| Sub-router | `express.Router()` | `chi.NewRouter()` + `Mount()` |

---

## Chi-Specific Features

### 1. Method-Based Routing
```go
router.Route("/api/users", func(subRouter chi.Router) {
    subRouter.Get("/", ListUsers)
    subRouter.Post("/", CreateUser)
    subRouter.Get("/{id}", GetUser)
    subRouter.Put("/{id}", UpdateUser)
    subRouter.Delete("/{id}", DeleteUser)
})
```

### 2. Built-in Middleware
```go
router.Use(middleware.RequestID)      // Unique ID per request
router.Use(middleware.RealIP)         // Get real client IP
router.Use(middleware.Logger)         // Log requests
router.Use(middleware.Recoverer)      // Catch panics
router.Use(middleware.Timeout(60 * time.Second)) // Request timeout
```

### 3. URL Pattern Matching
```go
router.Get("/articles/{year}/{month}/{day}/{slug}", GetArticle)
// Extract: chi.URLParam(request, "year"), etc.
```

---

## Next Steps

1. Read [chi Documentation](https://go-chi.io/) (features and examples)
2. Study [net/http](https://pkg.go.dev/net/http) (Go's standard HTTP library)
3. Practice building a CRUD API with chi to internalize patterns
4. Explore [chi middleware](https://github.com/go-chi/chi/tree/master/middleware) for common needs
