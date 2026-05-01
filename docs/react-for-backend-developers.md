# React for Backend Developers

A practical guide mapping backend concepts to React patterns.

---

## Official Resources

| Resource | Link | Purpose |
|----------|------|---------|
| **React Documentation** | https://react.dev/ | Complete reference (start here) |
| **Thinking in React** | https://react.dev/learn/thinking-in-react | Mental model for React |
| **Hooks Reference** | https://react.dev/reference/react/hooks | All built-in hooks |
| **Components & Props** | https://react.dev/learn/passing-props-to-a-component | Data flow fundamentals |
| **State Management** | https://react.dev/learn/managing-state | When and how to use state |

---

## Mental Model Shift

| Backend Concept | React Equivalent | Key Difference |
|-----------------|------------------|----------------|
| Server renders HTML | Component returns JSX | JSX is compiled to `React.createElement()` calls |
| Request → Response | Props → JSX | Components are pure functions of their inputs |
| Session state | `useState` | State triggers re-render when updated |
| Middleware | Higher-order components / Hooks | Wrap components with shared logic |
| Template engine (EJS/Pug) | JSX | JSX is JavaScript, not a separate templating language |
| Server-side routing | `react-router-dom` | Client-side routing, no page reloads |

---

## Components as Pure Functions

**Backend (Express route handler):**
```typescript
app.get("/user/:id", (request, response) => {
  const user = getUserById(request.params.id)
  response.render("user-profile", { user })
})
```

**React (Component):**
```tsx
interface UserProfileProps {
  userName: string
  userEmail: string
}

function UserProfile({ userName, userEmail }: UserProfileProps) {
  return (
    <article>
      <h1>{userName}</h1>
      <p>{userEmail}</p>
    </article>
  )
}
```

**Key takeaway:** A React component is `f(props) => JSX`. Given the same props, it always returns the same JSX. This is identical to a pure function in your backend code.

---

## State vs Props

**Props** are like function arguments—passed down, read-only.
**State** is like local variables—managed internally, triggers re-render.

```tsx
function Counter() {
  // State: managed internally, triggers re-render
  const [count, setCount] = useState(0)

  return (
    <button onClick={() => setCount(count + 1)}>
      Count: {count}
    </button>
  )
}
```

**Key takeaway:** `useState` returns `[currentValue, setterFunction]`. Calling the setter triggers a re-render with the new value. This is like updating a variable and re-running your template.

---

## Effects (Side Effects)

**Backend:** Side effects happen in route handlers (DB queries, file I/O).

**React:** Side effects happen in `useEffect` or event handlers—never during render.

```tsx
function UserProfile({ userId }: { userId: string }) {
  const [userData, setUserData] = useState<User | null>(null)

  // Runs after render, when userId changes
  useEffect(() => {
    fetchUser(userId).then(setUserData)
  }, [userId]) // Dependency array: re-run when userId changes

  if (!userData) return <p>Loading...</p>

  return <UserProfileDisplay user={userData} />
}
```

**Key takeaway:** `useEffect` is your "async boundary." It runs after render, not during. The dependency array `[userId]` means "re-run this effect when userId changes."

---

## Custom Hooks (Reusable Logic)

**Backend:** Extract shared logic into a service function.

**React:** Extract shared logic into a custom hook.

```tsx
// Custom hook: encapsulates data fetching logic
function useUserData(userId: string) {
  const [userData, setUserData] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    setIsLoading(true)
    fetchUser(userId)
      .then(setUserData)
      .finally(() => setIsLoading(false))
  }, [userId])

  return { userData, isLoading }
}

// Component: stays pure, just renders
function UserProfile({ userId }: { userId: string }) {
  const { userData, isLoading } = useUserData(userId)

  if (isLoading) return <p>Loading...</p>
  return <UserProfileDisplay user={userData} />
}
```

**Key takeaway:** Custom hooks let you share stateful logic between components. They're like service functions that return state + setters.

---

## TanStack Query (Server State)

For data fetching, you rarely need `useEffect` + `fetch`. TanStack Query handles caching, background updates, and loading states.

```tsx
// Hook: declarative data fetching
function useUserData(userId: string) {
  return useQuery({
    queryKey: ["user", userId],
    queryFn: () => fetchUser(userId),
  })
}

// Component: just renders
function UserProfile({ userId }: { userId: string }) {
  const { data: userData, isLoading, isError } = useUserData(userId)

  if (isLoading) return <p>Loading...</p>
  if (isError) return <p>Error loading user</p>

  return <UserProfileDisplay user={userData} />
}
```

**Key takeaway:** TanStack Query replaces manual `useEffect` + `fetch` + caching logic. You declare "this component needs this data" and it handles the rest.

---

## Event Handlers

**Backend:** Route handlers process requests.

**React:** Event handlers process user interactions.

```tsx
function DeleteButton({ userId }: { userId: string }) {
  const handleDelete = async () => {
    await deleteUser(userId)
    // Trigger re-fetch or update local state
  }

  return (
    <button onClick={handleDelete}>
      Delete User
    </button>
  )
}
```

**Key takeaway:** Event handlers are where side effects live. They don't run during render—they run when the user interacts.

---

## Quick Reference: React Patterns

| Pattern | Purpose | Example |
|---------|---------|---------|
| `useState` | Local component state | Form inputs, toggles |
| `useEffect` | Side effects after render | Data fetching, subscriptions |
| `useContext` | Share data without prop drilling | Theme, auth user |
| `useRef` | Hold mutable value without re-render | DOM nodes, timers |
| `useMemo` | Cache expensive calculation | Filtered lists, derived data |
| `useCallback` | Cache function reference | Pass to optimized child components |
| Custom hooks | Reusable stateful logic | `useUserData`, `useTheme` |

---

## Common Pitfalls for Backend Developers

1. **Don't mutate state**: Use `setState(newValue)`, not `state.property = value`.
2. **Don't put side effects in render**: Effects belong in `useEffect` or event handlers.
3. **Don't over-use `useEffect`**: If you're fetching data, use TanStack Query instead.
4. **Don't forget dependency arrays**: Missing deps cause stale closures or infinite loops.
5. **Don't treat components like classes**: They're functions. No `this`, no lifecycle methods.

---

## Next Steps

1. Read [Thinking in React](https://react.dev/learn/thinking-in-react) (mental model)
2. Complete [React Tutorial](https://react.dev/learn/tutorial-tic-tac-toe) (interactive)
3. Study [Hooks Reference](https://react.dev/reference/react/hooks) (API reference)
4. Build a small form with validation to practice state management
