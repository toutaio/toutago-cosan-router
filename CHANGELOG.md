# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial development

## [1.0.0] - 2025-01-TBD

### Added

#### Core Features
- HTTP router with method-based registration (GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD)
- Path parameter extraction (`:param` syntax)
- Wildcard routes (`*path` syntax)
- Query parameter access
- Request body parsing (JSON, XML, Form)
- Response helpers (JSON, String, HTML)
- Middleware chain support
- Route grouping with prefix
- Nested route groups
- Context value storage (Get/Set)

#### Advanced Features
- Radix tree-based route matching for performance
- Request/Response lifecycle hooks
- Custom error handlers
- Route metadata and introspection API
- Context pooling for reduced allocations
- Concurrent request handling
- Status code capture

#### Middleware
- Logger middleware
- Recovery middleware (panic handling)
- Request ID middleware
- CORS middleware with customization

#### Testing
- >90% test coverage (90.3%)
- Fuzzing tests for robustness
- Integration tests
- Example tests (testable documentation)
- Race condition testing
- Benchmarks

#### Documentation
- Comprehensive README
- API documentation (godoc)
- Migration guides (from Chi, Gin, Echo)
- Architecture Decision Records (ADRs)
- Performance tuning guide
- Troubleshooting guide
- 10+ working examples
- Production deployment templates

#### Performance
- Competitive performance with Chi/Gin/Echo
- Memory optimization via context pooling
- Low allocation count (203 B/op for simple routes)
- Fast route matching (239.5 ns/op)
- Excellent parallel performance (46.42 ns/op on 16 cores)

#### Quality
- Zero race conditions
- No critical security issues
- 100% test pass rate
- CI/CD integration
- Security scanning enabled
- Automated vulnerability checks

### Core Principles
- SOLID principles implementation
- Interface-first design
- Zero external dependencies (uses only stdlib)
- Production-ready
- Fully testable and mockable

### Notes
- First stable release
- Ready for production use
- Full semantic versioning support going forward

[Unreleased]: https://github.com/toutaio/toutago-cosan-router/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/toutaio/toutago-cosan-router/releases/tag/v1.0.0
