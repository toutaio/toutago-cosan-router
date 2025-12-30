# Cosan Troubleshooting Guide

Common issues and solutions when using Cosan router.

## Installation Issues

### Issue: `go get` fails

**Error:**
```
go: github.com/toutaio/toutago-cosan-router@latest: reading github.com/toutaio/toutago-cosan-router/go.mod at revision latest: unknown revision latest
```

**Solutions:**
1. Check if repository exists and is accessible
2. Verify your Go version: `go version` (requires 1.21+)
3. Try with specific version: `go get github.com/toutaio/toutago-cosan-router@v0.1.0`
4. Clear module cache: `go clean -modcache`

### Issue: Import path not found

**Error:**
```
package github.com/toutaio/toutago-cosan-router/pkg/cosan is not in GOROOT
```

**Solutions:**
1. Run `go mod tidy`
2. Initialize go module if needed: `go mod init yourproject`
3. Verify import path matches package name

## Routing Issues

### Issue: 404 Not Found for valid routes

**Problem:**
```go
router.GET("/users/:id", getUser)
// Request to /users/123 returns 404
```

**Solutions:**

1. **Check route pattern syntax:**
```go
// ✅ Correct
router.GET("/users/:id", getUser)

// ❌ Wrong
router.GET("/users/{id}", getUser)  // Use :id not {id}
router.GET("/users/<id>", getUser)  // Use :id not <id>
```

2. **Verify HTTP method:**
```go
router.GET("/users", handler)
// POST request will return 404/405
```

3. **Check middleware blocking:**
```go
router.Use(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // If you return early without calling next, request stops
        if !authorized {
            return ctx.String(401, "Unauthorized")
        }
        return next(ctx) // ✅ Must call next
    }
})
```

### Issue: Route conflicts

**Problem:**
```go
router.GET("/users/admin", adminHandler)
router.GET("/users/:id", userHandler)
// Which one matches /users/admin?
```

**Solution:**

Routes are matched by priority:
1. Static routes first: `/users/admin`
2. Parameter routes second: `/users/:id`
3. Wildcard routes last: `/users/*`

Register specific routes before generic ones:
```go
router.GET("/users/admin", adminHandler)    // Matched first
router.GET("/users/:id", userHandler)       // Matched if not admin
```

### Issue: Wildcard capturing too much

**Problem:**
```go
router.GET("/files/*", fileHandler)
router.GET("/api/users", usersHandler)
// /api/users might match /files/*
```

**Solution:**

Wildcards match everything after them. Be specific:
```go
router.GET("/api/users", usersHandler)     // Register specific routes first
router.GET("/files/*filepath", fileHandler) // Wildcard last
```

## Handler Issues

### Issue: Handler not called

**Problem:**
```go
router.GET("/test", func(ctx cosan.Context) error {
    fmt.Println("Handler called")
    return ctx.String(200, "OK")
})
// Handler never executes
```

**Solutions:**

1. **Check middleware isn't blocking:**
```go
router.Use(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        return next(ctx) // ✅ Must call next
    }
})
```

2. **Verify router is being served:**
```go
func main() {
    router := cosan.New()
    router.GET("/test", handler)
    
    // ✅ Start server
    router.Listen(":8080")
}
```

3. **Check for panic:**
```go
router.Use(middleware.Recovery()) // Add recovery middleware
```

### Issue: Error not being handled

**Problem:**
```go
func handler(ctx cosan.Context) error {
    return errors.New("something failed")
    // Error not appearing in response
}
```

**Solution:**

Add error handling middleware:
```go
router.Use(middleware.Recovery())
router.Use(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        err := next(ctx)
        if err != nil {
            log.Error(err)
            return ctx.JSON(500, map[string]string{
                "error": err.Error(),
            })
        }
        return nil
    }
})
```

## Context Issues

### Issue: Cannot read request body

**Problem:**
```go
func handler(ctx cosan.Context) error {
    var data MyStruct
    err := ctx.BodyParser(&data)
    // err != nil or data is empty
}
```

**Solutions:**

1. **Check Content-Type header:**
```bash
# ✅ Correct
curl -X POST -H "Content-Type: application/json" -d '{"name":"test"}' http://localhost:8080/api

# ❌ Wrong - missing Content-Type
curl -X POST -d '{"name":"test"}' http://localhost:8080/api
```

