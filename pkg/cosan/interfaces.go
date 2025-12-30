// Package cosan provides a production-ready HTTP router for Go that embodies
// SOLID principles and demonstrates interface-first architectural design.
//
// Cosan (Irish for "pathway") is an independent component that can be used in any
// Go project, with optional integrations for the Toutā ecosystem.
//
// # Core Principles
//
// - SOLID Principles: Demonstrates all five principles in practice
// - Interface-First: Every component is mockable and testable
// - Zero Dependencies: Works with standard net/http, usable anywhere
// - Performance First: Competitive with Chi, Gin, Echo (within 10%)
// - Pluggable: Swap matchers, middleware, context implementations
//
// # Quick Start
//
//	router := cosan.New()
//
//	router.GET("/", func(ctx cosan.Context) error {
//	    return ctx.JSON(200, map[string]string{
//	        "message": "Hello from Cosan!",
//	    })
//	})
//
//	router.GET("/users/:id", func(ctx cosan.Context) error {
//	    id := ctx.Param("id")
//	    return ctx.JSON(200, map[string]string{
//	        "id": id,
//	    })
//	})
//
//	log.Fatal(router.Listen(":8080"))
//
// # Optional Ecosystem Integrations
//
// Cosan can integrate with Toutā ecosystem components:
//
//	router := cosan.New(
//	    cosan.WithBinder(datamapper.NewBinder()),      // Optional parameter binding
//	    cosan.WithRenderer(fith.NewRenderer()),        // Optional template rendering
//	    cosan.WithContainer(nasc.New()),               // Optional DI container
//	)
//
// All integrations are optional. Cosan works perfectly standalone.
package cosan

import (
	"net/http"
)

// HandlerFunc defines the signature for HTTP request handlers.
// Handlers receive a Context and return an error.
//
// The error return allows for centralized error handling:
//   - nil: Request handled successfully
//   - error: Error occurred, will be handled by error handler
//
// Example:
//
//	func GetUser(ctx cosan.Context) error {
//	    id := ctx.Param("id")
//	    user, err := userService.Get(id)
//	    if err != nil {
//	        return err // Error will be handled centrally
//	    }
//	    return ctx.JSON(200, user)
//	}
type HandlerFunc func(Context) error

// Router defines the interface for HTTP routing and server management.
// It follows the Single Responsibility Principle by focusing solely on
// route registration and HTTP request handling.
//
// Implementations must be thread-safe for concurrent route registration
// before compilation and concurrent request handling after compilation.
//
// Example:
//
//	router := cosan.New()
//	router.GET("/users", ListUsers)
//	router.POST("/users", CreateUser)
//	router.GET("/users/:id", GetUser)
//	router.Use(LoggingMiddleware, RecoveryMiddleware)
//	router.Listen(":8080")
type Router interface {
	// GET registers a handler for GET requests matching the pattern.
	GET(pattern string, handler HandlerFunc)

	// POST registers a handler for POST requests matching the pattern.
	POST(pattern string, handler HandlerFunc)

	// PUT registers a handler for PUT requests matching the pattern.
	PUT(pattern string, handler HandlerFunc)

	// DELETE registers a handler for DELETE requests matching the pattern.
	DELETE(pattern string, handler HandlerFunc)

	// PATCH registers a handler for PATCH requests matching the pattern.
	PATCH(pattern string, handler HandlerFunc)

	// OPTIONS registers a handler for OPTIONS requests matching the pattern.
	OPTIONS(pattern string, handler HandlerFunc)

	// HEAD registers a handler for HEAD requests matching the pattern.
	HEAD(pattern string, handler HandlerFunc)

	// Use registers middleware to be applied to all routes.
	// Middleware is executed in the order registered (outer to inner).
	Use(middleware ...Middleware)

	// Group creates a route group with the given prefix.
	// Groups support scoped middleware and nested grouping.
	Group(prefix string) Router

	// ServeHTTP implements http.Handler interface.
	// This allows the router to be used with the standard library:
	//   http.ListenAndServe(":8080", router)
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Listen starts the HTTP server on the specified address.
	// This is a convenience method equivalent to:
	//   http.ListenAndServe(addr, router)
	Listen(addr string) error
}

