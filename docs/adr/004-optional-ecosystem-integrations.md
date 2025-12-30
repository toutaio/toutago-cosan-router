# ADR 004: Optional Ecosystem Integrations

**Status:** Accepted  
**Date:** 2025-12-29  
**Deciders:** Core Team  

## Context

Cosan is part of the Toutā ecosystem, which includes:
- **toutago-datamapper**: Parameter binding and validation
- **toutago-fith-renderer**: Template rendering engine  
- **toutago-nasc-dependency-injector**: Dependency injection container

We needed to decide how Cosan should integrate with these components.

## Options Considered

### Option 1: Hard Dependencies
Require all Toutā components, bundle them with Cosan.

**Pros:**
- Easier integration
- Consistent behavior
- Simpler documentation

**Cons:**
- Not usable standalone
- Large dependency tree
- Forces users into ecosystem
- Violates independence principle

### Option 2: No Integration
Keep Cosan completely separate, no integration support.

**Pros:**
- Truly independent
- Minimal dependencies
- Simple codebase

**Cons:**
- Users can't benefit from ecosystem
- Have to manually integrate
- Fragmented experience

### Option 3: Optional Adapters (CHOSEN)
Provide adapter interfaces, make integrations optional.

**Pros:**
- Works standalone
- Works with ecosystem
- Users choose what to use
- Clean separation

**Cons:**
- More complex design
- Need adapter layer
- More documentation

## Decision

**Optional integration via adapter interfaces.**

Cosan works perfectly standalone. Users can optionally integrate ecosystem components through well-defined adapter interfaces.

## Design

### Core Principle

```
Cosan core = Zero dependencies on Toutā ecosystem
Adapters = Optional bridge to ecosystem components
```

### Adapter Interfaces

```go
// In cosan package
type Binder interface {
    Bind(ctx Context, v interface{}) error
}

type Renderer interface {
    Render(w io.Writer, template string, data interface{}) error
}

type Container interface {
    Resolve(ctx Context, handler interface{}) (HandlerFunc, error)
}
```

### Usage

**Standalone (no integrations):**
```go
router := cosan.New()

router.POST("/users", func(ctx cosan.Context) error {
    var user User
    if err := ctx.BodyParser(&user); err != nil {
        return err
    }
    // Manual validation
    if user.Email == "" {
        return NewHTTPError(400, "Email required")
    }
    return ctx.JSON(200, user)
})
```

**With datamapper integration:**
```go
router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),
)

router.POST("/users", func(ctx cosan.Context, user *User) error {
    // user is automatically bound and validated
    return ctx.JSON(200, user)
})
```

**With fith-renderer integration:**
```go
router := cosan.New(
    cosan.WithRenderer(fith.NewRenderer()),
)

router.GET("/users/:id", func(ctx cosan.Context) error {
    user := getUser(ctx.Param("id"))
    return ctx.Render("user-profile.html", user)
})
```

**With nasc DI integration:**
```go
container := nasc.New()
container.Register(&UserService{})

router := cosan.New(
    cosan.WithContainer(container),
)

router.GET("/users", func(ctx cosan.Context, svc *UserService) error {
    // svc is automatically injected
    users := svc.List()
    return ctx.JSON(200, users)
})
```

**With all integrations:**
```go
router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),
    cosan.WithRenderer(fith.NewRenderer()),
    cosan.WithContainer(nasc.New()),
)

router.POST("/users", func(ctx cosan.Context, 
    user *User,           // From binder
    svc *UserService,     // From container
) error {
    created, err := svc.Create(user)
    if err != nil {
        return err
    }
    return ctx.Render("user-created.html", created)
})
```

## Implementation

### Functional Options Pattern

```go
type RouterOptions struct {
    binder    Binder
    renderer  Renderer
    container Container
    matcher   Matcher
}

type Option func(*RouterOptions)

func WithBinder(b Binder) Option {
    return func(o *RouterOptions) {
        o.binder = b
    }
}

func WithRenderer(r Renderer) Option {
    return func(o *RouterOptions) {
        o.renderer = r
    }
}

func WithContainer(c Container) Option {
    return func(o *RouterOptions) {
        o.container = c
    }
}

func New(opts ...Option) Router {
    options := &RouterOptions{
        matcher: NewRadixMatcher(), // Default
    }
    for _, opt := range opts {
        opt(options)
    }
    return &DefaultRouter{
        options: options,
    }
}
```