2. **Verify JSON structure:**
```go
type MyStruct struct {
    Name string `json:"name"` // ✅ Exported field with json tag
}

// ❌ Wrong - unexported field
type MyStruct struct {
    name string `json:"name"` // Won't be unmarshaled
}
```

3. **Check body wasn't already read:**
```go
// Body can only be read once
ctx.BodyParser(&data1) // ✅ Works
ctx.BodyParser(&data2) // ❌ Body already consumed
```

### Issue: Parameters not found

**Problem:**
```go
router.GET("/users/:id", func(ctx cosan.Context) error {
    id := ctx.Param("id")
    // id is empty
})
```

**Solutions:**

1. **Check parameter name matches:**
```go
// Route pattern
router.GET("/users/:userId", handler)

// Handler - must match exactly
func handler(ctx cosan.Context) error {
    id := ctx.Param("userId") // ✅ Correct
    id := ctx.Param("id")     // ❌ Wrong - parameter is "userId"
}
```

2. **Verify route is matched:**
```go
// Add logging to check
router.Use(func(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        log.Printf("Path: %s, Params: %v", ctx.Path(), ctx.Params())
        return next(ctx)
    }
})
```

## Middleware Issues

### Issue: Middleware order matters

**Problem:**
```go
router.Use(authMiddleware)
router.Use(middleware.Logger())
// Requests blocked before being logged
```

**Solution:**

Order middleware from general to specific:
```go
router.Use(middleware.Logger())    // Log everything
router.Use(middleware.Recovery())  // Recover from panics
router.Use(corsMiddleware)         // Handle CORS
router.Use(authMiddleware)         // Then auth
```

### Issue: Middleware not applied to all routes

**Problem:**
```go
router.GET("/public", publicHandler)

router.Use(authMiddleware) // Applied after route registration!

router.GET("/private", privateHandler)
// authMiddleware not applied to /public
```

**Solution:**

Register middleware before routes:
```go
router.Use(authMiddleware)           // ✅ First

router.GET("/public", publicHandler)
router.GET("/private", privateHandler)
```

Or use route groups:
```go
// Public routes
router.GET("/public", publicHandler)

// Private routes with auth
api := router.Group("/api")
api.Use(authMiddleware)
api.GET("/users", listUsers)
```

### Issue: Middleware blocking requests

**Problem:**
```go
func myMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // Do something
        // Forgot to call next!
        return nil
    }
}
```

**Solution:**

Always call `next()`:
```go
func myMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        // Before
        log.Println("Before handler")
        
        err := next(ctx) // ✅ Call next
        
        // After
        log.Println("After handler")
        
        return err
    }
}
```

## Testing Issues

### Issue: Tests fail with nil pointer

**Problem:**
```go
func TestHandler(t *testing.T) {
    ctx := &cosan.DefaultContext{} // nil pointer panic
    handler(ctx)
}
```

**Solution:**

Use test context constructor:
```go
func TestHandler(t *testing.T) {
    ctx := cosan.NewTestContext()
    ctx.SetParam("id", "123")
    
    err := handler(ctx)
    assert.NoError(t, err)
}
```

### Issue: Cannot mock Context

**Problem:**
```go
// Can't create mock because Context is concrete type
```

**Solution:**

Context is an interface - use any mock library:
```go
type MockContext struct {
    mock.Mock
}

func (m *MockContext) JSON(code int, data interface{}) error {
    args := m.Called(code, data)
    return args.Error(0)
}

func TestHandler(t *testing.T) {
    mockCtx := new(MockContext)
    mockCtx.On("JSON", 200, mock.Anything).Return(nil)
    
    err := handler(mockCtx)
    
    assert.NoError(t, err)
    mockCtx.AssertExpectations(t)
}
```

## Performance Issues

### Issue: Slow response times

**Diagnosis:**

1. **Profile the application:**
```go
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

Then:
```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

2. **Check middleware overhead:**
```go
// Temporarily remove middleware and test
router := cosan.New()
// router.Use(middleware.Logger()) // Comment out
// router.Use(middleware.Metrics()) // Comment out
```

3. **Database queries:**
- Add query logging
- Check for N+1 problems
- Verify indexes

**Solutions:**
- See [Performance Guide](performance.md)
- Reduce middleware
- Optimize database queries
- Add caching layer

