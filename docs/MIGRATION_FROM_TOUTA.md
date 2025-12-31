# Migration Guide: From Toutā Framework to Cosan Router

This guide helps you migrate existing Toutā applications to use Cosan Router.

## Why Migrate?

- **SOLID Principles**: Learn and apply architectural best practices
- **Better Testing**: Interface-first design improves testability
- **Performance**: Optimized routing with radix tree
- **Independence**: Zero external dependencies (stdlib only)
- **Rich Features**: Hooks, introspection, metadata

## Migration Strategy

### Option 1: Gradual Migration (Recommended)

Keep your existing Toutā app, replace router incrementally:

```go
// Before
app := touta.New()
app.GET("/users", ListUsers)
app.POST("/users", CreateUser)

// After (Gradual)
app := touta.New()
router := cosan.New()
app.SetRouter(router) // If Toutā supports custom routers

router.GET("/users", ListUsers)
router.POST("/users", CreateUser)
```

### Option 2: Complete Migration

Move entirely to Cosan, integrate Toutā components as needed:

```go
// New approach - Cosan first
router := cosan.New()

// Add only the Toutā components you need
router.Use(middleware.Logger())

// Your routes
router.GET("/users", ListUsers)
router.POST("/users", CreateUser)

// Start server
router.Listen(":8080")
```

## Step-by-Step Migration

### Step 1: Analyze Your Current Routes

List all routes in your Toutā application:

```bash
# Find all route registrations
grep -r "app\.GET\|app\.POST\|app\.PUT" .
```

Create a migration checklist:
- [ ] GET /users
- [ ] POST /users
- [ ] GET /users/:id
- [ ] PUT /users/:id
- [ ] DELETE /users/:id

### Step 2: Install Cosan

```bash
go get github.com/toutaio/toutago-cosan-router/pkg/cosan
```

### Step 3: Create Parallel Router

Set up Cosan alongside existing routes:

```go
package main

import (
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
    "github.com/toutaio/touta-framework/pkg/touta"
)

func main() {
    // Existing Toutā app
    app := touta.New()
    
    // New Cosan router
    router := cosan.New()
    
    // Gradually move routes from app to router
    
    // Option A: Mount Cosan under a prefix
    app.Mount("/api/v2", router)
    
    // Option B: Full replacement (when ready)
    // app.SetRouter(router)
    
    app.Run(":8080")
}
```

### Step 4: Migrate Routes One by One

#### Before (Toutā)
```go
app.GET("/users/:id", func(c *touta.Context) error {
    id := c.Param("id")
    user := userService.Get(id)
    return c.JSON(200, user)
})
```

#### After (Cosan)
```go
router.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    user := userService.Get(id)
    return ctx.JSON(200, user)
})
```

### Step 5: Migrate Middleware

#### Before (Toutā)
```go
app.Use(touta.Logger())
app.Use(touta.Recovery())
app.Use(touta.CORS())
```

#### After (Cosan)
```go
import "github.com/toutaio/toutago-cosan-router/pkg/middleware"

router.Use(middleware.Logger())
router.Use(middleware.Recovery())
router.Use(middleware.CORS())
```

### Step 6: Migrate Route Groups

#### Before (Toutā)
```go
api := app.Group("/api")
v1 := api.Group("/v1")
v1.GET("/users", ListUsers)
```

#### After (Cosan)
```go
api := router.Group("/api")
v1 := api.Group("/v1")
v1.GET("/users", ListUsers)
```

## Feature Mapping

### Context Methods

| Toutā | Cosan | Notes |
|-------|-------|-------|
| `c.Param("id")` | `ctx.Param("id")` | ✅ Same |
| `c.Query("q")` | `ctx.Query("q")` | ✅ Same |
| `c.Bind(&user)` | `ctx.Bind(&user)` | ✅ Same |
| `c.JSON(200, data)` | `ctx.JSON(200, data)` | ✅ Same |
| `c.String(200, "OK")` | `ctx.String(200, "OK")` | ✅ Same |
| `c.HTML(200, html)` | `ctx.HTML(200, html)` | ✅ Same |
| `c.Request()` | `ctx.Request()` | ✅ Same |
| `c.Response()` | `ctx.Response()` | ✅ Same |
| `c.Set("key", val)` | `ctx.Set("key", val)` | ✅ Same |
| `c.Get("key")` | `ctx.Get("key")` | ✅ Same |

### Router Methods

| Toutā | Cosan | Notes |
|-------|-------|-------|
| `app.GET(path, h)` | `router.GET(path, h)` | ✅ Same |
| `app.POST(path, h)` | `router.POST(path, h)` | ✅ Same |
| `app.PUT(path, h)` | `router.PUT(path, h)` | ✅ Same |
| `app.DELETE(path, h)` | `router.DELETE(path, h)` | ✅ Same |
| `app.PATCH(path, h)` | `router.PATCH(path, h)` | ✅ Same |
| `app.Use(mw)` | `router.Use(mw)` | ✅ Same |
| `app.Group(prefix)` | `router.Group(prefix)` | ✅ Same |
| `app.Listen(addr)` | `router.Listen(addr)` | ✅ Same |

### New Features in Cosan

Features not available in Toutā:

