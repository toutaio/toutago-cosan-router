# Toutā Ecosystem Integration Guide

This guide explains how to integrate Cosan Router with the Toutā framework and its ecosystem components.

## Overview

Cosan is designed to work **standalone** or as part of the Toutā ecosystem. All integrations are **optional** - Cosan works perfectly without any external dependencies.

## Architecture

```
┌─────────────────────────────────────────┐
│           Application Layer              │
│  (Your Business Logic & Handlers)       │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│         Cosan Router (Core)              │
│      Zero Dependencies - Stdlib Only     │
└─────────────────────────────────────────┘
                    ↓
        ┌───────────┴───────────┐
        ↓                       ↓
┌───────────────┐      ┌───────────────┐
│   Standalone  │  OR  │   Toutā       │
│   Usage       │      │   Ecosystem   │
└───────────────┘      └───────────────┘
```

## Standalone Usage (Recommended for Most Cases)

Cosan works out-of-the-box with zero configuration:

```go
package main

import (
    "log"
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

func main() {
    router := cosan.New()
    
    router.GET("/hello", func(ctx cosan.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello, World!",
        })
    })
    
    log.Fatal(router.Listen(":8080"))
}
```

## Toutā Framework Integration

### As a Toutā Router Provider

If you're using the Toutā framework, Cosan can serve as the HTTP router:

```go
package main

import (
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
    "github.com/toutaio/touta-framework/pkg/touta"
)

func main() {
    // Create Toutā application
    app := touta.New()
    
    // Use Cosan as the router
    router := cosan.New()
    app.SetRouter(router)
    
    // Register routes through Cosan's clean API
    router.GET("/users", ListUsers)
    router.POST("/users", CreateUser)
    
    // Start application
    app.Run(":8080")
}
```

### With Message Bus Integration

Integrate with Toutā's message bus for event-driven architecture:

```go
package main

import (
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
    "github.com/toutaio/touta-messagebus/pkg/bus"
)

func main() {
    router := cosan.New()
    messageBus := bus.New()
    
    // Use AfterResponse hook to publish events
    router.AfterResponse(func(req *http.Request, statusCode int) {
        if statusCode == 201 {
            messageBus.Publish("resource.created", map[string]interface{}{
                "path":   req.URL.Path,
                "method": req.Method,
            })
        }
    })
    
    router.POST("/users", func(ctx cosan.Context) error {
        // Create user
        return ctx.JSON(201, user)
        // AfterResponse hook will publish event
    })
    
    router.Listen(":8080")
}
```

## Optional Ecosystem Component Adapters

Cosan can integrate with optional Toutā ecosystem components through adapters. **All adapters are optional** - use only what you need.

### 1. DataMapper Integration (Parameter Binding)

**Without DataMapper** (manual binding):
```go
router.POST("/users", func(ctx cosan.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return err
    }
    // Use user
    return ctx.JSON(201, user)
})
```

**With DataMapper** (advanced binding):
```go
import "github.com/toutaio/toutago-datamapper/pkg/datamapper"

// Setup (once)
mapper := datamapper.New()
// router.SetBinder(mapper) // Future: if we add binder support

// Usage with enhanced validation
router.POST("/users", func(ctx cosan.Context) error {
    var user User
    if err := mapper.BindAndValidate(ctx.Request(), &user); err != nil {
        return err
    }
    // User is validated and bound
    return ctx.JSON(201, user)
})
```

### 2. Fith Renderer Integration (Template Rendering)

**Without Fith** (manual rendering):
```go
router.GET("/page", func(ctx cosan.Context) error {
    html := "<h1>Hello</h1>"
    return ctx.HTML(200, html)
})
```

**With Fith** (template rendering):
```go
import "github.com/toutaio/toutago-fith-renderer/pkg/fith"

// Setup
renderer := fith.New()
renderer.LoadTemplates("./templates")

// Usage
router.GET("/page", func(ctx cosan.Context) error {
    data := map[string]interface{}{
        "title": "Welcome",
        "user":  currentUser,
    }
    html, err := renderer.Render("page.html", data)
    if err != nil {
        return err
    }
    return ctx.HTML(200, html)
})
```

### 3. Nasc DI Container Integration

**Without Nasc** (manual dependencies):
```go
// Pass dependencies manually
userService := services.NewUserService(db)

router.GET("/users", func(ctx cosan.Context) error {
    users := userService.List()
    return ctx.JSON(200, users)
})
```

**With Nasc** (dependency injection):
```go
import "github.com/toutaio/toutago-nasc-dependency-injector/pkg/nasc"

// Setup
container := nasc.New()
container.Singleton(func() *UserService {
    return services.NewUserService(db)
})

// Middleware to inject dependencies
router.Use(cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // Resolve dependencies
        userService := container.Make(&UserService{}).(*UserService)
        ctx.Set("userService", userService)
        return next(ctx)
    }
}))

// Usage
router.GET("/users", func(ctx cosan.Context) error {
    userService := ctx.Get("userService").(*UserService)
    users := userService.List()
    return ctx.JSON(200, users)
})
```

