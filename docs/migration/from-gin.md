# Migrating from Gin to Cosan

This guide helps you migrate from [gin-gonic/gin](https://github.com/gin-gonic/gin) to Cosan.

## Overview

Gin and Cosan both provide:
- Fast HTTP routing
- Middleware support
- Context-based request handling
- Built-in response helpers

**Key Differences:**
- Cosan is interface-driven for better testability
- Cosan emphasizes SOLID principles
- Cosan has no framework lock-in (builds on net/http)
- Cosan supports optional ecosystem integrations

## Quick Comparison

### Router Creation

**Gin:**
```go
r := gin.Default() // With logger and recovery
// or
r := gin.New()     // Without defaults
```

**Cosan:**
```go
r := cosan.New()
r.Use(middleware.Logger())
r.Use(middleware.Recovery())
```

### Route Registration

**Gin:**
```go
r.GET("/users/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{
        "id": id,
        "name": "User " + id,
    })
})
```

**Cosan:**
```go
r.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{
        "id": id,
        "name": "User " + id,
    })
})
```

### Route Groups

**Gin:**
```go
api := r.Group("/api")
api.Use(authMiddleware())
{
    api.GET("/users", listUsers)
    api.POST("/users", createUser)
}
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
import "github.com/gin-gonic/gin"
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
    r := gin.Default()
    r.Run(":8080")
}
```

**After:**
```go
func main() {
    r := cosan.New()
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    r.Listen(":8080")
}
```

### Step 3: Convert Handler Signatures

**Before:**
```go
func getUser(c *gin.Context) {
    id := c.Param("id")
    
    user, err := db.GetUser(id)
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }
    
    c.JSON(200, user)
}
```

**After:**
```go
func getUser(ctx cosan.Context) error {
    id := ctx.Param("id")
    
    user, err := db.GetUser(id)
    if err != nil {
        return ctx.JSON(404, map[string]string{
            "error": "User not found",
        })
    }
    
    return ctx.JSON(200, user)
}
```

### Step 4: Update Middleware

**Before:**
```go
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

**After:**
```go
func authMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        token := ctx.Request().Header.Get("Authorization")
        if token == "" {
            return ctx.JSON(401, map[string]string{
                "error": "Unauthorized",
            })
        }
        return next(ctx)
    }
}
```

### Step 5: Update Request Binding

**Gin:**
```go
var user User
if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
c.JSON(200, user)
```

**Cosan (built-in):**
```go
var user User
if err := ctx.BodyParser(&user); err != nil {
    return ctx.JSON(400, map[string]string{
        "error": err.Error(),
    })
}
return ctx.JSON(200, user)
```

**Cosan (with datamapper integration):**
```go
// Automatic binding with validation
func createUser(ctx cosan.Context, user *User) error {
    // user is automatically bound and validated
    return ctx.JSON(200, user)
}
```

### Step 6: Update Query Parameters

**Gin:**
```go
page := c.Query("page")
limit := c.DefaultQuery("limit", "10")
```

**Cosan:**
```go
page := ctx.Query("page")
limit := ctx.QueryDefault("limit", "10")
```

## Common Patterns

### JSON Response

**Gin:**
```go
c.JSON(200, gin.H{
    "message": "Success",
    "data": data,
})
```

**Cosan:**
```go
return ctx.JSON(200, map[string]interface{}{
    "message": "Success",
    "data": data,
})
```

### String Response

**Gin:**
```go
c.String(200, "Hello, %s", name)
```

**Cosan:**
```go
return ctx.String(200, fmt.Sprintf("Hello, %s", name))
```

### HTML Response

**Gin:**
```go
c.HTML(200, "index.html", gin.H{
    "title": "Home",
})
```

**Cosan (with fith-renderer integration):**
```go
return ctx.Render("index.html", map[string]interface{}{
    "title": "Home",
})
```

### File Upload

**Gin:**
```go
file, err := c.FormFile("file")
if err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
c.SaveUploadedFile(file, dst)
```

**Cosan:**
```go
file, err := ctx.FormFile("file")
if err != nil {
    return ctx.JSON(400, map[string]string{"error": err.Error()})
}
// Save file manually or use helper
return ctx.String(200, "File uploaded")
```

### Cookies

**Gin:**
```go
// Set cookie
c.SetCookie("name", "value", 3600, "/", "localhost", false, true)

