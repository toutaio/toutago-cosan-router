package cosan

import (
	"net/http"
	"sync"
)

// router is the default implementation of the Router interface.
// It provides method-based routing, middleware support, and exact path matching.
type router struct {
	routes     []*route
	middleware []Middleware
	matcher    Matcher
	compiled   bool
	mu         sync.RWMutex
}

// route represents a registered HTTP route.
type route struct {
	method  string
	pattern string
	handler HandlerFunc
}

// Pattern returns the route pattern.
func (r *route) Pattern() string {
	return r.pattern
}

// Method returns the HTTP method.
func (r *route) Method() string {
	return r.method
}

// Handler returns the handler function.
func (r *route) Handler() HandlerFunc {
	return r.handler
}

// New creates a new Router instance with default configuration.
//
// Example:
//
//	router := cosan.New()
//	router.GET("/", HomeHandler)
//	router.POST("/users", CreateUserHandler)
//	router.Listen(":8080")
func New(opts ...Option) Router {
	r := &router{
		routes:     make([]*route, 0),
		middleware: make([]Middleware, 0),
		matcher:    newRadixMatcher(), // Radix tree matcher with path parameters
		compiled:   false,
	}

	// Apply options
	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Option is a functional option for configuring the router.
type Option func(*router)

// WithMatcher sets a custom matcher implementation.
func WithMatcher(m Matcher) Option {
	return func(r *router) {
		r.matcher = m
	}
}

// GET registers a handler for GET requests.
func (r *router) GET(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodGet, pattern, handler)
}

// POST registers a handler for POST requests.
func (r *router) POST(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodPost, pattern, handler)
}

// PUT registers a handler for PUT requests.
func (r *router) PUT(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodPut, pattern, handler)
}

// DELETE registers a handler for DELETE requests.
func (r *router) DELETE(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodDelete, pattern, handler)
}

// PATCH registers a handler for PATCH requests.
func (r *router) PATCH(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodPatch, pattern, handler)
}

// OPTIONS registers a handler for OPTIONS requests.
func (r *router) OPTIONS(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodOptions, pattern, handler)
}

// HEAD registers a handler for HEAD requests.
func (r *router) HEAD(pattern string, handler HandlerFunc) {
	r.registerRoute(http.MethodHead, pattern, handler)
}

// Use registers middleware to be applied to all routes.
// Middleware is executed in the order registered (outer to inner).
//
// Example:
//
//	router.Use(LoggingMiddleware, RecoveryMiddleware)
func (r *router) Use(middleware ...Middleware) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.compiled {
		panic("cosan: cannot add middleware after router is compiled")
	}

	r.middleware = append(r.middleware, middleware...)
}

// Group creates a new route group with the given prefix.
// Groups support scoped middleware and nested grouping.
//
// Example:
//
//	api := router.Group("/api/v1")
//	api.GET("/users", ListUsers)
//	api.POST("/users", CreateUser)
func (r *router) Group(prefix string) Router {
	// For Phase 1, we'll return a simple group wrapper
	return &routerGroup{
		router: r,
		prefix: prefix,
	}
}

// ServeHTTP implements http.Handler interface.
// This allows the router to be used with the standard library.
//
// Example:
//
//	http.ListenAndServe(":8080", router)
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Ensure router is compiled
	r.ensureCompiled()

	// Match route
	routeInterface, params, found := r.matcher.Match(req.Method, req.URL.Path)
	if !found {
		// No route found - return 404
		http.NotFound(w, req)
		return
	}

	// Create context
	ctx := newContext(w, req, params)

	// Get handler from route interface
	handler := (*routeInterface).Handler()

	// Apply middleware chain
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i].Process(handler)
	}

	// Execute handler
	if err := handler(ctx); err != nil {
		// TODO: Phase 1.4 - Proper error handling
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Listen starts the HTTP server on the specified address.
// This is a convenience method equivalent to http.ListenAndServe(addr, router).
//
// Example:
//
//	router.Listen(":8080")
func (r *router) Listen(addr string) error {
	return http.ListenAndServe(addr, r)
}

// registerRoute registers a new route with the router.
func (r *router) registerRoute(method, pattern string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.compiled {
		panic("cosan: cannot register routes after router is compiled")
	}

	// Check for conflicts
	for _, existing := range r.routes {
		if existing.method == method && existing.pattern == pattern {
			panic("cosan: duplicate route registration: " + method + " " + pattern)
		}
	}

	// Create and store route
	rt := &route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}
	r.routes = append(r.routes, rt)

	// Register with matcher
	if err := r.matcher.Register(method, pattern, handler); err != nil {
		panic("cosan: failed to register route: " + err.Error())
	}
}

// ensureCompiled ensures the router is compiled before serving requests.
func (r *router) ensureCompiled() {
	r.mu.RLock()
	if r.compiled {
		r.mu.RUnlock()
		return
	}
	r.mu.RUnlock()

	// Acquire write lock to compile
	r.mu.Lock()
	defer r.mu.Unlock()

	// Double-check after acquiring write lock
	if r.compiled {
		return
	}

	// Compile the matcher
	if err := r.matcher.Compile(); err != nil {
		panic("cosan: failed to compile router: " + err.Error())
	}

	r.compiled = true
}

// routerGroup represents a route group with a common prefix.
type routerGroup struct {
	router *router
	prefix string
}

// GET registers a GET route in the group.
func (g *routerGroup) GET(pattern string, handler HandlerFunc) {
	g.router.GET(g.prefix+pattern, handler)
}

// POST registers a POST route in the group.
func (g *routerGroup) POST(pattern string, handler HandlerFunc) {
	g.router.POST(g.prefix+pattern, handler)
}

// PUT registers a PUT route in the group.
func (g *routerGroup) PUT(pattern string, handler HandlerFunc) {
	g.router.PUT(g.prefix+pattern, handler)
}

// DELETE registers a DELETE route in the group.
func (g *routerGroup) DELETE(pattern string, handler HandlerFunc) {
	g.router.DELETE(g.prefix+pattern, handler)
}

// PATCH registers a PATCH route in the group.
func (g *routerGroup) PATCH(pattern string, handler HandlerFunc) {
	g.router.PATCH(g.prefix+pattern, handler)
}

// OPTIONS registers an OPTIONS route in the group.
func (g *routerGroup) OPTIONS(pattern string, handler HandlerFunc) {
	g.router.OPTIONS(g.prefix+pattern, handler)
}

// HEAD registers a HEAD route in the group.
func (g *routerGroup) HEAD(pattern string, handler HandlerFunc) {
	g.router.HEAD(g.prefix+pattern, handler)
}

// Use adds middleware to the group (currently global, will be scoped in Phase 2).
func (g *routerGroup) Use(middleware ...Middleware) {
	g.router.Use(middleware...)
}

// Group creates a nested group.
func (g *routerGroup) Group(prefix string) Router {
	return g.router.Group(g.prefix + prefix)
}

// ServeHTTP implements http.Handler (delegates to parent router).
func (g *routerGroup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.router.ServeHTTP(w, r)
}

// Listen starts the server (delegates to parent router).
func (g *routerGroup) Listen(addr string) error {
	return g.router.Listen(addr)
}