// Route represents a registered HTTP route.
// This interface allows route introspection and metadata access.
type Route interface {
	// Pattern returns the route pattern (e.g., "/users/:id").
	Pattern() string

	// Method returns the HTTP method (e.g., "GET", "POST").
	Method() string

	// Handler returns the associated handler function.
	Handler() HandlerFunc
}

// ParamReader provides access to URL path parameters.
// This segregated interface follows the Interface Segregation Principle.
//
// Example:
//
//	// For route "/users/:id"
//	id := ctx.Param("id")
//	allParams := ctx.Params() // map[string]string{"id": "123"}
type ParamReader interface {
	// Param returns the value of the named path parameter.
	// Returns empty string if parameter doesn't exist.
	Param(key string) string

	// Params returns all path parameters as a map.
	Params() map[string]string
}

// QueryReader provides access to URL query parameters.
// This segregated interface follows the Interface Segregation Principle.
//
// Example:
//
//	// For URL "?name=John&tag=go&tag=web"
//	name := ctx.Query("name")           // "John"
//	tags := ctx.QueryAll("tag")         // []string{"go", "web"}
type QueryReader interface {
	// Query returns the first value of the named query parameter.
	// Returns empty string if parameter doesn't exist.
	Query(key string) string

	// QueryAll returns all values of the named query parameter.
	// Returns empty slice if parameter doesn't exist.
	QueryAll(key string) []string
}

// BodyReader provides access to request body content.
// This segregated interface follows the Interface Segregation Principle.
//
// Example:
//
//	var user User
//	if err := ctx.Bind(&user); err != nil {
//	    return err
//	}
type BodyReader interface {
	// Bind parses the request body into the provided struct.
	// Automatically detects Content-Type (JSON, XML, form).
	// Returns error if parsing fails.
	Bind(v interface{}) error

	// BodyBytes returns the raw request body as bytes.
	// Body can only be read once unless cached.
	BodyBytes() ([]byte, error)
}

// ResponseWriter provides methods for writing HTTP responses.
// This segregated interface follows the Interface Segregation Principle.
//
// Example:
//
//	ctx.JSON(200, map[string]string{"status": "ok"})
//	ctx.String(201, "Created resource %s", resourceID)
//	ctx.Status(204)
type ResponseWriter interface {
	// JSON writes a JSON response with the given status code.
	JSON(code int, v interface{}) error

	// String writes a formatted string response with the given status code.
	String(code int, format string, args ...interface{})

	// Status sets the HTTP status code.
	// Must be called before writing response body.
	Status(code int)

	// Header returns the response header map for modification.
	Header() http.Header

	// Write writes the response body bytes.
	// Implements io.Writer interface.
	Write([]byte) (int, error)
}

// Context represents the context of an HTTP request/response cycle.
// It composes smaller interfaces following the Interface Segregation Principle,
// allowing handlers to depend only on the methods they actually use.
//
// Context provides:
//   - Path parameter access (ParamReader)
//   - Query parameter access (QueryReader)
//   - Request body parsing (BodyReader)
//   - Response writing (ResponseWriter)
//   - Access to underlying http.Request and http.ResponseWriter
//
// Example:
//
//	func GetUser(ctx cosan.Context) error {
//	    id := ctx.Param("id")
//	    user, err := db.GetUser(id)
//	    if err != nil {
//	        return err
//	    }
//	    return ctx.JSON(200, user)
//	}
type Context interface {
	ParamReader
	QueryReader
	BodyReader
	ResponseWriter

	// Request returns the underlying *http.Request.
	// Useful for accessing headers, cookies, etc.
	Request() *http.Request

	// Response returns the underlying http.ResponseWriter.
	// Useful for low-level response manipulation.
	Response() http.ResponseWriter

	// Set stores a value in the context for the request lifetime.
	Set(key string, value interface{})

	// Get retrieves a value from the context.
	// Returns nil if key doesn't exist.
	Get(key string) interface{}
}

