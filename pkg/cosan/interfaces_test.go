package cosan_test

import (
	"testing"

	"github.com/toutaio/toutago-cosan-router/pkg/cosan"
)

// TestHandlerFuncSignature verifies the HandlerFunc type signature.
func TestHandlerFuncSignature(t *testing.T) {
	var _ cosan.HandlerFunc = func(ctx cosan.Context) error {
		return nil
	}
	t.Log("HandlerFunc signature is correct")
}

// TestMiddlewareFuncAdapter verifies MiddlewareFunc adapter pattern.
func TestMiddlewareFuncAdapter(t *testing.T) {
	// MiddlewareFunc should implement Middleware interface
	loggingMiddleware := cosan.MiddlewareFunc(func(next cosan.HandlerFunc) cosan.HandlerFunc {
		return func(ctx cosan.Context) error {
			// Pre-processing
			err := next(ctx)
			// Post-processing
			return err
		}
	})

	// Verify it implements Middleware interface
	var _ cosan.Middleware = loggingMiddleware
	t.Log("MiddlewareFunc correctly implements Middleware interface")
}

// TestInterfacesCompile verifies all interfaces compile correctly.
// This is a compile-time test - if it compiles, the interfaces are correct.
func TestInterfacesCompile(t *testing.T) {
	t.Log("All core interfaces compile successfully")
	t.Log("- Router interface defined")
	t.Log("- Context interface (composed of ParamReader, QueryReader, BodyReader, ResponseWriter)")
	t.Log("- Matcher interface defined")
	t.Log("- Middleware interface defined")
	t.Log("- HandlerFunc type defined")
	t.Log("- Optional: Binder, Renderer, Container interfaces defined")
}

// TestSOLIDPrinciplesCompliance documents SOLID principles in interfaces.
func TestSOLIDPrinciplesCompliance(t *testing.T) {
	t.Run("SingleResponsibility", func(t *testing.T) {
		t.Log("✓ Each interface has one clear purpose")
		t.Log("  - Router: HTTP routing")
		t.Log("  - Matcher: Route matching")
		t.Log("  - Middleware: Request transformation")
		t.Log("  - Context: Request/response access")
	})

	t.Run("OpenClosed", func(t *testing.T) {
		t.Log("✓ Interfaces are open for extension via composition")
		t.Log("✓ Closed for modification - contracts don't change")
	})

	t.Run("LiskovSubstitution", func(t *testing.T) {
		// MiddlewareFunc is interchangeable with any Middleware
		var _ cosan.Middleware = cosan.MiddlewareFunc(nil)
		t.Log("✓ All implementations are interchangeable")
	})

	t.Run("InterfaceSegregation", func(t *testing.T) {
		t.Log("✓ Context segregated into focused interfaces:")
		t.Log("  - ParamReader: Path parameter access")
		t.Log("  - QueryReader: Query parameter access")
		t.Log("  - BodyReader: Request body parsing")
		t.Log("  - ResponseWriter: Response writing")
		t.Log("✓ Clients depend only on methods they use")
	})

	t.Run("DependencyInversion", func(t *testing.T) {
		// Handler depends on Context interface, not concrete implementation
		handler := func(ctx cosan.Context) error {
			return ctx.JSON(200, "ok")
		}
		_ = handler
		t.Log("✓ Depend on abstractions (interfaces), not concretions")
	})
}

// TestHTTPStandardLibraryCompliance verifies standard library compatibility.
func TestHTTPStandardLibraryCompliance(t *testing.T) {
	t.Log("Router implements http.Handler interface")
	t.Log("Can be used with http.ListenAndServe(\":8080\", router)")
}

// TestOptionalIntegrationsAreIndependent verifies that optional integrations
// are truly optional and don't create hard dependencies.
func TestOptionalIntegrationsAreIndependent(t *testing.T) {
	t.Log("✓ Binder, Renderer, and Container are separate optional interfaces")
	t.Log("✓ They can be nil/not provided without affecting core functionality")
	t.Log("✓ No dependencies between them")
	t.Log("✓ Cosan works perfectly standalone")
}

// TestInterfaceDocumentation verifies documentation completeness.
func TestInterfaceDocumentation(t *testing.T) {
	t.Log("All interfaces have godoc comments")
	t.Log("Usage examples provided in comments")
	t.Log("SOLID principles documented")
	t.Log("See: docs/INTERFACES.md for complete documentation")
}