// Get cookie
value, err := c.Cookie("name")
```

**Cosan:**
```go
// Set cookie
ctx.SetCookie(&http.Cookie{
    Name:   "name",
    Value:  "value",
    MaxAge: 3600,
    Path:   "/",
})

// Get cookie
value := ctx.Cookie("name")
```

### Redirects

**Gin:**
```go
c.Redirect(302, "/new-url")
```

**Cosan:**
```go
return ctx.Redirect(302, "/new-url")
```

### Abort Chain

**Gin:**
```go
if unauthorized {
    c.JSON(401, gin.H{"error": "Unauthorized"})
    c.Abort()
    return
}
```

**Cosan:**
```go
if unauthorized {
    return ctx.JSON(401, map[string]string{"error": "Unauthorized"})
    // Returning stops execution automatically
}
```

## Context Values

**Gin:**
```go
// Set value
c.Set("user", user)

// Get value
user, exists := c.Get("user")
if exists {
    u := user.(*User)
}
```

**Cosan:**
```go
// Set value
ctx.Set("user", user)

// Get value
if user, ok := ctx.Get("user"); ok {
    u := user.(*User)
}
```

## Testing

**Gin:**
```go
func TestGetUser(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.GET("/users/:id", getUser)
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/users/1", nil)
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

**Cosan:**
```go
func TestGetUser(t *testing.T) {
    r := cosan.New()
    r.GET("/users/:id", getUser)
    
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/users/1", nil)
    r.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}

// Or using mock context
func TestGetUserHandler(t *testing.T) {
    ctx := cosan.NewTestContext()
    ctx.SetParam("id", "1")
    
    err := getUser(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 200, ctx.StatusCode())
}
```

## Performance Considerations

**Gin** is known for excellent performance due to:
- httprouter-based routing
- Minimal allocations
- Optimized parameter parsing

**Cosan** provides comparable performance:
- Radix tree routing
- Pool-based context reuse
- Efficient middleware chain
- Typically within 5-10% of Gin

Trade-off: Cosan's interface abstraction adds minimal overhead but provides significant testability benefits.

## Benefits of Migration

### Better Architecture
- SOLID principles enforced
- Interface-driven design
- Better separation of concerns

### Enhanced Testability
- Mock Context interface
- No global state
- Easier unit testing

### Flexibility
- Pluggable components
- No framework lock-in
- Works with standard net/http

### Optional Integrations
- Data mapper for automatic binding
- Template renderer
- Dependency injection
- All completely optional

## Migration Checklist

- [ ] Update imports
- [ ] Change `gin.Context` to `cosan.Context`
- [ ] Update handler signatures (add `error` return)
- [ ] Convert middleware signatures
- [ ] Replace `c.JSON()` with `return ctx.JSON()`
- [ ] Replace `c.Abort()` with early `return`
- [ ] Update `gin.H{}` to `map[string]interface{}{}`
- [ ] Convert binding methods
- [ ] Update tests
- [ ] Remove gin-specific features

## Gradual Migration

You can run Gin and Cosan side-by-side:

```go
func main() {
    ginRouter := gin.Default()
    cosanRouter := cosan.New()
    
    // Old routes on Gin
    ginRouter.GET("/old/*any", oldHandler)
    
    // New routes on Cosan
    cosanRouter.GET("/api/v2/*", newHandler)
    
    // Combine routers
    ginRouter.Any("/api/v2/*any", gin.WrapH(cosanRouter))
    
    ginRouter.Run(":8080")
}
```

## Common Gotchas

### 1. Error Returns

**Gin:**
```go
func handler(c *gin.Context) {
    // No return value
}
```

**Cosan:**
```go
func handler(ctx cosan.Context) error {
    // Must return error (or nil)
    return nil
}
```

### 2. Middleware Chaining

**Gin:**
```go
func middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before
        c.Next()
        // After
    }
}
```

**Cosan:**
```go
func middleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // Before
        err := next(ctx)
        // After
        return err
    }
}
```

### 3. Aborting Request

**Gin:**
```go
c.JSON(400, gin.H{"error": "Bad request"})
c.Abort()
return
```

**Cosan:**
```go
return ctx.JSON(400, map[string]string{"error": "Bad request"})
// Return automatically stops execution
```

## Getting Help

- [Cosan Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
- [Examples](../../examples/)
- [GitHub Issues](https://github.com/toutaio/toutago-cosan-router/issues)

## See Also

- [Migration from Chi](./from-chi.md)
- [Migration from Echo](./from-echo.md)
- [Performance Guide](../guides/performance.md)