// Matcher defines the interface for route matching strategies.
// This allows pluggable matching algorithms (e.g., radix tree, trie, hash map).
//
// Implementations must be thread-safe after Compile() is called.
//
// Example:
//
//	matcher := NewRadixMatcher()
//	matcher.Register("GET", "/users/:id", handler)
//	matcher.Compile()
//	route, params, found := matcher.Match("GET", "/users/123")
type Matcher interface {
	// Match finds a route matching the method and path.
	// Returns the route, extracted parameters, and whether a match was found.
	Match(method, path string) (*Route, map[string]string, bool)

	// Register adds a route to the matcher.
	// Must be called before Compile().
	Register(method, pattern string, handler HandlerFunc) error

	// Compile optimizes the route tree for matching.
	// Must be called before Match() and after all routes are registered.
	// Route registration is not allowed after compilation.
	Compile() error
}

// Middleware defines the interface for request/response transformation.
// Middleware wraps a handler and returns a new handler.
//
// Middleware is executed in the order registered (outer to inner):
//   router.Use(LoggingMiddleware, AuthMiddleware)
//   // Request flow: Logging → Auth → Handler → Auth → Logging
//
// Example:
//
//	type loggingMiddleware struct{}
//
//	func (m *loggingMiddleware) Process(next HandlerFunc) HandlerFunc {
//	    return func(ctx Context) error {
//	        start := time.Now()
//	        err := next(ctx)
//	        log.Printf("Request took %v", time.Since(start))
//	        return err
//	    }
//	}
type Middleware interface {
	// Process wraps the next handler in the chain.
	// Middleware can execute code before and/or after calling next.
	Process(next HandlerFunc) HandlerFunc
}

// MiddlewareFunc is a function adapter for the Middleware interface.
// This allows functions to implement Middleware without creating a type.
//
// Example:
//
//	var LoggingMiddleware = cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
//	    return func(ctx cosan.Context) error {
//	        log.Printf("Request: %s %s", ctx.Request().Method, ctx.Request().URL.Path)
//	        return next(ctx)
//	    }
//	})
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// Process implements the Middleware interface.
func (mw MiddlewareFunc) Process(next HandlerFunc) HandlerFunc {
	return mw(next)
}

// ============================================================================
// Optional Integration Interfaces
// ============================================================================
//
// The following interfaces define optional integrations with ecosystem
// components. Cosan works perfectly without them - they enable advanced
// features when available.

// Binder defines the interface for advanced parameter binding.
// This is an optional integration for components like toutago-datamapper.
//
// When a Binder is configured, it enables automatic parameter binding:
//
//	router := cosan.New(cosan.WithBinder(datamapper.NewBinder()))
//
//	router.GET("/users/:id", func(ctx cosan.Context, user *User) error {
//	    // user is automatically populated from request
//	    return ctx.JSON(200, user)
//	})
//
// Without a Binder, manual parsing is required:
//
//	router.GET("/users/:id", func(ctx cosan.Context) error {
//	    var user User
//	    if err := ctx.Bind(&user); err != nil {
//	        return err
//	    }
//	    return ctx.JSON(200, user)
//	})
type Binder interface {
	// Bind maps source data to destination struct.
	// Source can be request params, query, body, etc.
	Bind(src interface{}, dst interface{}) error
}

// Renderer defines the interface for template rendering.
// This is an optional integration for components like toutago-fith-renderer.
//
// When a Renderer is configured, it enables template rendering:
//
//	router := cosan.New(cosan.WithRenderer(fith.NewRenderer()))
//
//	router.GET("/users/:id", func(ctx cosan.Context) error {
//	    user := getUser(ctx.Param("id"))
//	    return ctx.Render("user-profile", user)
//	})
//
// Without a Renderer, manual template handling is required.
type Renderer interface {
	// Render renders a template with the given data.
	// Returns the rendered string or an error.
	Render(template string, data interface{}) (string, error)
}

// Container defines the interface for dependency injection.
// This is an optional integration for components like toutago-nasc-dependency-injector.
//
// When a Container is configured, it enables automatic dependency injection:
//
//	router := cosan.New(cosan.WithContainer(nasc.New()))
//
//	router.GET("/users", func(ctx cosan.Context, userService UserService) error {
//	    // userService is automatically injected
//	    users := userService.List()
//	    return ctx.JSON(200, users)
//	})
//
// Without a Container, dependencies must be passed explicitly.
type Container interface {
	// Make resolves and returns an instance of the given type.
	Make(typ interface{}) interface{}

	// Bind registers a type mapping in the container.
	Bind(typ interface{}, impl interface{})
}
