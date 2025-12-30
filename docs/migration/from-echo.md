# Migrating from Echo to Cosan

This guide helps you migrate from [labstack/echo](https://github.com/labstack/echo) to Cosan.

## Overview

Echo and Cosan share several design philosophies:
- Context-based request handling
- Middleware chain support
- Built-in response helpers
- High performance

**Key Differences:**
- Cosan is interface-driven for maximum testability
- Cosan emphasizes SOLID principles
- Cosan has optional ecosystem integrations
- Cosan's middleware signature is slightly different

## Quick Comparison

### Router Creation

**Echo:**
```go
e := echo.New()
```

**Cosan:**
```go
r := cosan.New()
```

### Route Registration

**Echo:**
```go
e.GET("/users/:id", func(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(200, map[string]string{
        "id": id,
    })
})
```

**Cosan:**
```go
r.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, map[string]string{
        "id": id,
    })
})
```

Very similar! The main difference is the router variable name convention.

### Route Groups

**Echo:**
```go
api := e.Group("/api")
api.Use(middleware.Logger())
api.GET("/users", listUsers)
```

**Cosan:**
```go
api := r.Group("/api")
api.Use(middleware.Logger())
api.GET("/users", listUsers)
```

## Step-by-Step Migration

### Step 1: Update Imports

**Before:**
```go
import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
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
    e := echo.New()
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    e.Logger.Fatal(e.Start(":8080"))
}
```

**After:**
```go
func main() {
    r := cosan.New()
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    
    log.Fatal(r.Listen(":8080"))
}
```

### Step 3: Convert Handler Signatures

Good news! Echo and Cosan handler signatures are very similar:

**Echo:**
```go
func getUser(c echo.Context) error {
    id := c.Param("id")
    return c.JSON(200, user)
}
```

**Cosan:**
```go
func getUser(ctx cosan.Context) error {
    id := ctx.Param("id")
    return ctx.JSON(200, user)
}
```

Just change `echo.Context` to `cosan.Context` and `c` to `ctx` (optional, for consistency).

### Step 4: Update Middleware

**Echo:**
```go
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        token := c.Request().Header.Get("Authorization")
        if token == "" {
            return c.JSON(401, map[string]string{"error": "Unauthorized"})
        }
        return next(c)
    }
}
```

**Cosan:**
```go
func authMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        token := ctx.Request().Header.Get("Authorization")
        if token == "" {
            return ctx.JSON(401, map[string]string{"error": "Unauthorized"})
        }
        return next(ctx)
    }
}
```

Almost identical! Just type changes.

### Step 5: Update Request Binding

**Echo:**
```go
var user User
if err := c.Bind(&user); err != nil {
    return c.JSON(400, map[string]string{"error": err.Error()})
}
return c.JSON(200, user)
```

**Cosan:**
```go
var user User
if err := ctx.BodyParser(&user); err != nil {
    return ctx.JSON(400, map[string]string{"error": err.Error()})
}
return ctx.JSON(200, user)
```

Change `c.Bind()` to `ctx.BodyParser()`.

### Step 6: Update Validation

**Echo (with validator):**
```go
var user User
if err := c.Bind(&user); err != nil {
    return err
}
if err := c.Validate(&user); err != nil {
    return err
}
```

**Cosan (with datamapper integration):**
```go
// Automatic binding and validation
func createUser(ctx cosan.Context, user *User) error {
    // user is automatically bound and validated
    return ctx.JSON(200, user)
}
```

## Common Patterns

### JSON Response

**Echo:**
```go
return c.JSON(200, data)
```

**Cosan:**
```go
return ctx.JSON(200, data)
```

Identical!

### String Response

**Echo:**
```go
return c.String(200, "Hello, World!")
```

**Cosan:**
```go
return ctx.String(200, "Hello, World!")
```

Identical!

### HTML Response

**Echo:**
```go
return c.HTML(200, "<h1>Hello</h1>")
```

**Cosan:**
```go
return ctx.HTML(200, "<h1>Hello</h1>")
```

Identical!

### Query Parameters

**Echo:**
```go
page := c.QueryParam("page")
limit := c.QueryParamDefault("limit", "10")
```

**Cosan:**
```go
page := ctx.Query("page")
limit := ctx.QueryDefault("limit", "10")
```

Slightly shorter method names.

### Path Parameters

**Echo:**
```go
id := c.Param("id")
```

**Cosan:**
```go
id := ctx.Param("id")
```

Identical!

### Form Values

**Echo:**
```go
name := c.FormValue("name")
```

**Cosan:**
```go
name := ctx.FormValue("name")
```

Identical!

### File Upload

**Echo:**
```go
file, err := c.FormFile("file")
if err != nil {
    return err
}
```

**Cosan:**
```go
file, err := ctx.FormFile("file")
if err != nil {
    return err
}
```

Identical!

### Cookies

**Echo:**
```go
cookie := &http.Cookie{
    Name:  "session",
    Value: "token",
}
c.SetCookie(cookie)

cookie, err := c.Cookie("session")
```

**Cosan:**
```go
cookie := &http.Cookie{
    Name:  "session",
    Value: "token",
}
ctx.SetCookie(cookie)

value := ctx.Cookie("session")
```

Very similar, `ctx.Cookie()` returns value directly.

### Redirects

**Echo:**
```go
return c.Redirect(302, "/new-url")
```

**Cosan:**
```go
return ctx.Redirect(302, "/new-url")
```

Identical!

### Static Files

**Echo:**
```go
e.Static("/static", "assets")
```

**Cosan:**
```go
r.Static("/static", "assets")
```

Identical!

## Context Values

**Echo:**
```go
// Set value
c.Set("user", user)

// Get value
user := c.Get("user").(*User)
```

**Cosan:**
```go
// Set value
ctx.Set("user", user)

// Get value
if val, ok := ctx.Get("user"); ok {
    user := val.(*User)
}
```

Similar, but Cosan's Get returns (value, bool) for safer access.

## Testing

**Echo:**
```go
func TestGetUser(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest("GET", "/users/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPath("/users/:id")
    c.SetParamNames("id")
    c.SetParamValues("1")
    
    if assert.NoError(t, getUser(c)) {
        assert.Equal(t, 200, rec.Code)
    }
}
```

**Cosan:**
```go
func TestGetUser(t *testing.T) {
    r := cosan.New()
    req := httptest.NewRequest("GET", "/users/1", nil)
    rec := httptest.NewRecorder()
    r.ServeHTTP(rec, req)
    
    assert.Equal(t, 200, rec.Code)
}

// Or using test context
func TestGetUserHandler(t *testing.T) {
    ctx := cosan.NewTestContext()
    ctx.SetParam("id", "1")
    
    err := getUser(ctx)
    assert.NoError(t, err)
    assert.Equal(t, 200, ctx.StatusCode())
}
```

## Middleware

Most Echo middleware can be easily ported:

**Echo CORS:**
```go
e.Use(middleware.CORS())
```

**Cosan CORS:**
```go
r.Use(middleware.CORS())
```

**Echo Rate Limiter:**
```go
e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
```

**Cosan Rate Limiter:**
```go
// Implement custom or use third-party
```

## Performance Considerations

Both Echo and Cosan offer excellent performance:

**Echo:**
- Optimized routing with radix tree
- Minimal allocations
- Context pooling

**Cosan:**
- Similar radix tree routing
- Context pooling
- Efficient middleware chain
- Typically within 5% of Echo

The main overhead in Cosan comes from the interface abstraction, which is minimal but provides significant testability benefits.

## Benefits of Migration

### SOLID Principles
- Better architecture
- More maintainable code
- Easier to test and extend

### Interface-Driven Design
- Mock any component
- Swap implementations
- Better testability

### Optional Integrations
- Data mapper for automatic binding
- Template renderer
- Dependency injection
- All completely optional

### No Framework Lock-in
- Builds on standard net/http
- Works anywhere
- Easy to integrate with other tools

## Migration Checklist

- [ ] Update imports (`echo` → `cosan`)
- [ ] Change variable names (`e` → `r`, `c` → `ctx`)
- [ ] Update `c.Bind()` to `ctx.BodyParser()`
- [ ] Change `c.QueryParam()` to `ctx.Query()`
- [ ] Update `c.Get()` to handle (value, bool) return
- [ ] Convert Echo-specific middleware
- [ ] Update tests
- [ ] Test thoroughly

## Echo Features Not (Yet) in Cosan

Some Echo features may not have direct equivalents:

1. **Built-in Template Rendering**: Use fith-renderer integration instead
2. **Built-in Validation**: Use datamapper integration instead
3. **Auto TLS**: Implement manually or use separate package
4. **HTTP/2 Server Push**: Use standard Go HTTP/2 features

## Gradual Migration

Run Echo and Cosan side-by-side during migration:

```go
func main() {
    echoRouter := echo.New()
    cosanRouter := cosan.New()
    
    // Old routes on Echo
    echoRouter.GET("/old/*", oldHandler)
    
    // New routes on Cosan
    cosanRouter.GET("/api/v2/*", newHandler)
    
    // Proxy Cosan through Echo
    echoRouter.Any("/api/v2/*", echo.WrapHandler(cosanRouter))
    
    echoRouter.Start(":8080")
}
```

## Common Gotchas

### 1. Context Get Return Type

**Echo:**
```go
user := c.Get("user").(*User) // Direct type assertion
```

**Cosan:**
```go
if val, ok := ctx.Get("user"); ok {
    user := val.(*User) // Safe type assertion
}
```

### 2. Error Handling

Both use similar error handling, but ensure you:
- Return errors from handlers
- Use error middleware for centralized handling
- Don't swallow errors

### 3. Binder Interface

**Echo:**
```go
if err := c.Bind(&user); err != nil {
    return err
}
```

**Cosan:**
```go
if err := ctx.BodyParser(&user); err != nil {
    return err
}
```

## Getting Help

- [Cosan Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
- [Examples](../../examples/)
- [GitHub Issues](https://github.com/toutaio/toutago-cosan-router/issues)

## See Also

- [Migration from Chi](./from-chi.md)
- [Migration from Gin](./from-gin.md)
- [Performance Guide](../guides/performance.md)

## Conclusion

Migration from Echo to Cosan is straightforward due to their similar designs. The main changes are:
1. Import paths
2. Variable naming conventions
3. Minor API differences

Most code can be migrated with search-and-replace, making the transition smooth and low-risk.