```go
// Request/Response hooks
router.BeforeRequest(func(req *http.Request) error {
    // Validate, authenticate, log
    return nil
})

router.AfterResponse(func(req *http.Request, statusCode int) {
    // Log, metrics, cleanup
})

// Custom error handling
router.SetErrorHandler(func(ctx cosan.Context, err error) {
    // Custom error transformation
    ctx.JSON(500, map[string]string{"error": err.Error()})
})

// Route introspection
routes := router.GetRoutes()
for _, r := range routes {
    fmt.Printf("%s %s\n", r.Method, r.Pattern)
}

// Route metadata
router.GET("/users", ListUsers, 
    cosan.WithName("list-users"),
    cosan.WithDescription("Lists all users"),
    cosan.WithTags("users", "api"),
)
```

## Common Patterns

### Pattern 1: Service Layer Integration

#### Before (Toutā with DI)
```go
app.GET("/users", func(c *touta.Context) error {
    userService := c.Get("userService").(*UserService)
    users := userService.List()
    return c.JSON(200, users)
})
```

#### After (Cosan with middleware)
```go
// Setup middleware to inject services
router.Use(cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        ctx.Set("userService", userService)
        return next(ctx)
    }
}))

router.GET("/users", func(ctx cosan.Context) error {
    userService := ctx.Get("userService").(*UserService)
    users := userService.List()
    return ctx.JSON(200, users)
})
```

### Pattern 2: Authentication

#### Before (Toutā)
```go
app.Use(touta.Auth())

app.GET("/protected", func(c *touta.Context) error {
    user := c.Get("user").(*User)
    return c.JSON(200, user)
})
```

#### After (Cosan)
```go
authMiddleware := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        token := ctx.Request().Header.Get("Authorization")
        user, err := authService.Verify(token)
        if err != nil {
            return ctx.JSON(401, map[string]string{"error": "unauthorized"})
        }
        ctx.Set("user", user)
        return next(ctx)
    }
})

router.Use(authMiddleware)

router.GET("/protected", func(ctx cosan.Context) error {
    user := ctx.Get("user").(*User)
    return ctx.JSON(200, user)
})
```

### Pattern 3: Error Handling

#### Before (Toutā)
```go
app.SetErrorHandler(func(c *touta.Context, err error) {
    c.JSON(500, touta.Error{Message: err.Error()})
})
```

#### After (Cosan)
```go
router.SetErrorHandler(func(ctx cosan.Context, err error) {
    ctx.JSON(500, map[string]string{"error": err.Error()})
})
```

## Testing Changes

### Before (Toutā)
```go
func TestListUsers(t *testing.T) {
    app := touta.New()
    app.GET("/users", ListUsers)
    
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()
    app.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### After (Cosan)
```go
func TestListUsers(t *testing.T) {
    router := cosan.New()
    router.GET("/users", ListUsers)
    
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

## Breaking Changes

### Minimal Breaking Changes

Cosan intentionally maintains API compatibility with common patterns:

✅ **Compatible:**
- Handler signature: `func(Context) error`
- Method names: GET, POST, PUT, DELETE, etc.
- Context methods: Param, Query, JSON, etc.
- Middleware pattern
- Route grouping

⚠️ **May Need Changes:**
- Custom middleware (adapt to `cosan.Middleware` interface)
- Framework-specific features (replace with Cosan equivalents)
- Error handling (adapt to new error handler signature)

## Rollback Plan

If migration issues arise:

1. **Keep both routers** during transition
2. **Feature flag** new routes:
   ```go
   if useNewRouter {
       router.GET("/users", ListUsers)
   } else {
       app.GET("/users", ListUsers)
   }
   ```
3. **A/B test** with traffic splitting
4. **Monitor metrics** before/after

## Performance Considerations

Cosan generally performs better due to:
- Radix tree routing
- Context pooling
- Optimized parameter extraction

Run benchmarks before/after:

```bash
# Before migration
go test -bench=. -benchmem ./... > before.txt

# After migration  
go test -bench=. -benchmem ./... > after.txt

# Compare
benchstat before.txt after.txt
```

## Checklist

- [ ] Inventory all routes
- [ ] List all middleware
- [ ] Identify Toutā-specific features
- [ ] Install Cosan
- [ ] Create parallel router
- [ ] Migrate routes incrementally
- [ ] Migrate middleware
- [ ] Update tests
- [ ] Run benchmarks
- [ ] Update documentation
- [ ] Deploy gradually
- [ ] Monitor metrics

## Getting Help

- **Documentation**: [Cosan README](../README.md)
- **Examples**: [examples/](../examples/)
- **Issues**: [GitHub Issues](https://github.com/toutaio/toutago-cosan-router/issues)
- **Discussions**: [GitHub Discussions](https://github.com/toutaio/toutago-cosan-router/discussions)

## Success Stories

After migration, teams typically see:
- ✅ Better test coverage (interfaces are easier to mock)
- ✅ Clearer architecture (SOLID principles)
- ✅ Improved performance (optimized routing)
- ✅ Less coupling (zero dependencies)
- ✅ Easier maintenance (simpler codebase)

## Conclusion

Migrating from Toutā to Cosan is straightforward due to API compatibility. The process can be gradual, allowing you to validate each step. Start with a small set of routes, gain confidence, then migrate the rest.

Happy migrating! 🚀
