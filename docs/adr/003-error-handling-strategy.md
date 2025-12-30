# ADR 003: Error Handling Strategy

**Status:** Accepted  
**Date:** 2025-12-29  
**Deciders:** Core Team  

## Context

We needed to decide how handlers should handle errors. Options included:

1. **No error returns** (like Chi): Handlers write directly to response
2. **Error returns** (like Echo): Handlers return errors
3. **Panic/Recover** (like Gin): Panic on error, recover in middleware
4. **Result types**: Return (value, error) tuples

## Decision

Handlers return `error` and centralized error handling via middleware.

```go
type HandlerFunc func(ctx Context) error
```

## Rationale

### Benefits

1. **Explicit Error Handling**
   - Errors are visible in function signature
   - Can't accidentally ignore errors
   - Easy to see what can fail
   - Idiomatic Go

2. **Centralized Error Processing**
   - One place to handle all errors
   - Consistent error responses
   - Easy to add logging, monitoring
   - Simple to customize per environment

3. **Testability**
   - Easy to test error paths
   - No need to check response writer
   - Can assert on returned errors
   - Mock-friendly

4. **Composition**
   - Middleware can wrap and transform errors
   - Error stack traces easy to build
   - Can add context to errors
   - Supports error wrapping

### Trade-offs

**Compared to direct response writing:**
- ✅ More explicit
- ✅ Better testability
- ❌ Slightly more verbose
- ❌ Can't write partial response then error

**Compared to panic/recover:**
- ✅ More idiomatic Go
- ✅ Better control flow
- ✅ Easier to reason about
- ❌ Slightly more boilerplate

## Implementation

### Handler Signature

```go
type HandlerFunc func(ctx Context) error

func getUser(ctx Context) error {
    user, err := db.GetUser(ctx.Param("id"))
    if err != nil {
        return err // Will be handled by error middleware
    }
    return ctx.JSON(200, user)
}
```

### Error Middleware

```go
func ErrorHandler() Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx Context) error {
            err := next(ctx)
            if err != nil {
                // Log error
                log.Error(err)
                
                // Send appropriate response
                if httpErr, ok := err.(*HTTPError); ok {
                    return ctx.JSON(httpErr.Code, httpErr)
                }
                return ctx.JSON(500, map[string]string{
                    "error": "Internal Server Error",
                })
            }
            return nil
        }
    }
}
```

### Custom HTTP Errors

```go
type HTTPError struct {
    Code    int    `json:"-"`
    Message string `json:"message"`
}

func (e *HTTPError) Error() string {
    return e.Message
}

func NewHTTPError(code int, message string) *HTTPError {
    return &HTTPError{Code: code, Message: message}
}

// Usage
func getUser(ctx Context) error {
    if !authorized {
        return NewHTTPError(401, "Unauthorized")
    }
    // ...
}
```

## Error Handling Patterns

### 1. Return Error Directly

```go
func handler(ctx Context) error {
    user, err := db.GetUser(id)
    if err != nil {
        return err // Middleware will handle it
    }
    return ctx.JSON(200, user)
}
```

### 2. Return HTTP Error

```go
func handler(ctx Context) error {
    if !isValid {
        return NewHTTPError(400, "Invalid input")
    }
    return ctx.JSON(200, result)
}
```

### 3. Write Response Directly

```go
func handler(ctx Context) error {
    if !authorized {
        return ctx.JSON(401, map[string]string{
            "error": "Unauthorized",
        })
    }
    return ctx.JSON(200, data)
}
```

### 4. Wrap Errors

```go
func handler(ctx Context) error {
    user, err := db.GetUser(id)
    if err != nil {
        return fmt.Errorf("failed to get user %s: %w", id, err)
    }
    return ctx.JSON(200, user)
}
```

## Consequences

### Positive

- ✅ Idiomatic Go error handling
- ✅ Explicit and visible errors
- ✅ Easy to test
- ✅ Centralized error processing
- ✅ Composable middleware
- ✅ Support for error wrapping
- ✅ Clear control flow

### Negative

- ❌ More verbose than panic/recover
- ❌ Can't write partial response then error
- ❌ Need error middleware for complete solution
- ❌ Requires discipline to return errors correctly

### Neutral

- 🔶 Different from some frameworks (Gin)
- 🔶 Similar to others (Echo)
- 🔶 Matches standard library patterns

## Best Practices

### DO:

```go
// Return errors from business logic
func handler(ctx Context) error {
    user, err := getUser(id)
    if err != nil {
        return err
    }
    return ctx.JSON(200, user)
}

// Use HTTP errors for client errors
func handler(ctx Context) error {
    if !valid {
        return NewHTTPError(400, "Invalid request")
    }
    return ctx.JSON(200, result)
}

// Wrap errors with context
func handler(ctx Context) error {
    err := doSomething()
    if err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }
    return nil
}
```

### DON'T:

```go
// Don't ignore errors
func handler(ctx Context) error {
    result, _ := doSomething() // Bad!
    return ctx.JSON(200, result)
}

// Don't panic in handlers
func handler(ctx Context) error {
    if err != nil {
        panic(err) // Bad! Return error instead
    }
    return ctx.JSON(200, data)
}

// Don't swallow errors
func handler(ctx Context) error {
    if err := doSomething(); err != nil {
        log.Error(err) // Logged but not returned - Bad!
    }
    return ctx.JSON(200, "ok")
}
```

## Testing Impact

### Before (direct response writing):

```go
func TestHandler(t *testing.T) {
    w := httptest.NewRecorder()
    req := httptest.NewRequest("GET", "/", nil)
    handler(w, req)
    
    assert.Equal(t, 200, w.Code)
    // Have to parse response body to check error
}
```

### After (error returns):

```go
func TestHandler(t *testing.T) {
    mockCtx := &MockContext{}
    
    err := handler(mockCtx)
    
    // Can assert directly on error
    assert.Error(t, err)
    assert.Equal(t, 404, err.(*HTTPError).Code)
}
```

## Future Considerations

### Error Codes/Types

Could add error types for common scenarios:

```go
var (
    ErrNotFound     = NewHTTPError(404, "Not Found")
    ErrUnauthorized = NewHTTPError(401, "Unauthorized")
    ErrBadRequest   = NewHTTPError(400, "Bad Request")
)
```

### Error Translation

Could add middleware to translate different error types:

```go
func TranslateErrors() Middleware {
    return func(next HandlerFunc) HandlerFunc {
        return func(ctx Context) error {
            err := next(ctx)
            if err == nil {
                return nil
            }
            
            // Translate database errors
            if errors.Is(err, sql.ErrNoRows) {
                return ErrNotFound
            }
            
            // Translate validation errors
            if _, ok := err.(validator.ValidationErrors); ok {
                return NewHTTPError(400, "Validation failed")
            }
            
            return err
        }
    }
}
```

## References

- [Effective Go - Errors](https://golang.org/doc/effective_go#errors)
- [Go Blog - Error Handling](https://blog.golang.org/error-handling-and-go)
- [Echo Error Handling](https://echo.labstack.com/guide/error-handling/)

## Notes

This approach balances idiomatic Go practices with modern web framework ergonomics. It's explicit, testable, and composable while avoiding the pitfalls of exception-based error handling.