## Integration Patterns

### Pattern 1: Hooks for Cross-Cutting Concerns

Use Cosan's hooks to integrate with ecosystem services:

```go
// Logging integration
router.BeforeRequest(func(req *http.Request) error {
    logger.Info("Request started", req.Method, req.URL.Path)
    return nil
})

// Metrics integration
router.AfterResponse(func(req *http.Request, statusCode int) {
    metrics.RecordRequest(req.Method, req.URL.Path, statusCode)
})

// Event publishing
router.AfterResponse(func(req *http.Request, statusCode int) {
    if statusCode >= 200 && statusCode < 300 {
        events.Publish("http.success", req)
    }
})
```

### Pattern 2: Middleware for Service Integration

```go
// Authentication service integration
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
```

### Pattern 3: Context for Request-Scoped Data

```go
// Store request-scoped services in context
router.Use(cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // Create request-scoped transaction
        tx := db.BeginTransaction()
        defer tx.Rollback()
        
        ctx.Set("tx", tx)
        
        err := next(ctx)
        if err == nil {
            tx.Commit()
        }
        return err
    }
}))
```

## Migration from Pure Toutā to Cosan

### Before (Pure Toutā Framework)

```go
app := touta.New()
app.GET("/users", ListUsers)
app.POST("/users", CreateUser)
app.Run(":8080")
```

### After (Cosan with Toutā Integration)

```go
app := touta.New()
router := cosan.New()
app.SetRouter(router)

// Enjoy Cosan's clean API
router.GET("/users", ListUsers)
router.POST("/users", CreateUser)

app.Run(":8080")
```

## Benefits of Integration

| Feature | Standalone Cosan | With Toutā Ecosystem |
|---------|-----------------|---------------------|
| Routing | ✅ Full support | ✅ Full support |
| Middleware | ✅ Built-in | ✅ Built-in + ecosystem |
| Testing | ✅ Easy to mock | ✅ Easy to mock |
| Dependencies | ✅ Zero | ⚠️ Optional extras |
| Parameter Binding | ✅ Basic (JSON/XML) | ✅ Advanced (validation) |
| Templates | ✅ HTML strings | ✅ Full rendering engine |
| DI Container | ⚠️ Manual | ✅ Automatic |
| Message Bus | ⚠️ Manual integration | ✅ Built-in |
| Learning Curve | ✅ Simple | ⚠️ More concepts |

## Best Practices

### 1. Start Simple

Begin with standalone Cosan. Add integrations only when needed:

```go
// Week 1: Simple API
router := cosan.New()
router.GET("/health", HealthCheck)

// Week 2: Add middleware
router.Use(middleware.Logger())

// Week 3: Add ecosystem integration (if needed)
// Only add DataMapper if you need complex validation
// Only add Fith if you need templates
// Only add Nasc if you need DI
```

### 2. Keep Core Logic Independent

Don't couple your business logic to ecosystem components:

```go
// ✅ Good - Logic independent of framework
type UserService struct {
    db Database
}

func (s *UserService) CreateUser(u User) error {
    return s.db.Save(u)
}

// ❌ Bad - Logic coupled to framework
func CreateUser(ctx cosan.Context, mapper datamapper.Mapper) error {
    // Logic mixed with framework concerns
}
```

### 3. Use Interfaces

Define interfaces for ecosystem integrations:

```go
type Renderer interface {
    Render(template string, data interface{}) (string, error)
}

// Works with Fith or any other renderer
func RenderPage(ctx cosan.Context, renderer Renderer) error {
    html, err := renderer.Render("page.html", data)
    if err != nil {
        return err
    }
    return ctx.HTML(200, html)
}
```

## Testing with Integrations

### Testing Standalone Code

```go
func TestUserHandler(t *testing.T) {
    router := cosan.New()
    router.GET("/users", ListUsers)
    
    req := httptest.NewRequest("GET", "/users", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

### Testing with Mock Integrations

```go
type MockRenderer struct{}

func (m *MockRenderer) Render(t string, d interface{}) (string, error) {
    return "<html>Mock</html>", nil
}

func TestWithRenderer(t *testing.T) {
    renderer := &MockRenderer{}
    // Test handler with mock renderer
}
```

## Conclusion

Cosan Router is designed to:

1. **Work standalone** - Zero dependencies, production-ready
2. **Integrate optionally** - Add Toutā ecosystem components as needed
3. **Stay testable** - Interfaces and mocks work seamlessly
4. **Remain flexible** - Choose your level of integration

Start simple, add integrations when they provide clear value.

## Next Steps

- **Standalone**: See `examples/basic/` for getting started
- **With DataMapper**: See `examples/integration-datamapper/`
- **With Fith**: See `examples/integration-renderer/`
- **With Nasc**: See `examples/integration-di/`
- **Full Integration**: See `examples/full-integration/`

## Resources

- [Cosan Documentation](../README.md)
- [Toutā Framework](https://github.com/toutaio/touta-framework)
- [Migration Guide](./MIGRATION_FROM_TOUTA.md)
- [Examples](../examples/)
