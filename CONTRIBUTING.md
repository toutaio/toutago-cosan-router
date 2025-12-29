# Contributing to Cosan

Thank you for your interest in contributing to Cosan! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How to Contribute

### Reporting Issues

- Use the GitHub issue tracker
- Check if the issue already exists
- Provide detailed information:
  - Go version
  - Operating system
  - Steps to reproduce
  - Expected vs actual behavior

### Pull Requests

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Write or update tests
5. Ensure tests pass (`go test ./...`)
6. Ensure code is formatted (`go fmt ./...`)
7. Run linter (`golangci-lint run`)
8. Commit with conventional commit format
9. Push to your fork
10. Open a Pull Request

### Commit Convention

We use [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `perf`: Performance improvement
- `refactor`: Code restructuring
- `test`: Test additions or modifications
- `docs`: Documentation changes
- `chore`: Build, CI, or tooling changes

**Examples:**
```
feat(router): add support for route groups
fix(matcher): handle trailing slash correctly
perf(context): use sync.Pool for context reuse
docs(readme): add quick start guide
test(middleware): add chain execution tests
```

## Development Setup

### Prerequisites

- Go 1.21.5 or higher
- Git
- golangci-lint (for linting)

### Getting Started

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/toutago-cosan-router
cd toutago-cosan-router

# Install dependencies
go mod download

# Run tests
go test ./... -v -race

# Run linter
golangci-lint run
```

### Project Structure

```
cosan/
├── pkg/cosan/        # Public API
├── internal/         # Private implementation
├── cmd/              # CLI tools
├── examples/         # Usage examples
├── benchmarks/       # Performance tests
├── docs/             # Documentation
└── openspec/         # Planning and specs
```

## Design Principles

### SOLID Principles

All code must demonstrate SOLID principles:

1. **Single Responsibility**: Each component has one clear purpose
2. **Open/Closed**: Extensible via interfaces, not modification
3. **Liskov Substitution**: Implementations must be interchangeable
4. **Interface Segregation**: Small, focused interfaces
5. **Dependency Inversion**: Depend on abstractions

### Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Keep functions focused and small
- Write clear, descriptive names
- Comment exported functions and types
- Avoid global state

### Testing

- Write tests for all new code
- Maintain >90% coverage
- Use table-driven tests
- Test edge cases and errors
- Run with race detection (`-race`)

### Performance

- Profile before optimizing
- Minimize allocations in hot paths
- Use benchmarks (`go test -bench`)
- Target: within 10% of Chi/Gin performance

## Review Process

### What We Look For

- ✅ Follows SOLID principles
- ✅ Includes tests
- ✅ Documentation updated
- ✅ No breaking changes (unless justified)
- ✅ Passes CI checks
- ✅ Clear commit messages

### Review Timeline

- Initial response: within 48 hours
- Full review: within 1 week
- Merge: after approval and CI success

## Questions?

- Open an issue with the `question` label
- Join discussions on GitHub Discussions
- Check existing documentation

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Cosan! 🎉
