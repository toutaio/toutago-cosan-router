# Cosan - HTTP Router for Go

[![CI](https://github.com/toutaio/toutago-cosan-router/actions/workflows/ci.yml/badge.svg)](https://github.com/toutaio/toutago-cosan-router/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/toutaio/toutago-cosan-router.svg)](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
[![Go Report Card](https://goreportcard.com/badge/github.com/toutaio/toutago-cosan-router)](https://goreportcard.com/report/github.com/toutaio/toutago-cosan-router)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

> **Cosan** (Irish for "pathway") - A production-ready HTTP router for Go that embodies SOLID principles and interface-first design.

## Features

- 🎯 **SOLID Principles Compliance** - Demonstrates all five SOLID principles in practice
- ⚡ **High Performance** - Competitive with Chi, Gin, Echo (within 10%)
- 🔌 **Interface-Driven Design** - Every component is mockable and testable
- 🔧 **Zero Framework Dependencies** - Works with standard `net/http`, usable anywhere
- 📦 **Pluggable Architecture** - Swap matchers, middleware, context implementations
- ✅ **Complete Testability** - >90% test coverage, fully mockable
- 🌐 **Optional Ecosystem Integrations** - Works with toutago-datamapper, fith-renderer, nasc

## Quick Start

```go
package main

import (
    "log"
    "github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

func main() {
    // Create a new router
    router := cosan.New()

    // Register routes
    router.GET("/", func(ctx cosan.Context) error {
        return ctx.JSON(200, map[string]string{
            "message": "Hello from Cosan!",
        })
    })

    router.GET("/users/:id", func(ctx cosan.Context) error {
        id := ctx.Param("id")
        return ctx.JSON(200, map[string]string{
            "id": id,
            "name": "User " + id,
        })
    })

    // Start the server
    log.Println("Server starting on :8080")
    log.Fatal(router.Listen(":8080"))
}
```

## Installation

```bash
go get github.com/toutaio/toutago-cosan-router
```

## Core Concepts

### Interface-First Design

Cosan is built on well-defined interfaces following the Interface Segregation Principle:

- `Router` - HTTP routing and server management
- `Context` - Request/response abstraction (composed of smaller interfaces)
- `Matcher` - Route matching strategy (pluggable)
- `Middleware` - Request transformation chain

### SOLID Principles

- **Single Responsibility**: Each component has one clear purpose
- **Open/Closed**: Extensible via functional options and interfaces
- **Liskov Substitution**: All implementations are fully interchangeable
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

### Optional Ecosystem Integrations

Cosan can integrate with Toutā ecosystem components through adapter interfaces:

```go
router := cosan.New(
    cosan.WithBinder(datamapper.NewBinder()),      // Optional parameter binding
    cosan.WithRenderer(fith.NewRenderer()),        // Optional template rendering
    cosan.WithContainer(nasc.New()),               // Optional DI container
)

router.GET("/users/:id", func(ctx cosan.Context, user *User) error {
    // user automatically injected and bound from request
    return ctx.Render("user-profile", user)
})
```

**All integrations are optional.** Cosan works perfectly as a standalone router.

## Documentation

### Getting Started
- [Installation & Quick Start](#quick-start)
- [Examples](./examples/) - Working code examples
- [API Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)

### Guides
- [Performance Tuning](./docs/guides/performance.md)
- [Troubleshooting](./docs/guides/troubleshooting.md)
- [Interface Documentation](./docs/INTERFACES.md)

### Migration
- [From Chi](./docs/migration/from-chi.md)
- [From Gin](./docs/migration/from-gin.md)
- [From Echo](./docs/migration/from-echo.md)

### Architecture
- [ADR 001: Interface-First Design](./docs/adr/001-interface-first-design.md)
- [ADR 002: Radix Tree Routing](./docs/adr/002-radix-tree-routing.md)
- [ADR 003: Error Handling Strategy](./docs/adr/003-error-handling-strategy.md)
- [ADR 004: Optional Ecosystem Integrations](./docs/adr/004-optional-ecosystem-integrations.md)

### Project Planning
- [Complete Implementation Plan](./openspec/IMPLEMENTATION_PLAN.md)
- [Project Context & Conventions](./openspec/project.md)

## Development Status

**Current Phase:** Phase 4 - Documentation & Ecosystem Integration

### Completed Phases
- [x] Phase 1: Foundation & Core Routing
  - [x] Project setup and structure
  - [x] Core interfaces definition
  - [x] Basic router implementation
  - [x] Middleware chain support
- [x] Phase 2: Advanced Routing Features
  - [x] Route groups
  - [x] Path parameters
  - [x] Wildcard routes
  - [x] Static file serving
- [x] Phase 3: Performance & Optimization
  - [x] Context pooling
  - [x] Route caching
  - [x] Benchmarks and profiling
  - [x] Performance optimization

### Current Phase
- [x] Phase 4.1: Documentation
  - [x] Comprehensive README
  - [x] API documentation
  - [x] Migration guides (Chi, Gin, Echo)
  - [x] Architecture Decision Records
  - [x] Performance tuning guide
  - [x] Troubleshooting guide
- [x] Phase 4.2: Examples & Templates
- [x] Phase 4.3: Testing & Quality
- [x] Phase 4.4: Community & Release
- [x] Phase 4.5: Ecosystem Integration

See [IMPLEMENTATION_PLAN.md](./openspec/IMPLEMENTATION_PLAN.md) for complete roadmap and release history.

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

Part of the [Toutā Framework](https://github.com/toutaio/toutago) ecosystem.

**Toutā** (Proto-Celtic for "people" or "tribe") - A message-driven Go web framework emphasizing:
- Interface-first design for pluggability
- Message-passing architecture
- Dependency injection for testability

---

**Project Status:** 🟢 Production Ready  
**Current Version:** See [Releases](https://github.com/toutaio/toutago-cosan-router/releases) for latest version  
**Go Version:** 1.22+
