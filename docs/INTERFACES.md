# Cosan Core Interfaces

This document provides detailed documentation for Cosan's core interfaces, demonstrating how they embody SOLID principles.

## Table of Contents

- [Overview](#overview)
- [Core Interfaces](#core-interfaces)
  - [HandlerFunc](#handlerfunc)
  - [Router](#router)
  - [Context](#context)
  - [Matcher](#matcher)
  - [Middleware](#middleware)
- [Optional Integration Interfaces](#optional-integration-interfaces)
  - [Binder](#binder)
  - [Renderer](#renderer)
  - [Container](#container)
- [SOLID Principles](#solid-principles)
- [Usage Examples](#usage-examples)

## Overview

Cosan's interfaces follow the **Interface Segregation Principle**, splitting functionality into focused, composable interfaces. This design enables:

- **Complete testability** - All components are mockable
- **Pluggable architecture** - Swap implementations without code changes
- **Clear contracts** - Each interface has one responsibility
- **Optional dependencies** - Use only what you need

## Core Interfaces

### HandlerFunc

```go
type HandlerFunc func(Context) error
```

The fundamental building block of Cosan routing. Handlers receive a Context and return an error.

**Why error return?**
- Centralized error handling
- Clear failure semantics
- Middleware can intercept errors

**Example:**
```go
func GetUser(ctx cosan.Context) error {
    id := ctx.Param("id")
    user, err := db.GetUser(id)
    if err != nil {
        return err // Will be handled by error handler
    }
    return ctx.JSON(200, user)
}
```

### Router

The Router interface defines HTTP routing and server management.

**Key Methods:**
- `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD` - HTTP method routing
- `Use(middleware ...)` - Register global middleware
- `Group(prefix string)` - Create route groups
- `ServeHTTP(w, r)` - Implements http.Handler
- `Listen(addr string)` - Convenience server start

**Design Decisions:**
- **Single Responsibility**: Only handles routing, not business logic
- **http.Handler compliance**: Works with standard library
- **Immutable after compilation**: Thread-safe serving

**Example:**
```go
router := cosan.New()

// Method-based routing
router.GET("/users", ListUsers)
router.POST("/users", CreateUser)
router.GET("/users/:id", GetUser)

// Middleware
router.Use(LoggingMiddleware, AuthMiddleware)

// Route groups
api := router.Group("/api/v1")
api.GET("/posts", ListPosts)

// Start server
router.Listen(":8080")
```

### Context

Context represents the request/response cycle. It's **composed** of smaller interfaces following the Interface Segregation Principle.

**Component Interfaces:**

#### ParamReader
```go
type ParamReader interface {
    Param(key string) string
    Params() map[string]string
}
```

Access to URL path parameters (`/users/:id`).

#### QueryReader
```go
type QueryReader interface {
    Query(key string) string
    QueryAll(key string) []string
}
```

Access to URL query parameters (`?name=value&tag=go&tag=web`).

#### BodyReader
```go
type BodyReader interface {
    Bind(v interface{}) error
    BodyBytes() ([]byte, error)
}
```

Access to request body with automatic parsing.

#### ResponseWriter
```go
type ResponseWriter interface {
    JSON(code int, v interface{}) error
    String(code int, format string, args ...interface{})
    Status(code int)
    Header() http.Header
    Write([]byte) (int, error)
}
```

Methods for writing HTTP responses.

**Full Context Interface:**
```go
type Context interface {
    ParamReader
    QueryReader
    BodyReader
    ResponseWriter
    
    Request() *http.Request
    Response() http.ResponseWriter
    Set(key string, value interface{})
    Get(key string) interface{}
}
```

**Why segregated?**
- **Testability**: Mock only what you need
- **Clarity**: Handler dependencies are explicit
- **Flexibility**: Implementations can be partial

**Example:**
```go
// Handler only needs ParamReader and ResponseWriter
func GetUser(ctx cosan.Context) error {
    // ParamReader
    id := ctx.Param("id")
    
    // ResponseWriter
    return ctx.JSON(200, map[string]string{
        "id": id,
        "name": "User " + id,
    })
}

// Another handler needs BodyReader
func CreateUser(ctx cosan.Context) error {
    // BodyReader
    var user User
    if err := ctx.Bind(&user); err != nil {
        return err
    }
    
    // Business logic...
    
    // ResponseWriter
    return ctx.JSON(201, user)
}
```

### Matcher

Matcher defines the route matching strategy. This allows pluggable algorithms.

**Interface:**
```go
type Matcher interface {
    Match(method, path string) (*Route, map[string]string, bool)
    Register(method, pattern string, handler HandlerFunc) error
    Compile() error
}
```

**Design Decisions:**
- **Strategy Pattern**: Different matching algorithms (radix tree, hash map, etc.)
- **Compile step**: Enables optimization before serving
- **Immutable after compile**: Thread-safe matching

**Example Implementation Strategy:**
```go
// Radix tree matcher (Phase 2)
type RadixMatcher struct {
    tree *radixTree
}

// Hash map matcher (simple, for exact matches only)
type HashMatcher struct {
    routes map[string]HandlerFunc
}

// Both implement Matcher interface
```

### Middleware

Middleware transforms requests/responses in a composable chain.

**Interface:**
```go
type Middleware interface {
    Process(next HandlerFunc) HandlerFunc
}
```

**Adapter Function:**
```go
type MiddlewareFunc func(HandlerFunc) HandlerFunc

func (mw MiddlewareFunc) Process(next HandlerFunc) HandlerFunc {
    return mw(next)
}
```

**Why both?**
- **Interface**: For stateful middleware (holds configuration)
- **Function**: For simple, stateless middleware

**Example:**
```go
// Stateless middleware using function
var LoggingMiddleware = cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        start := time.Now()
        log.Printf("→ %s %s", ctx.Request().Method, ctx.Request().URL.Path)
        
        err := next(ctx)
        
        log.Printf("← %s %s (%v)", ctx.Request().Method, ctx.Request().URL.Path, time.Since(start))
        return err
    }
})

// Stateful middleware using type
type RateLimiter struct {
    limiter *rate.Limiter
}

func (rl *RateLimiter) Process(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        if !rl.limiter.Allow() {
            ctx.Status(429)
            return nil
        }
        return next(ctx)
    }
}
```

## Optional Integration Interfaces

These interfaces enable ecosystem integrations but are **not required**.

### Binder

```go
type Binder interface {
    Bind(src interface{}, dst interface{}) error
}
```

**Purpose:** Advanced parameter binding (e.g., toutago-datamapper)

**Without Binder:**
```go
func GetUser(ctx cosan.Context) error {
    id := ctx.Param("id")
    name := ctx.Query("name")
    user := User{ID: id, Name: name}
    // ...
}
```

**With Binder:**
```go
func GetUser(ctx cosan.Context, user *User) error {
    // user automatically populated from request
    // ...
}
```

### Renderer

```go
type Renderer interface {
    Render(template string, data interface{}) (string, error)
}
```

**Purpose:** Template rendering (e.g., toutago-fith-renderer)

**Example:**
```go
router := cosan.New(cosan.WithRenderer(fith.NewRenderer()))

router.GET("/users/:id", func(ctx cosan.Context) error {
    user := getUser(ctx.Param("id"))
    return ctx.Render("user-profile", user)
})
```

### Container

```go
type Container interface {
    Make(typ interface{}) interface{}
    Bind(typ interface{}, impl interface{})
}
```

**Purpose:** Dependency injection (e.g., toutago-nasc-dependency-injector)

**Example:**
```go
router := cosan.New(cosan.WithContainer(nasc.New()))

router.GET("/users", func(ctx cosan.Context, userService UserService) error {
    // userService automatically injected
    users := userService.List()
    return ctx.JSON(200, users)
})
```

## SOLID Principles

### Single Responsibility Principle (SRP)

Each interface has one clear purpose:
- `Router` - HTTP routing only
- `Matcher` - Route matching only
- `Context` - Request/response access only
- `Middleware` - Request transformation only

### Open/Closed Principle (OCP)

Interfaces are:
- **Open for extension**: Via functional options, composition
- **Closed for modification**: Interface contracts don't change

Example: Adding custom middleware doesn't modify Router interface.

### Liskov Substitution Principle (LSP)

All implementations of an interface must be fully interchangeable:
- Any `Matcher` implementation works with `Router`
- Any `Middleware` implementation works in the chain
- `MiddlewareFunc` and custom types both implement `Middleware`

### Interface Segregation Principle (ISP)

Clients depend only on methods they use:
- `Context` is segregated into `ParamReader`, `QueryReader`, etc.
- Handlers can accept just `ParamReader` if that's all they need
- Tests mock only required interfaces

### Dependency Inversion Principle (DIP)

High-level modules depend on abstractions:
- `Router` depends on `Matcher` interface, not concrete implementation
- Handlers depend on `Context` interface, not concrete context
- No dependencies on concrete types in public API

## Usage Examples

### Basic Routing

```go
router := cosan.New()

router.GET("/", func(ctx cosan.Context) error {
    return ctx.String(200, "Welcome to Cosan!")
})

router.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    user := db.GetUser(id)
    return ctx.JSON(200, user)
})

router.POST("/users", func(ctx cosan.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        return err
    }
    db.CreateUser(&user)
    return ctx.JSON(201, user)
})

router.Listen(":8080")
```

### Middleware Chain

```go
router := cosan.New()

// Global middleware
router.Use(
    LoggingMiddleware,
    RecoveryMiddleware,
    CORSMiddleware,
)

// Route-specific via groups
api := router.Group("/api")
api.Use(AuthMiddleware)

api.GET("/profile", GetProfile)  // Protected by auth
```

### Route Groups

```go
router := cosan.New()

// API v1
v1 := router.Group("/api/v1")
v1.GET("/users", ListUsersV1)
v1.POST("/users", CreateUserV1)

// API v2 with different middleware
v2 := router.Group("/api/v2")
v2.Use(NewVersionMiddleware)
v2.GET("/users", ListUsersV2)
v2.POST("/users", CreateUserV2)

// Nested groups
admin := v2.Group("/admin")
admin.Use(AdminAuthMiddleware)
admin.DELETE("/users/:id", DeleteUser)
```

### Optional Integrations

```go
// Standalone (no integrations)
router := cosan.New()

// With all integrations
router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),
    cosan.WithRenderer(fith.NewRenderer()),
    cosan.WithContainer(nasc.New()),
)

// Mixed - only what you need
router := cosan.New(
    cosan.WithRenderer(fith.NewRenderer()),
)
```

---

For complete API documentation, see: https://pkg.go.dev/github.com/toutaio/toutago-cosan-router