### Issue: High memory usage

**Diagnosis:**
```bash
go tool pprof http://localhost:6060/debug/pprof/heap
```

**Common Causes:**
1. Memory leaks in middleware
2. Not releasing resources
3. Large response bodies
4. Goroutine leaks

**Solutions:**
```go
// Use defer to ensure cleanup
func handler(ctx cosan.Context) error {
    resource := acquireResource()
    defer releaseResource(resource) // ✅ Always cleanup
    
    // Use resource
    return ctx.JSON(200, data)
}

// Limit request body size
router.Use(middleware.BodyLimit(10 * 1024 * 1024)) // 10 MB
```

## Deployment Issues

### Issue: Port already in use

**Error:**
```
listen tcp :8080: bind: address already in use
```

**Solutions:**
1. Find and kill process: `lsof -i :8080` then `kill <PID>`
2. Use different port: `router.Listen(":8081")`
3. Check for another instance running

### Issue: CORS errors in browser

**Error:**
```
Access to fetch at 'http://localhost:8080' has been blocked by CORS policy
```

**Solution:**

Add CORS middleware:
```go
router.Use(middleware.CORS(middleware.CORSConfig{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
}))
```

### Issue: 502/504 errors behind proxy

**Problem:**
Nginx/ALB returns 502/504 Gateway Timeout

**Solutions:**

1. **Increase timeouts:**
```go
router := cosan.New(
    cosan.WithReadTimeout(30 * time.Second),
    cosan.WithWriteTimeout(30 * time.Second),
)
```

2. **Check health endpoint:**
```go
router.GET("/health", func(ctx cosan.Context) error {
    return ctx.String(200, "OK")
})
```

3. **Verify proxy configuration:**
```nginx
# Nginx
proxy_read_timeout 90s;
proxy_send_timeout 90s;
```

## Integration Issues

### Issue: Datamapper not binding

**Problem:**
```go
router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),
)

router.POST("/users", func(ctx cosan.Context, user *User) error {
    // user is nil or empty
})
```

**Solutions:**

1. **Check handler signature:**
```go
// ✅ Correct - pointer parameter
func handler(ctx cosan.Context, user *User) error

// ❌ Wrong - value parameter
func handler(ctx cosan.Context, user User) error
```

2. **Verify JSON tags:**
```go
type User struct {
    Name  string `json:"name"`  // ✅ json tag
    Email string `json:"email"` // ✅ json tag
}
```

3. **Check request:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}' \
  http://localhost:8080/users
```

## Getting Help

### Before asking for help:

1. ✅ Check this troubleshooting guide
2. ✅ Read the relevant documentation
3. ✅ Search existing issues
4. ✅ Create minimal reproduction
5. ✅ Check logs and error messages

### When reporting issues:

Include:
- Go version: `go version`
- Cosan version
- Minimal code to reproduce
- Expected vs actual behavior
- Full error messages
- Stack trace if panic

### Resources:

- [GitHub Issues](https://github.com/toutaio/toutago-cosan-router/issues)
- [Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
- [Examples](../../examples/)
- [Performance Guide](./performance.md)

## Common Patterns

### Health Check

```go
router.GET("/health", func(ctx cosan.Context) error {
    return ctx.JSON(200, map[string]string{
        "status": "healthy",
        "version": "1.0.0",
    })
})
```

### Request ID Tracking

```go
func requestIDMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        requestID := ctx.Request().Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateID()
        }
        ctx.Set("requestID", requestID)
        ctx.Response().Header().Set("X-Request-ID", requestID)
        return next(ctx)
    }
}
```

### Error Logging

```go
func errorLogger(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        err := next(ctx)
        if err != nil {
            log.Printf("[ERROR] %s %s: %v",
                ctx.Request().Method,
                ctx.Request().URL.Path,
                err,
            )
        }
        return err
    }
}
```

## Debug Mode

Enable debug logging:

```go
router := cosan.New(
    cosan.WithDebug(true),
)

// Logs:
// - Route registration
// - Request matching
// - Middleware execution
// - Error details
```

## Conclusion

Most issues can be resolved by:
1. Checking middleware order
2. Verifying route patterns
3. Ensuring next() is called
4. Using proper error handling
5. Reading error messages carefully

If you're still stuck, create an issue with a minimal reproduction!