### Adapter Packages

Create separate adapter packages:

```
pkg/
  cosan/           # Core (no Toutā deps)
  adapters/
    datamapper/    # Adapter for datamapper
    fith/          # Adapter for fith-renderer
    nasc/          # Adapter for nasc DI
```

Each adapter:
1. Depends on the specific Toutā component
2. Implements Cosan's interface
3. Provides `New()` function returning interface

## Benefits

### For Standalone Users

```go
// Just use Cosan, no ecosystem
go get github.com/toutaio/toutago-cosan-router

router := cosan.New()
// Works perfectly, zero Toutā dependencies
```

### For Ecosystem Users

```go
// Opt into ecosystem components
go get github.com/toutaio/toutago-cosan-router
go get github.com/toutaio/toutago-datamapper
go get github.com/toutaio/toutago-cosan-router/pkg/adapters/datamapper

router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),
)
```

### For Custom Implementations

```go
// Implement your own adapter
type MyBinder struct{}

func (b *MyBinder) Bind(ctx cosan.Context, v interface{}) error {
    // Your custom binding logic
}

router := cosan.New(
    cosan.WithBinder(&MyBinder{}),
)
```

## Trade-offs

### Advantages

- ✅ Cosan is truly independent
- ✅ No forced dependencies
- ✅ Works in any project
- ✅ Users choose what to use
- ✅ Easy to add custom implementations
- ✅ Clean separation of concerns
- ✅ Testable (can mock interfaces)

### Disadvantages

- ❌ More complex than hard integration
- ❌ Need to maintain adapter layer
- ❌ Users need to know about adapters
- ❌ More documentation required
- ❌ Potential version compatibility issues

### Mitigations

- Document clearly which features require which adapters
- Provide examples with and without integrations
- Test adapters separately
- Version adapters with Cosan
- Clear error messages when adapter is needed but missing

## Documentation Requirements

### 1. README

Show both standalone and integrated usage prominently.

### 2. Examples

Provide examples for:
- Standalone Cosan
- With datamapper only
- With fith-renderer only
- With nasc DI only
- With all integrations
- With custom adapters

### 3. Comparison Guide

| Feature | Standalone | With datamapper | With fith | With nasc |
|---------|------------|----------------|-----------|-----------|
| JSON responses | ✅ | ✅ | ✅ | ✅ |
| Auto binding | ❌ | ✅ | ❌ | ❌ |
| Validation | Manual | ✅ | ❌ | ❌ |
| Templates | Manual | ❌ | ✅ | ✅ |
| DI | Manual | ❌ | ❌ | ✅ |

## Testing Strategy

### Test Cosan Without Adapters

```go
func TestRouterStandalone(t *testing.T) {
    router := cosan.New()
    // Should work perfectly
}
```

### Test With Mock Adapters

```go
func TestRouterWithBinder(t *testing.T) {
    mockBinder := &MockBinder{}
    router := cosan.New(cosan.WithBinder(mockBinder))
    // Test binding behavior
}
```

### Integration Tests

```go
func TestDatamapperIntegration(t *testing.T) {
    router := cosan.New(
        cosan.WithBinder(datamapper.NewBinder()),
    )
    // Test real integration
}
```

## Consequences

### Positive

- ✅ Cosan is independent and reusable
- ✅ Ecosystem integrations available when needed
- ✅ Users have full control
- ✅ Easy to test
- ✅ Clean architecture
- ✅ Can integrate with any data mapper, renderer, or DI container

### Negative

- ❌ More complex than monolithic design
- ❌ Need adapter maintenance
- ❌ Requires good documentation

### Neutral

- 🔶 Need to version adapters carefully
- 🔶 Adapter interfaces become API contracts
- 🔶 Breaking changes more complex

## Future Considerations

### More Adapters

Could add adapters for:
- Other validation libraries (go-playground/validator)
- Other template engines (standard html/template)
- Other DI containers (uber/fx, google/wire)
- Message bus integration
- Caching layers
- Session management

### Adapter Registry

Could create adapter registry for discovery:

```go
adapters.Register("datamapper", datamapper.NewBinder())
router := cosan.New(
    cosan.WithAdapter("datamapper"),
)
```

## References

- [Functional Options Pattern](https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)

## Notes

This approach ensures Cosan remains independent while enabling powerful integrations. It's the best balance between flexibility and integration.
