# Migrating from Chi to Cosan

This guide helps you migrate from [go-chi/chi](https://github.com/go-chi/chi) to Cosan.

## Overview

Cosan and Chi share similar philosophies:
- Both build on `net/http`
- Both emphasize simplicity and performance
- Both support middleware composition

**Key Differences:**
- Cosan uses a `Context` interface for request/response handling
- Cosan's middleware signature is different
- Cosan emphasizes SOLID principles and testability
- Cosan supports optional ecosystem integrations

## Quick Comparison

### Router Creation

**Chi:**
```go
r := chi.NewRouter()
```

**Cosan:**
```go
r := cosan.New()
```

### Route Registration

**Chi:**
```go
r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    w.Write([]byte("User: " + id))
})
```

**Cosan:**
```go
r.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    return ctx.String(200, "User: " + id)
})
```

### Middleware

**Chi:**
```go
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
```

**Cosan:**
```go
r.Use(middleware.Logger())
r.Use(middleware.Recovery())
```

### Route Groups

**Chi:**
```go
r.Route("/api", func(r chi.Router) {
    r.Use(authMiddleware)
    r.Get("/users", listUsers)
    r.Post("/users", createUser)
})
```

**Cosan:**
```go
api := r.Group("/api")
api.Use(authMiddleware)
api.GET("/users", listUsers)
api.POST("/users", createUser)
```

## Step-by-Step Migration

### Step 1: Update Imports

**Before:**
```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"
)
```

**After:**
```go
import (
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
    "github.com/toutaio/toutago-cosan-router/pkg/middleware"
)
```

### Step 2: Update Router Initialization

**Before:**
```go
func main() {
    r := chi.NewRouter()
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    
    http.ListenAndServe(":3000", r)
}
```

**After:**
```go
func main() {
    r := cosan.New()
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    
    r.Listen(":3000")
}
```

### Step 3: Convert Handler Signatures

**Before:**
```go
func getUser(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    
    user, err := db.GetUser(id)
    if err != nil {
        http.Error(w, "Not found", 404)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
```

**After:**
```go
func getUser(ctx cosan.Context) error {
    id := ctx.Param("id")
    
    user, err := db.GetUser(id)
    if err != nil {
        return ctx.String(404, "Not found")
    }
    
    return ctx.JSON(200, user)
}
```

### Step 4: Update Middleware

**Before:**
```go
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

**After:**
```go
func authMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        token := ctx.Request().Header.Get("Authorization")
        if token == "" {
            return ctx.String(401, "Unauthorized")
        }
        return next(ctx)
    }
}
```

### Step 5: Update URL Parameters

**Chi:**
```go
id := chi.URLParam(r, "id")
name := chi.URLParam(r, "name")
```

**Cosan:**
```go
id := ctx.Param("id")
name := ctx.Param("name")

// Or get all params
params := ctx.Params()
```

### Step 6: Update Query Parameters

**Chi:**
```go
page := r.URL.Query().Get("page")
limit := r.URL.Query().Get("limit")
```

**Cosan:**
```go
page := ctx.Query("page")
limit := ctx.Query("limit")

// With default value
page := ctx.QueryDefault("page", "1")
```

## Common Patterns

### JSON Response

**Chi:**
```go
w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(data)
```

**Cosan:**
```go
return ctx.JSON(200, data)
```

### Error Handling

**Chi:**
```go
if err != nil {
    http.Error(w, err.Error(), 500)
    return
}
```

**Cosan:**
```go
if err != nil {
    return err // Will be handled by error middleware
}

// Or explicitly:
if err != nil {
    return ctx.String(500, err.Error())
}
```

### Static Files

**Chi:**
```go
r.Handle("/static/*", http.FileServer(http.Dir("./public")))
```

**Cosan:**
```go
r.Static("/static", "./public")
```

### Subrouters

**Chi:**
```go
r.Mount("/api/v1", apiV1Router())

func apiV1Router() chi.Router {
    r := chi.NewRouter()
    r.Get("/users", listUsers)
    return r
}
```

**Cosan:**
```go
api := r.Group("/api/v1")
api.GET("/users", listUsers)

// Or with a separate function
func setupAPIRoutes(r cosan.Router) {
    api := r.Group("/api/v1")
    api.GET("/users", listUsers)
}
```

## Testing Changes

### Chi Tests

**Before:**
```go
func TestHandler(t *testing.T) {
    r := chi.NewRouter()
    r.Get("/test", handler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

**After:**
```go
func TestHandler(t *testing.T) {
    r := cosan.New()
    r.GET("/test", handler)
    
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}

// Or using Cosan's test helpers
func TestHandler(t *testing.T) {
    ctx := cosan.NewTestContext()
    ctx.SetPath("/test")
    
    err := handler(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 200, ctx.StatusCode())
}
```

## Performance Considerations

Both Chi and Cosan have similar performance characteristics:
- Both use radix tree for routing
- Both have minimal allocations
- Both support middleware composition

Cosan may have slightly more overhead due to the Context abstraction, but it's typically <5% and provides significant benefits in testability and flexibility.

## Benefits of Migration

### Better Testability
- Mock the Context interface for unit tests
- No need for httptest.NewRequest in every test
- Easy to test middleware in isolation

### Cleaner Error Handling
- Return errors instead of calling http.Error
- Centralized error handling via middleware
- More composable and testable

### Enhanced Features
- Built-in JSON/String/HTML helpers
- Query parameter helpers with defaults
- Optional ecosystem integrations (datamapper, renderer, DI)

### SOLID Principles
- Better separation of concerns
- More maintainable code
- Easier to extend and customize

## Gradual Migration Strategy

You don't have to migrate everything at once:

1. **Start with new routes**: Use Cosan for new endpoints
2. **Migrate route groups**: Move one group at a time
3. **Proxy between routers**: Use Chi and Cosan side-by-side temporarily
4. **Complete migration**: Remove Chi when all routes migrated

```go
func main() {
    chiRouter := chi.NewRouter()
    cosanRouter := cosan.New()
    
    // Old routes on Chi
    chiRouter.Get("/old-endpoint", oldHandler)
    
    // New routes on Cosan
    cosanRouter.GET("/new-endpoint", newHandler)
    
    // Proxy Cosan through Chi temporarily
    chiRouter.Mount("/api/v2", cosanRouter)
    
    http.ListenAndServe(":8080", chiRouter)
}
```

## Getting Help

- [Cosan Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
- [Examples](../../examples/)
- [GitHub Issues](https://github.com/toutaio/toutago-cosan-router/issues)

## See Also

- [Migration from Gin](./from-gin.md)
- [Migration from Echo](./from-echo.md)
- [Performance Guide](../guides/performance.md)
