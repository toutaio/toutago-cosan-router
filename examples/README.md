# Cosan Router Examples

This directory contains comprehensive examples demonstrating various features and integration capabilities of the Cosan router.

## Basic Examples

### 1. Basic Usage (`basic/`)
Demonstrates fundamental router features:
- Route registration (GET, POST, etc.)
- Path parameters
- Route groups
- JSON responses
- Request binding

**Run:**
```bash
cd examples/basic
go run main.go
```

**Try:**
```bash
curl http://localhost:8080/
curl http://localhost:8080/users/123
curl -X POST http://localhost:8080/users -d '{"name":"John","id":1}'
```

### 2. Middleware (`middleware/`)
Shows middleware usage patterns:
- Global middleware (logging, recovery)
- Route-specific middleware
- Group middleware (authentication, authorization)
- Custom middleware creation
- Middleware chaining

**Run:**
```bash
cd examples/middleware
go run main.go
```

**Try:**
```bash
# Public endpoint
curl http://localhost:8080/public

# Protected endpoint (requires auth)
curl -H "Authorization: Bearer valid-token" http://localhost:8080/api/profile

# Admin endpoint (requires admin role)
curl -H "Authorization: Bearer valid-token" http://localhost:8080/admin/dashboard
```

### 3. Path Parameters (`path-parameters/`)
Demonstrates advanced path parameter handling:
- Simple path parameters (`:id`)
- Multiple parameters (`:userId/posts/:postId`)
- Wildcard parameters (`*filepath`)
- Nested resources
- Query parameters combined with path parameters

**Run:**
```bash
cd examples/path-parameters
go run main.go
```

**Try:**
```bash
curl http://localhost:8080/users/123
curl http://localhost:8080/users/1/posts/42
curl http://localhost:8080/files/path/to/file.txt
curl http://localhost:8080/search/electronics?q=laptop&page=1
```

## Integration Examples

### 4. DataMapper Integration (`integration-datamapper/`)
Shows integration with `toutago-datamapper` for:
- Automatic request binding
- Data validation
- Query parameter binding
- Form data binding
- Multi-source binding (path + body + query)

**Run:**
```bash
cd examples/integration-datamapper
go run main.go
```

**Try:**
```bash
# Create user with validation
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","age":25}'

# Search with query parameters
curl "http://localhost:8080/search?username=john&min_age=18&max_age=65"
```

### 5. Renderer Integration (`integration-renderer/`)
Demonstrates integration with `toutago-fith-renderer` for:
- HTML page rendering
- Template rendering
- Dynamic content
- Form handling
- Partial rendering (HTMX/AJAX support)
- Content negotiation (HTML vs JSON)

**Run:**
```bash
cd examples/integration-renderer
go run main.go
```

**Visit in browser:**
- http://localhost:8080/
- http://localhost:8080/login
- http://localhost:8080/products

### 6. Dependency Injection (`integration-di/`)
Shows integration with `toutago-nasc-dependency-injector` for:
- Service layer architecture
- Automatic dependency resolution
- Constructor injection
- Controller pattern with DI
- Middleware with injected dependencies

**Run:**
```bash
cd examples/integration-di
go run main.go
```

**Try:**
```bash
curl http://localhost:8080/users/1
curl http://localhost:8080/products
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"username":"jane","email":"jane@example.com"}'
```

### 7. Full Integration (`full-integration/`)
Complete example integrating all toutago components:
- **Router**: Request routing and handling
- **DataMapper**: Parameter binding and validation
- **Fith Renderer**: HTML template rendering
- **NASC DI**: Dependency injection container

Features demonstrated:
- Layered architecture (Controllers, Services, Repositories)
- Content negotiation (HTML/JSON)
- RESTful API design
- Web pages with forms
- Middleware stack
- Error handling

**Run:**
```bash
cd examples/full-integration
go run main.go
```

**Try:**
```bash
# Web pages
curl http://localhost:8080/
curl http://localhost:8080/users/1

# JSON API
curl http://localhost:8080/api/v1/users/1
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","first_name":"Alice","last_name":"Smith"}'

# Query products
curl "http://localhost:8080/api/v1/products?category=electronics"
```

## Project Structure

Each example follows this structure:
```
example-name/
├── main.go          # Main application file
└── README.md        # Example-specific documentation (optional)
```

## Running Examples

All examples are self-contained and can be run independently:

```bash
# Navigate to example directory
cd examples/<example-name>

# Run the example
go run main.go

# The server will start on http://localhost:8080
```

## Integration Notes

### Using with DataMapper

To integrate `toutago-datamapper`:

```go
import "github.com/toutaio/toutago-datamapper"

// Configure router with datamapper
binder := datamapper.NewBinder()
validator := datamapper.NewValidator()

router.SetBinder(binder)
router.SetValidator(validator)

// Use in handlers
func CreateUser(ctx cosan.Context) error {
    var user User
    if err := ctx.Bind(&user); err != nil {
        // Validation errors are included
        return ctx.JSON(400, err)
    }
    // user is now validated and bound
}
```

### Using with Fith Renderer

To integrate `toutago-fith-renderer`:

```go
import "github.com/toutaio/toutago-fith-renderer"

// Configure renderer
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    Extension:   ".html",
    Cache:       true,
})

router.SetRenderer(renderer)

// Use in handlers
func ShowPage(ctx cosan.Context) error {
    return ctx.Render(200, "page", map[string]interface{}{
        "title": "Page Title",
        "data":  data,
    })
}
```

### Using with NASC DI

To integrate `toutago-nasc-dependency-injector`:

```go
import "github.com/toutaio/toutago-nasc-dependency-injector"

// Setup container
container := nasc.NewContainer()
container.Register(NewLogger)
container.Register(NewUserService)
container.Register(NewUserController)

router.SetContainer(container)

// Controllers are auto-injected
router.GET("/users/:id", container.Resolve(UserController).GetUser)
```

## Learning Path

Recommended order for learning:

1. **basic/** - Start here to understand router fundamentals
2. **middleware/** - Learn middleware patterns
3. **path-parameters/** - Master routing patterns
4. **integration-datamapper/** - Add data binding and validation
5. **integration-renderer/** - Add HTML rendering
6. **integration-di/** - Add dependency injection
7. **full-integration/** - See everything working together

## Best Practices Demonstrated

- **Separation of concerns**: Controllers, services, repositories
- **Error handling**: Consistent error responses
- **Middleware composition**: Reusable middleware functions
- **RESTful design**: Standard HTTP methods and status codes
- **Content negotiation**: Supporting multiple response formats
- **Dependency injection**: Loose coupling and testability
- **Validation**: Input validation at the entry point

## Contributing

When adding new examples:

1. Create a new directory under `examples/`
2. Include a `main.go` with clear comments
3. Make it self-contained (runnable without external setup)
4. Add example curl commands in comments
5. Update this README with the new example

## Support

For issues or questions:
- GitHub Issues: https://github.com/toutaio/toutago-cosan-router/issues
- Documentation: https://github.com/toutaio/toutago-cosan-router/docs

## License

All examples are part of the Cosan Router project and follow the same license.
