// Package cosan provides a production-ready HTTP router for Go that embodies SOLID principles.
//
// Cosan (Irish: "pathway") is designed as an independent, framework-agnostic router that
// demonstrates interface-first architectural design. It works standalone or integrates
// seamlessly with the Toutā framework ecosystem.
//
// # Features
//
//   - SOLID Principles Compliance - All five principles demonstrated in practice
//   - High Performance - Competitive with Chi, Gin, Echo (within 10%)
//   - Interface-Driven Design - Every component mockable and testable
//   - Zero Framework Dependencies - Works with standard net/http
//   - Pluggable Architecture - Swap matchers, middleware, context implementations
//   - Complete Testability - >90% test coverage
//   - Optional Ecosystem Integrations - Works with toutago-datamapper, fith-renderer, nasc
//
// # Quick Start
//
//	router := cosan.New()
//
//	router.GET("/", func(ctx cosan.Context) error {
//	    return ctx.JSON(200, map[string]string{"message": "Hello"})
//	})
//
//	router.GET("/users/:id", func(ctx cosan.Context) error {
//	    id := ctx.Param("id")
//	    return ctx.JSON(200, map[string]string{"id": id})
//	})
//
//	log.Fatal(router.Listen(":8080"))
//
// # SOLID Principles
//
//   - Single Responsibility: Each component has one clear purpose
//   - Open/Closed: Extensible via functional options and interfaces
//   - Liskov Substitution: All implementations fully interchangeable
//   - Interface Segregation: Small, focused interfaces
//   - Dependency Inversion: Depend on abstractions, not concretions
//
// # Architecture
//
// Cosan uses pluggable matchers for route resolution:
//
//	router := cosan.New(
//	    cosan.WithMatcher(matcher.NewRadixMatcher()),
//	    cosan.WithLogger(logger),
//	)
//
// Context provides segregated interfaces for different concerns:
//
//   - ParamReader: Path parameter access
//   - QueryReader: Query string access
//   - BodyReader: Request body parsing
//   - ResponseWriter: Response rendering
//
// # Middleware
//
// Middleware uses chain of responsibility pattern:
//
//	router.Use(loggingMiddleware)
//	router.Use(authMiddleware)
//	router.GET("/protected", handler)
//
// # Thread Safety
//
// Routes are immutable after compilation. The router is thread-safe for concurrent
// request handling with no locks in the hot path.
//
// # Version
//
// This is version 1.0.0 - production ready with 90%+ test coverage.
// Requires Go 1.22 or higher.
package cosan
