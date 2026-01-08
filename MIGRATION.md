# Migration Guide

## v2.0.0 - Package Structure Refactoring

### Breaking Changes

**Import path change**: The main `cosan` package has moved from `pkg/cosan/` to the module root.

#### Before (v1.x)

```go
import "github.com/toutaio/toutago-cosan-router/pkg/cosan"
```

#### After (v2.x)

```go
import "github.com/toutaio/toutago-cosan-router"
```

### Middleware Import Change

Middleware has also moved:

#### Before (v1.x)

```go
import "github.com/toutaio/toutago-cosan-router/pkg/middleware"
```

#### After (v2.x)

```go
import "github.com/toutaio/toutago-cosan-router/middleware"
```

### Migration Steps

1. Update your `go.mod` to require v2:
   ```bash
   go get github.com/toutaio/toutago-cosan-router@v2
   ```

2. Update all import statements:
   ```bash
   # Find and replace in your codebase
   find . -name "*.go" -type f -exec sed -i 's|github.com/toutaio/toutago-cosan-router/pkg/cosan|github.com/toutaio/toutago-cosan-router|g' {} +
   find . -name "*.go" -type f -exec sed -i 's|github.com/toutaio/toutago-cosan-router/pkg/middleware|github.com/toutaio/toutago-cosan-router/middleware|g' {} +
   ```

3. Run `go mod tidy`

4. Test your application thoroughly

### Why This Change?

This change aligns with Go's best practices for library modules. The `pkg/` pattern is appropriate for applications with multiple binaries, but for standalone library modules, the primary package should be at the module root. This makes imports cleaner and more idiomatic.

### No Functional Changes

All functionality remains identical. Only the import paths have changed. This is purely a structural refactoring to follow Go community standards.

### Need Help?

If you encounter issues during migration, please [open an issue](https://github.com/toutaio/toutago-cosan-router/issues).
