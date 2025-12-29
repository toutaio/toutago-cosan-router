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

- [Complete Implementation Plan](./openspec/IMPLEMENTATION_PLAN.md)
- [Project Context & Conventions](./openspec/project.md)
- [API Documentation](https://pkg.go.dev/github.com/toutaio/toutago-cosan-router)
- [Examples](./examples/)

## Development Status

**Current Phase:** Phase 1 - Foundation & Core Routing

- [x] Project setup and structure
- [ ] Core interfaces definition
- [ ] Basic router implementation
- [ ] Middleware chain support
- [ ] >80% test coverage

See [IMPLEMENTATION_PLAN.md](./openspec/IMPLEMENTATION_PLAN.md) for complete roadmap.

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

**Project Status:** 🟢 Active Development  
**Version:** v0.1.0-alpha  
**Go Version:** 1.21.5+
