# ADR 001: Interface-First Design

**Status:** Accepted  
**Date:** 2025-12-29  
**Deciders:** Core Team  

## Context

We needed to decide on the fundamental architectural approach for Cosan. The options were:
1. Concrete implementation with limited extensibility
2. Interface-first design with pluggable components
3. Hybrid approach with some interfaces and some concrete types

## Decision

We chose **interface-first design** where every major component is defined by an interface.

## Rationale

### Benefits

1. **Testability**
   - Every component can be mocked
   - Unit tests don't require httptest setup
   - Middleware can be tested in isolation
   - Handlers can be tested without a full router

2. **Flexibility**
   - Users can swap out implementations
   - Easy to create custom matchers, contexts, etc.
   - Enables decorator pattern for extensions
   - Supports multiple implementations side-by-side

3. **SOLID Compliance**
   - Dependency Inversion: Depend on abstractions
   - Interface Segregation: Small, focused interfaces
   - Liskov Substitution: All implementations interchangeable
   - Open/Closed: Extend without modifying

4. **Maintainability**
   - Clear contracts between components
   - Easier to refactor implementations
   - Better documentation through interfaces
   - Reduces coupling

### Trade-offs

1. **Slight Performance Overhead**
   - Interface calls have minimal indirection cost
   - Typically <2% overhead vs concrete types
   - Mitigated by inlining and compiler optimizations
   - **Decision:** Acceptable for the testability gains

2. **More Complex Type System**
   - More types to understand initially
   - Can be confusing for beginners
   - Requires good documentation
   - **Mitigation:** Comprehensive docs and examples

3. **Potential Over-Engineering**
   - Risk of creating unnecessary abstractions
   - Can make simple things complex
   - **Mitigation:** Only interface-ize components that benefit from it

## Implementation

### Core Interfaces

```go
type Router interface {
    GET(pattern string, handler HandlerFunc)
    POST(pattern string, handler HandlerFunc)
    // ... other HTTP methods
    Use(middleware ...Middleware)
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Context interface {
    Request() *http.Request
    Response() http.ResponseWriter
    Param(key string) string
    Query(key string) string
    JSON(code int, data interface{}) error
    // ... other methods
}

type Matcher interface {
    Match(path string, method string) (Match, bool)
    Add(pattern string, method string, handler HandlerFunc)
}
```

### Concrete Implementations

- `DefaultRouter` implements `Router`
- `DefaultContext` implements `Context`
- `RadixMatcher` implements `Matcher`

Users interact with interfaces but use concrete implementations under the hood.

## Consequences

### Positive

- ✅ Excellent testability - can mock everything
- ✅ Highly extensible - swap any component
- ✅ SOLID compliant architecture
- ✅ Clear separation of concerns
- ✅ Easy to add new features
- ✅ Ecosystem integrations possible (datamapper, renderer, DI)

### Negative

- ❌ Slightly higher learning curve
- ❌ Minimal performance overhead
- ❌ More types to document

### Neutral

- 🔶 Requires more upfront design
- 🔶 Need to maintain interface stability
- 🔶 Documentation becomes critical

## Examples

### Before (concrete types):
```go
func TestHandler(t *testing.T) {
    r := NewRouter()
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    // Assert on w
}
```

### After (interface-driven):
```go
func TestHandler(t *testing.T) {
    mockCtx := &MockContext{}
    mockCtx.On("JSON", 200, mock.Anything).Return(nil)
    
    err := handler(mockCtx)
    assert.NoError(t, err)
    mockCtx.AssertExpectations(t)
}
```

## References

- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [Interface Segregation Principle](https://en.wikipedia.org/wiki/Interface_segregation_principle)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

## Notes

This decision is fundamental to Cosan's architecture and should not be reversed without significant justification. It enables all other architectural decisions and is a core differentiator from other routers.
